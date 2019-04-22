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
	Resources chan<- Resource //爬取资源
	Results   chan<- Result   //爬取结果
	Urls      <-chan string   //请求链接
}

//下载器任务
type DownLoadTask struct {
	Name     string
	DirPath  string
	Resource <-chan Resource //chan,协程使用
}

//数据持久化任务
type PersistenceTask struct {
	Name        string
	Persistence Persistence
	Results     <-chan Result
}

func CreateSpiderTask(resources chan Resource, results chan Result, urls chan string) Task {
	task := SpiderTask{
		Name:      "default downLoad task",
		Resources: resources,
		Results:   results,
		Urls:      urls,
	}
	return &task
}

func (this *SpiderTask) Start() {
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				//从channel中取出url进行抓取
				url := <-this.Urls
				response, err := http.Get(url)
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
					queryParams := QueryParams(url)
					this.Resources <- Resource{
						Url:  tvs[e].Image,
						Type: queryParams["type"],
						Tag:  queryParams["tag"],
					}
				}
				response.Body.Close()
			}
		}()
	}
	log.Info("[蜘蛛任务启动完成]")
}

func CreateDownLoadTask(dirPath string, resouce chan Resource) Task {
	task := DownLoadTask{
		Name:     "default downLoad task",
		DirPath:  dirPath,
		Resource: resouce,
	}
	return &task
}

func (this *DownLoadTask) Start() {
	for i := 0; i < 1000; i++ {
		go func() {
			for true {
				//从channel取得图片url
				resource := <-this.Resource
				//切割url得到文件名
				urlSplits := strings.Split(resource.Url, "/")
				imgName := urlSplits[len(urlSplits)-1]
				log.Info("[正在下载图片]:", imgName)
				response, err := http.Get(resource.Url)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				// 计算图片保存路径
				savePath := this.DirPath + "/" + resource.Type + "/" + resource.Tag
				_ = os.MkdirAll(savePath, 0777)
				out, _ := os.Create(savePath + "/" + imgName)
				_, _ = io.Copy(out, response.Body)
				log.Info("[图片]:", imgName, "下载完成")
				response.Body.Close()
				out.Close()
			}
		}()
	}
	log.Info("[多线程下载器启动完成]")
}

func CreatePersistenceTask(persistence Persistence, results chan Result) Task {
	task := PersistenceTask{
		Name:        "default persistencee Task",
		Persistence: persistence,
		Results:     results,
	}
	return &task
}

func (this *PersistenceTask) Start() {
	for i := 0; i < 1000; i++ {
		go func() {
			for true {
				tv := <-this.Results
				id, err := this.Persistence.SaveOne(tv)
				if id == nil {
					log.Error("[持久化配置]:失效,routine退出")
					return
				}
				if err != nil {
					log.Error(err.Error())
					continue
				}
			}
		}()
	}
	log.Info("[持久化任务启动完成]")

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

func parseUrls(tp string, urls chan string) {
	//获取tag,解析为url提供蜘蛛任务进行爬取
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
