package main

import (
	"github.com/siddontang/go/log"
	"io"
	"net/http"
	"os"
	"strings"
)

//模块只暴露接口,隐藏具体实现,具体实现由对应的工厂方法注入依赖模块中
type Task interface {
	Start()
}

type downLoadTask struct {
	Name    string
	DirPath string
	URL     chan string //chan,协程使用
	Finish  chan bool   //chan,作为工作停止信号
}

func CreateDownLoadTask(dirPath string, url chan string, finish chan bool) Task {
	task := downLoadTask{
		Name:    "defalt downLoad task",
		DirPath: dirPath,
		URL:     url,
		Finish:  finish,
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
				url := <-this.URL                    //从channel取得图片url
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
	log.Info("[所有图片下载完成]")
}
