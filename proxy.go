package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var proxyPool []string

// registerProxy register proxy into proxy pool
func registerProxy(addr string) error {
	temp := strings.Split(addr, ":")
	if len(temp) != 2 {
		return fmt.Errorf("invalid net addr:%v", addr)
	}
	proxyPool = append(proxyPool, addr)
	return nil
}

// Request use http method to request the spec url
func Request(method, url string, body io.Reader) (resp *http.Response, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", userAgent())
	return proxyClient().Do(request)
}

func proxyClient() *http.Client {
	if len(proxyPool) == 0 {
		return http.DefaultClient
	}
	randNo := rand.Intn(len(proxyPool))
	proxyAddr := proxyPool[randNo]
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return nil
	}

	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		DialContext: func(context context.Context, netw, addr string) (net.Conn, error) {
			defer context.Done()
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(10))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}

	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func userAgent() string {
	agent := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := len(agent)
	return agent[r.Intn(length)]
}
