package main

import (
	"strings"
)

// QueryParams 解析url的查询参数工具方法
// expamle:https://movie.douban.com/j/search_subjects?type=tv&tag=hot
// result: {"type":"tv","tag":"hot"}
func QueryParams(url string) map[string]string {
	strs := strings.Split(url, "/")
	handle := strs[len(strs)-1] //得到最后一个"/"切片
	strs = strings.Split(handle, "?")
	handle = strs[len(strs)-1]        //type=tv&tag=hot
	strs = strings.Split(handle, "&") //["type=tv","tag=hot"]
	queryParams := make(map[string]string)
	for e := range strs {
		key := strings.Split(strs[e], "=")[0]
		value := strings.Split(strs[e], "=")[1]
		queryParams[key] = value
	}
	return queryParams
}
