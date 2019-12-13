package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Task 启动任务接口
type Task interface {
	Start()
}

// PrepareTask 预处理任务结构
type PrepareTask struct {
	Name string
	Urls chan string //爬取的链接
}

// CreatePrepareTask 工厂方法，返回预处理任务
func CreatePrepareTask(urls chan string) Task {
	task := PrepareTask{
		Name: "default PrepareTask",
		Urls: urls,
	}
	return &task
}

// Start 启动url预处理任务
func (task *PrepareTask) Start() {
	parseURL := func(kind string) {
		//获取tag,解析为url提供蜘蛛任务进行爬取
		type Tags struct {
			Tags []string `json:"tags"`
		}
		response, err := Request(http.MethodGet, "https://movie.douban.com/j/search_tags?type="+kind, nil)
		if err != nil {
			log.Fatalf("request failed: %v", err)
		}
		defer response.Body.Close()

		var tags Tags
		body, err := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &tags)
		if err != nil {
			log.Fatalf("unmarshal failed: %v", err)
			return
		}
		for e := range tags.Tags {
			var url = "https://movie.douban.com/j/search_subjects?type=" + kind + "&tag=" + tags.Tags[e] +
				"&page_limit=" + strconv.Itoa(math.MaxInt32) + "&page_start=0"
			log.Println("[parse url]:", url)
			task.Urls <- url //放入channel中
		}
	}
	go parseURL("tv")
	go parseURL("movie")
	log.Println("[PrepareTask]:超链接预处理器启动完成")
}

// SpiderTask 爬虫任务
type SpiderTask struct {
	Name      string
	Resources chan<- Resource //爬取资源
	Results   chan<- Result   //爬取结果
	Urls      <-chan string   //请求链接
}

// CreateSpiderTask 爬虫任务工厂方法
func CreateSpiderTask(resources chan Resource, results chan Result, urls chan string) Task {
	task := SpiderTask{
		Name:      "default downLoad task",
		Resources: resources,
		Results:   results,
		Urls:      urls,
	}
	return &task
}

// Start 启动爬虫任务
func (task *SpiderTask) Start() {
	for i := 0; i < globalConfig().TaskRoutines; i++ {
		go func() {
			for {
				// 从channel中取出url进行抓取 Start 启动下载器任务
				url := <-task.Urls
				response, err := Request(http.MethodGet, url, nil)
				if err != nil {
					log.Fatal(err)
				}
				body, err := ioutil.ReadAll(response.Body)
				result := make(map[string]interface{})
				err = json.Unmarshal(body, &result)
				if err != nil {
					log.Fatalf("unmarshal failed: %v", err)
				}
				//解析为model切片,供程序后续使用
				tvs, err := ParseJSON(body)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("[SpiderTask]:爬取到内容->", tvs)
				// 存入chnnel
				for e := range tvs {
					task.Results <- tvs[e]
					// Start 启动下载器任务
					queryParams := QueryParams(url)
					task.Resources <- Resource{
						URL:  tvs[e].Image,
						Type: queryParams["type"],
						Tag:  queryParams["tag"],
						Name: tvs[e].Title,
					}
				}
				response.Body.Close()
			}
		}()
	}
	log.Println("[SpiderTask]:蜘蛛任务启动完成")
}

// DownLoadTask 下载器任务
type DownLoadTask struct {
	Name     string
	DirPath  string
	Resource <-chan Resource //chan,协程使用
}

// CreateDownLoadTask 工厂方法 返回下载器
func CreateDownLoadTask(dirPath string, resouce chan Resource) Task {
	task := DownLoadTask{
		Name:     "default downLoad task",
		DirPath:  dirPath,
		Resource: resouce,
	}
	return &task
}

// Start 启动下载器任务
func (task *DownLoadTask) Start() {
	for i := 0; i < globalConfig().TaskRoutines; i++ {
		go func() {
			for true {
				// 从channel取得图片url Start 启动下载器任务
				resource := <-task.Resource
				//切割url得到文件类型
				urlSplits := strings.Split(resource.URL, ".")
				imgFileType := urlSplits[len(urlSplits)-1]
				imgName := resource.Name + "." + imgFileType
				response, err := Request(http.MethodGet, resource.URL, nil)
				if err != nil {
					log.Println(err)
					continue
				}
				// 判断图片保存路径 Start 启动下载器任务
				savePath := task.DirPath + "/" + resource.Type + "/" + resource.Tag
				_ = os.MkdirAll(savePath, 0777)
				out, _ := os.Create(savePath + "/" + imgName)
				io.Copy(out, response.Body)
				log.Println("[DownLoadTask]:图片", imgName, "下载完成")
				response.Body.Close()
				out.Close()
			}
		}()
	}
	log.Println("[DownLoadTask]多线程下载器启动完成")
}

// PersistenceTask 数据持久化任务
type PersistenceTask struct {
	Name        string
	Persistence Persistence
	Results     <-chan Result
}

// CreatePersistenceTask 工厂方法 返回Task
func CreatePersistenceTask(persistence Persistence, results chan Result) Task {
	task := PersistenceTask{
		Name:        "default persistencee Task",
		Persistence: persistence,
		Results:     results,
	}
	return &task
}

// Start 启动持久化认任务
func (task *PersistenceTask) Start() {
	for i := 0; i < globalConfig().TaskRoutines; i++ {
		go func() {
			for true {
				tv := <-task.Results
				id, err := task.Persistence.SaveOne(tv)
				if id == nil || err != nil {
					log.Println("[persistenceTask]:持久化失败")
					continue
				}
			}
		}()
	}
	log.Println("[PersistenceTask]:持久化任务启动完成")
}
