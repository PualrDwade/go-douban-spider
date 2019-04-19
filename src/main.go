package main

import (
	"github.com/siddontang/go/log"
	"time"
)

func main() {

	// 1.蜘蛛任务->得到results-tv(chan)*2 -生产者
	// 2.持久化任务->消费results1-tv(chan)->持久化->消费者
	// 3.下载器任务->消费results2-tv(chan)->下载model中的图片资源->消费者

	results := make(chan TV)       //对应results1
	resources := make(chan string) //对应results2
	finish := make(chan bool)      //控制信号

	//1.启动下载器任务
	downLoadTask := CreateDownLoadTask("./download", resources, finish)
	go downLoadTask.Start()

	//2.启动持久化任务
	persistenceTask := CreatePersistenceTask(CreateMonoPersistence(), results)
	go persistenceTask.Start()

	//3.启动蜘蛛任务
	spiderTask := CreateSpiderTask(resources, results, finish)
	go spiderTask.Start()

	time.Sleep(time.Second * 20)
	log.Info("爬虫程序退出")
}
