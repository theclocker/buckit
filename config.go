package main

import (
	"fmt"
	"io/ioutil"
	"sync"
)

type ApiConfig struct {
	Name   string
	Limits struct {
		Calls int32
		Frame string
	}
	Credentials map[string]struct {
		Secret string
	}
	Routes map[string]struct {
		Path   string
		Method string
		Block  bool
	}
}

var configs []ApiConfig
var once sync.Once

func LoadApiConfigs() {
	files, err := ioutil.ReadDir("./example_config")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
