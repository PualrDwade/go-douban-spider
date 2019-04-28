package main

import (
	"testing"
)

func TestCreateDownLoadTask(t *testing.T) {
	resource := make(chan Resource)
	task := CreateDownLoadTask("./download", resource)
	if task == nil {
		t.Fail()
	}
}

func TestDownLoadTask_Start(t *testing.T) {
	resource := make(chan Resource)
	task := CreateDownLoadTask("./download", resource)
	if task == nil {
		t.Fail()
	}
	go task.Start()
	for i := 0; i < 10; i++ {
		resource <- Resource{
			Url:  "https://www.baidu.com/img/bd_logo1.png",
			Type: "tv",
			Tag:  "hot",
			Name: "百度",
		}
	}
}

func TestCreatePrepareTask(t *testing.T) {
	urls := make(chan string)
	prepareTask := CreatePrepareTask(urls)
	if prepareTask == nil {
		t.Fail()
	}
}

func TestPrepareTask_Start(t *testing.T) {
	urls := make(chan string)
	prepareTask := CreatePrepareTask(urls)
	go prepareTask.Start()
	for i := 0; i < 10; i++ {
		<-urls
	}
}

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
