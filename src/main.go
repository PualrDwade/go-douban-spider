package main

import (
	"encoding/json"
	"github.com/siddontang/go/log"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

//豆瓣tv标签链接
const TAGS string = "https://movie.douban.com/j/search_tags?type=tv"

//豆瓣爬取链接模版
// type = tv:爬取电视剧
// tag  = ??? :爬取的类型,可以通过tags链接获取到
// page_limit = xxx :每爬取一次显示的items数量
// page_start = xxx: 当前爬取位置,设置为爬取次数i*page_limits

var UrlTemplate = "https://movie.douban.com/j/search_subjects?type=tv&tag=????&page_limit=xxx&page_start=0"

var URL = "https://movie.douban.com/j/search_subjects?type=tv&tag=热门&page_limit=" + strconv.Itoa(math.MaxInt32) + "&page_start=0"

// 数据持久化对象
var dao Persistence

// 爬取的标签类型
var tags []string

// tv
type TV struct {
	Id       string `json:"id"`       //豆瓣ID
	Rate     string `json:"rate"`     //评分
	Title    string `json:"title"`    //电视剧标题
	Url      string `json:"url"`      //TV链接
	Playable bool   `json:"playable"` //是否可观看
	Image    string `json:"cover"`    //封面链接
	IsNew    bool   `json:"is_new"`   //是否是新TV
}

// 解析json
func ParseJson(content []byte) ([]TV, error) {
	// 首先使用map 接受json内容
	var result map[string]interface{}
	err := json.Unmarshal(content, &result)
	// 从map中取出需要的内容
	jsonTVs := result["subjects"].([]interface{})
	// 作为返回结果
	var resultTvs []TV
	jsonString, err := json.Marshal(jsonTVs)
	err = json.Unmarshal(jsonString, &resultTvs)
	return resultTvs, err
}

func main() {
	// 发起网络请求,请求tv的url
	response, err := http.Get(URL)
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

	// 从map中取出需要的内容
	jsonTVs := result["subjects"].([]interface{})
	// 持久化,使用mongodb
	dao, err = CreateMonoPersistence()
	if err != nil {
		log.Error(err.Error())
		return
	}
	complete, err := dao.Save(jsonTVs)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("插入many mongodb完成,数量:", len(complete.InsertedIDs), "--返回结果:", complete.InsertedIDs)

	//解析为model切片,供程序后续使用
	tvs, err := ParseJson(body)
	if err != nil {
		log.Info(err.Error())
		return
	}
	log.Info("爬取到内容:", tvs)

	//启动多线程图片加载器,进行资源的下载
	urls := make(chan string) //使用make构造引用类型
	finish := make(chan bool) //结束flag
	downLoadTask := CreateDownLoadTask("./download", urls, finish)
	go downLoadTask.Start() //启动下载任务
	for i := 0; i < len(tvs); i++ {
		url := tvs[i].Image
		urls <- url //into channel
	}
	for i := 0; i < len(tvs); i++ {
		<-finish
	}
	log.Info("爬虫程序退出")
}
