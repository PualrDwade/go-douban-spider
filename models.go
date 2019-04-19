package main

import "encoding/json"

// tv-model
type Item struct {
	Id       string `json:"id"`       //豆瓣ID
	Rate     string `json:"rate"`     //评分
	Title    string `json:"title"`    //电视剧标题
	Url      string `json:"url"`      //TV链接
	Playable bool   `json:"playable"` //是否可观看
	Image    string `json:"cover"`    //封面链接
	IsNew    bool   `json:"is_new"`   //是否是新TV
}

// 解析json
func ParseJson(content []byte) ([]Item, error) {
	// 首先使用map 接受json内容
	var result map[string]interface{}
	err := json.Unmarshal(content, &result)
	// 从map中取出需要的内容
	jsonTVs := result["subjects"].([]interface{})
	// 作为返回结果
	var resultTvs []Item
	jsonString, err := json.Marshal(jsonTVs)
	err = json.Unmarshal(jsonString, &resultTvs)
	return resultTvs, err
}
