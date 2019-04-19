package main

import (
	"encoding/json"
	"github.com/siddontang/go/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Task interface {
	Start()
}

//爬虫任务
type SpiderTask struct {
	Name      string
	Resources chan string //爬取资源
	Results   chan Item   //爬取结果
	Urls      chan string //请求链接
	Finish    chan bool   //结束标签
}

func CreateSpiderTask(resources chan string, results chan Item, urls chan string, finish chan bool) Task {
	task := SpiderTask{
		Name:      "default downLoad task",
		Resources: resources,
		Results:   results,
		Urls:      urls,
		Finish:    finish,
	}
	return &task
}

func (this *SpiderTask) Start() {
	go func() {
		for true {
			response, err := http.Get(<-this.Urls)
			if err != nil {
				log.Info(err.Error())
				return
			}
			body, err := ioutil.ReadAll(response.Body)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				log.Error(err.Error())
				return
			}
			//解析为model切片,供程序后续使用
			tvs, err := ParseJson(body)
			if err != nil {
				log.Info(err.Error())
				return
			}
			log.Info("[爬取到内容]:", tvs)
			// 存入chnnel
			for e := range tvs {
				this.Results <- tvs[e]
				this.Resources <- tvs[e].Image
			}
			response.Body.Close()
		}
	}()
	log.Info("[蜘蛛任务启动完成]")
}

//下载器任务
type DownLoadTask struct {
	Name     string
	DirPath  string
	Resource chan string //chan,协程使用
	Finish   chan bool   //chan,作为工作停止信号
}

func CreateDownLoadTask(dirPath string, url chan string, finish chan bool) Task {
	task := DownLoadTask{
		Name:     "default downLoad task",
		DirPath:  dirPath,
		Resource: url,
		Finish:   finish,
	}
	return &task
}

func (this *DownLoadTask) Start() {
	// 首先创建文件夹
	err := os.Mkdir(this.DirPath, 0777)
	if err != nil {
		//如果已经存在则直接进入
		log.Info("文件夹", this.DirPath, "已经存在,直接进入.")
	}
	err = os.Chdir(this.DirPath)
	if err != nil {
		log.Error(err.Error())
		return
	}
	const routines = 10000 //设定rouines数量
	log.Info("[多线程下载器启动完成]")
	for i := 0; i < routines; i++ {
		go func() {
			for true {
				url := <-this.Resource               //从channel取得图片url
				urlSplits := strings.Split(url, "/") //切割url得到文件名
				imgName := urlSplits[len(urlSplits)-1]
				log.Info("[正在下载图片]:", imgName)
				response, err := http.Get(url)
				if err != nil {
					log.Error(err.Error())
					return
				}
				out, err := os.Create(imgName)
				if err != nil {
					log.Error(err.Error())
					return
				}
				_, err = io.Copy(out, response.Body)
				if err != nil {
					log.Error(err)
					return
				}
				log.Info("[图片]:", imgName, "下载完成")
				this.Finish <- true
				response.Body.Close()
				out.Close()
			}
		}()
	}
}

//数据持久化任务
type PersistenceTask struct {
	Name        string
	Persistence Persistence
	Results     chan Item
}

func CreatePersistenceTask(persistence Persistence, results chan Item) Task {
	task := PersistenceTask{
		Name:        "default persistencee Task",
		Persistence: persistence,
		Results:     results,
	}
	return &task
}

func (this *PersistenceTask) Start() {
	log.Info("[持久化任务启动完成]")
	for i := 0; i < 100; i++ {
		go func() {
			for true {
				tv := <-this.Results
				_, err := this.Persistence.SaveOne(tv)
				if err != nil {
					log.Error(err.Error())
					continue
				}
			}
		}()
	}
}

type PrepareTask struct {
	Name string
	Urls chan string //爬取的链接
}

func CreatePrepareTask(urls chan string) Task {
	task := PrepareTask{
		Name: "default PrepareTask",
		Urls: urls,
	}
	return &task
}

func (this *PrepareTask) Start() {
	go parseUrls("tv", this.Urls)
	go parseUrls("movie", this.Urls)
	log.Info("[超链接预处理器启动完成]")
}

//获取tag,解析为url提供蜘蛛任务进行爬取
func parseUrls(tp string, urls chan string) {
	type Tags struct {
		Tags []string `json:"tags"`
	}
	response, err := http.Get("https://movie.douban.com/j/search_tags?type=" + tp)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var tags Tags
	body, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &tags)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for e := range tags.Tags {
		var url = "https://movie.douban.com/j/search_subjects?type=" + tp + "&tag=" + tags.Tags[e] + "&page_limit=" + strconv.Itoa(1000000) + "&page_start=0"
		log.Info("[parse url:]", url)
		urls <- url //放入channel中
	}
	defer response.Body.Close()
}
