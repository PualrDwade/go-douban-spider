package main

import "encoding/json"

// 实体model
type Result struct {
	Id       string `json:"id"`       //豆瓣ID
	Rate     string `json:"rate"`     //评分
	Title    string `json:"title"`    //电视剧标题
	Url      string `json:"url"`      //TV链接
	Playable bool   `json:"playable"` //是否可观看
	Image    string `json:"cover"`    //封面链接
	IsNew    bool   `json:"is_new"`   //是否是新TV
}

// 资源model
type Resource struct {
	Url  string
	Type string
	Tag  string
}

// 解析json
func ParseJson(content []byte) ([]Result, error) {
	// 首先使用map 接受json内容
	var result map[string]interface{}
	err := json.Unmarshal(content, &result)
	// 从map中取出需要的内容
	jsonTVs := result["subjects"].([]interface{})
	// 作为返回结果
	var resultTvs []Result
	jsonString, err := json.Marshal(jsonTVs)
	err = json.Unmarshal(jsonString, &resultTvs)
	return resultTvs, err
}
