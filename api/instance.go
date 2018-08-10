package api

import (
	"fmt"
	"io/ioutil"
	"sync"
	"gopkg.in/yaml.v2"
	"strings"
	"path/filepath"
	"errors"
)

type Instance struct {
	Name string
	Limits struct {
		Calls int32
		Frame string
	}
	Credentials map[string]struct {
		Secret string
	}
	Endpoints map[string]struct {
		Path   string
		Method string
		Block  bool
	}
}

var configs = make(map[string]Instance)
var loadOnce sync.Once

func loadAll() {
	loadOnce.Do(func() {
		files, err := ioutil.ReadDir("./example_config")
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".yml" {
				conf := Instance{}
				content, _ := ioutil.ReadFile(fmt.Sprintf("./example_config/%s", file.Name()))
				if err := yaml.Unmarshal(content, &conf); err != nil {
					panic(err)
				}
				apiName := strings.Split(file.Name(), ".")[0]
				configs[apiName] = conf
			}
		}
	})
}

func GetApiConfig(name string) (*Instance, error) {
	loadAll()
	config, ok := configs[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Requested api (%s) does not exist", name))
	}
	return &config, nil
}

func (i *Instance) Call(endpoint string) interface{} {
	return nil
}
