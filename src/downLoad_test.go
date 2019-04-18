package main

import (
	"sync"
	"testing"
	"time"
)

func TestCreateDownLoadTask(t *testing.T) {
	url := make(chan string)
	task := CreateDownLoadTask("./download", url)
	if task == nil {
		t.Fail()
	}
}

func TestDownLoadTask_Start(t *testing.T) {
	url := make(chan string)
	task := CreateDownLoadTask("./download", url)
	if task == nil {
		t.Fail()
	}
	var lath sync.WaitGroup
	lath.Add(1)
	go task.Start(&lath) // 这里需要传递引用
	for i := 0; i < 100; i++ {
		url <- "https://www.baidu.com/img/bd_logo1.png"
	}
	time.Sleep(time.Second * 3)
	lath.Done()
}
