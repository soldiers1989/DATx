package server

import (
	"fmt"

	"github.com/go-ini/ini"
)

func GetConfig() (*ini.File, error) {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetCfgProducerName() string {
	var result string
	cfg, err := GetConfig()
	if err != nil {
		fmt.Printf("Get config err:%v\n", err)
		return result
	}

	result = cfg.Section("").Key("producer-name").String()
	return result
}
