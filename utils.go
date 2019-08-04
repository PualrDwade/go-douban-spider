package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// QueryParams 解析url的查询参数工具方法
// expamle:https://movie.douban.com/j/search_subjects?type=tv&tag=hot
// result: {"type":"tv","tag":"hot"}
func QueryParams(url string) map[string]string {
	strs := strings.Split(url, "/")
	handle := strs[len(strs)-1] //得到最后一个"/"切片
	strs = strings.Split(handle, "?")
	handle = strs[len(strs)-1]        //type=tv&tag=hot
	strs = strings.Split(handle, "&") //["type=tv","tag=hot"]
	queryParams := make(map[string]string)
	for e := range strs {
		key := strings.Split(strs[e], "=")[0]
		value := strings.Split(strs[e], "=")[1]
		queryParams[key] = value
	}
	return queryParams
}

// RandomUserAgent 随机返回一个user-agent
func RandomUserAgent() string {
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

// WrapperRequest返回预设http请求
func WrapperRequest(method, url string, body io.Reader) (resp *http.Response, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", RandomUserAgent())
	return http.DefaultClient.Do(request)
}
