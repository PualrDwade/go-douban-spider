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

var config Config

func globalConfig() Config {
	return config
}

func loadConfigOrDie() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("load config.json failed: %v, use default config", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("unmarshal config.json failed: %v, use default config\n", err)
	}
	// register proxy
	for name, addr := range config.ProxyPool {
		registerProxy(name, addr)
	}
}
