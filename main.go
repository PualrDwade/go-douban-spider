package main

import (
	"log"
	"time"
)

func main() {
	// load configuration for application
	config := loadConfig()

	// 1.预处理器->解析url->urls(chan)-生产者
	// 2.蜘蛛任务->得到results-tv(chan)*2 -(消费者,消费urls)+(生产者)
	// 3.持久化引擎->消费results1-tv(chan)->持久化->消费者
	// 4.下载器->消费results2-tv(chan)->下载model中的图片资源->消费者
	urls := make(chan string, config.ChanSize)
	results := make(chan Result, config.ChanSize)
	resources := make(chan Resource, config.ChanSize)

	//1.启动下载器任务
	downLoadTask := CreateDownLoadTask("./download", resources)
	go downLoadTask.Start()

	//2.启动持久化任务
	persistenceTask := CreatePersistenceTask(CreateMonoPersistence(), results)
	go persistenceTask.Start()

	//3.启动蜘蛛任务
	spiderTask := CreateSpiderTask(resources, results, urls)
	go spiderTask.Start()

	//4.启动预处理器任务
	prepareTask := CreatePrepareTask(urls)
	go prepareTask.Start()

	executeTime := config.Duration
	time.Sleep(config.Duration * time.Second)
	log.Println("[爬虫程序退出]:爬取时间:", executeTime, "秒")
}
