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
