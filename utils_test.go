package main

import (
	"testing"
)

func TestQueryParams(t *testing.T) {
	var url = "https://movie.douban.com/j/search_subjects?type=tv&tag=hot"
	queryParams := QueryParams(url)
	t.Log(queryParams)
	if queryParams["type"] != "tv" || queryParams["tag"] != "hot" {
		t.Fail()
	}
}

func TestQueryParams2(t *testing.T) {
	var url = "https://movie.douban.com/j/?type=电影&tag=热门"
	queryParams := QueryParams(url)
	t.Log(queryParams)
	if queryParams["type"] != "电影" || queryParams["tag"] != "热门" {
		t.Fail()
	}
}

func TestQueryParams3(t *testing.T) {
	var url = "https://movie.douban.com/j/?name=&tag=hot"
	queryParams := QueryParams(url)
	t.Log(queryParams)
	if queryParams["name"] != "" || queryParams["tag"] != "hot" {
		t.Fail()
	}
}
