package main

import (
	"encoding/json"
	"github.com/siddontang/go/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//模块只暴露接口,隐藏具体实现,具体实现由对应的工厂方法注入依赖模块中
type Task interface {
	Start()
}

//爬虫任务
type spiderTask struct {
	Name      string
	Resources chan string //爬取资源
	Results   chan TV     //爬取结果
	Finish    chan bool
}

func CreateSpiderTask(resources chan string, results chan TV, finish chan bool) Task {
	task := spiderTask{
		Name:      "default downLoad task",
		Resources: resources,
		Results:   results,
		Finish:    finish,
	}
	return &task
}

func (this *spiderTask) Start() {
	// 发起网络请求,请求tv的url
	response, err := http.Get(StartLink)
	if err != nil {
		log.Info(err.Error())
		return
	}
	defer response.Body.Close()

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
}

//下载器任务
type downLoadTask struct {
	Name     string
	DirPath  string
	Resource chan string //chan,协程使用
	Finish   chan bool   //chan,作为工作停止信号
}

func CreateDownLoadTask(dirPath string, url chan string, finish chan bool) Task {
	task := downLoadTask{
		Name:     "default downLoad task",
		DirPath:  dirPath,
		Resource: url,
		Finish:   finish,
	}
	return &task
}

func (this *downLoadTask) Start() {
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
	const routines = 100 //设定rouines数量
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
type persistenceTask struct {
	Name        string
	Persistence Persistence
	Results     chan TV
}

func CreatePersistenceTask(persistence Persistence, results chan TV) Task {
	task := persistenceTask{
		Name:        "default persistencee Task",
		Persistence: persistence,
		Results:     results,
	}
	return &task
}

func (this *persistenceTask) Start() {
	for true {
		tv := <-this.Results
		_, err := this.Persistence.SaveOne(tv)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
}
