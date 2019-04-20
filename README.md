# go-douban-spider

![](https://img.shields.io/badge/语言-go-blue.svg)  ![](https://img.shields.io/badge/爬虫-豆瓣-gren.svg)

## 简介

基于go的豆瓣电影多线程爬虫


## 架构
1. 预处理器->解析url并添加可爬取超链到channel中->作为消费者
2. 蜘蛛任务->得到爬取结果,解析可利用资源->存入channel->作为消费者,消费超链channel,同时作为生产者
3. 持久化引擎->消费爬取结果,进行持久化->作为消费者
4. 下载器->消费可利用资源,下载model中的图片资源->作为消费者


## 启动流程
```go
//1.启动下载器任务
downLoadTask := CreateDownLoadTask("./download", resources, finish)
go downLoadTask.Start()

//2.启动持久化任务
persistenceTask := CreatePersistenceTask(CreateMonoPersistence(), results)
go persistenceTask.Start()

//3.启动蜘蛛任务
spiderTask := CreateSpiderTask(resources, results, urls, finish)
go spiderTask.Start()

//4.启动预处理器任务
prepareTask := CreatePrepareTask(urls)
go prepareTask.Start()
```
