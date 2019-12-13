package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

// Config holds add configuration for app
type Config struct {
	PersistanceModel string            `json:"persistance_model"`
	MongoURL         string            `json:"mongo_url"`
	MySQLURL         string            `json:"mysql_url"`
	ProxyPool        map[string]string `json:"proxy_pool"`
	Duration         time.Duration     `json:"duration"`
	ChanSize         int               `json:"chan_size"`
	TaskRoutines     int               `json:"task_routines"`
}

var config = Config{
	PersistanceModel: "mongo",
	MongoURL:         "localhost:27017",
	MySQLURL:         "localhost:3306",
	ProxyPool:        make(map[string]string),
	Duration:         60,
	ChanSize:         5000,
	TaskRoutines:     100,
}

func globalConfig() Config {
	return config
}
func loadConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("load config.json failed, use default config")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("unmarshal config.json failed, use default config")
	}
	// register proxy
	for name, addr := range config.ProxyPool {
		registerProxy(name, addr)
	}
}
