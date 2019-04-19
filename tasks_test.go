package main

import (
	"testing"
)

func TestCreateDownLoadTask(t *testing.T) {
	url := make(chan string)
	finish := make(chan bool)
	task := CreateDownLoadTask("./download", url, finish)
	if task == nil {
		t.Fail()
	}
}

func TestDownLoadTask_Start(t *testing.T) {
	url := make(chan string)
	finish := make(chan bool)
	task := CreateDownLoadTask("./download", url, finish)
	if task == nil {
		t.Fail()
	}
	go task.Start()
	for i := 0; i < 10; i++ {
		url <- "https://www.baidu.com/img/bd_logo1.png"
	}
	for i := 0; i < 10; i++ {
		<-finish
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
