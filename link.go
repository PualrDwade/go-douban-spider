package main

import (
	"strconv"
)

//豆瓣tv标签链接
//const TAGS string = "https://movie.douban.com/j/search_tags?type=tv"

//豆瓣爬取链接模版
// type = tv:爬取电视剧
// tag  = ??? :爬取的类型,可以通过tags链接获取到
// page_limit = xxx :每爬取一次显示的items数量
// page_start = xxx: 当前爬取位置,设置为爬取次数i*page_limits

var StartLink = "https://movie.douban.com/j/search_subjects?type=tv&tag=热门&page_limit=" + strconv.Itoa(1000000) + "&page_start=0"
