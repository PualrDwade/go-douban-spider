package main

import (
	"github.com/siddontang/go/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// 多线程下载任务
type DownLoadTask struct {
	Name    string
	DirPath string
	URL     chan string //chan,协程使用
}

// 工厂方法提供任务构建
func CreateDownLoadTask(dirPath string, url chan string) *DownLoadTask {
	task := DownLoadTask{
		Name:    "defalt downLoad task",
		DirPath: dirPath,
		URL:     url,
	}
	return &task
}

// 开启下载协程,使用lath控制
func (this *DownLoadTask) Start(lath *sync.WaitGroup) {
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
	for i := 0; i < routines; i++ {
		var imgName = "downLoad-" + strconv.Itoa(i) + ".jpg"
		go func() {
			url := <-this.URL
			log.Info("正在下载图片:", url)
			response, err := http.Get(url)
			if err != nil {
				log.Error(err.Error())
			}
			defer response.Body.Close()
			out, err := os.Create(imgName)
			if err != nil {
				log.Error(err.Error())
			}
			defer out.Close()
			_, err = io.Copy(out, response.Body)
			if err != nil {
				log.Error(err)
			}
		}()
	}
	lath.Wait()
	log.Info("图片下载完成!")
}