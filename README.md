# go-douban-spider

![](https://img.shields.io/badge/语言-go-blue.svg)  ![](https://img.shields.io/badge/爬虫-豆瓣-gren.svg)

## 简介

基于go的豆瓣电影多线程爬虫,爬取豆瓣电影与电视剧内容、下载相应资源。支持自定义持久化(mongodb、mysql)。基于goroutine轻量级协程与channel通信机制实现高效并发,代码量很少，适合对于golang还不太熟悉的小伙伴学习使用~


## 架构&思路
- [预处理器]->解析url并添加可爬取超链到channel中->作为消费者

- [蜘蛛任务]->得到爬取结果,解析可利用资源->存入channel->作为消费者,消费超链channel,同时作为生产者

- [持久化引擎]->消费爬取结果,进行持久化->作为消费者

- [下载器]->消费可利用资源,下载model中的图片资源->作为消费者

## 项目配置

项目配置文件对应config.json
```json
// config.json
{
    "persistance_model": "mongo",
    "mongo_url": "localhost:27017",
    "mysql_url": "localhost:3306",
    "proxy_pool": {
        "proxy1": "11.11.11.11:7777"
     },
    "duration": 60,
    "chan_size": 1000,
    "task_routines": 30
}
```
项目的默认配置如下:
```go
var config = Config{
	PersistanceModel: "mongo",
	MongoURL:         "localhost:27017",
	MySQLURL:         "localhost:3306",
	ProxyPool:        make(map[string]string),
	Duration:         60,
	ChanSize:         5000,
	TaskRoutines:     100,
}
```

## 核心流程
```go
// 1.预处理器->解析url->urls(chan)-生产者
// 2.蜘蛛任务->得到results-tv(chan)*2 -(消费者,消费urls)+(生产者)
// 3.持久化引擎->消费results1-tv(chan)->持久化->消费者
// 4.下载器->消费results2-tv(chan)->下载model中的图片资源->消费者
urls := make(chan string, 5000)
results := make(chan Result, 5000)
resources := make(chan Resource, 5000)

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
```

## 构建&使用

1.gopath初始化

```shell
cd $GOPATH/src
mkdir -p github.com/PualrDwade
cd github.com/PualrDwade
```

2.获取项目

```shell
git clone https://github.com/PualrDwade/go-douban-spider.git
cd go-douban-spider
```

3.编译

```makefile
make
```

4.运行

```makefile
make run
```

5.测试

```	makefile
make test
```

6.清理

```makefile
make clean
```

## 结果展示

![](https://i.loli.net/2019/08/04/hvAMSKXjFUWB3Vy.png)

## todo

- [x] 提供更加健壮的运行时机制,减少panic错误
- [x] 加入随机User-Agent，减少handshake error
- [x] 提供代理池功能,提高反反爬虫能力
- [x] 抽取公共配置到json文件中,提供定制化