package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

type Configuration struct {
	RedisServer RedisServer `yaml:"redisServer"`
	Service     Service     `yaml:"service"`
	MysqlServer MysqlServer `yaml:"mysqlServer"`
}

type Service struct {
	ServerPort string `yaml:"serverPort"`
}

type MysqlServer struct {
	Address  string `yaml:"address"`
	UserName string `yaml:"userName"`
	PassWord string `yaml:"passWord"`
}

type RedisServer struct {
	Address       string `yaml:"address"`
	RedisPassword string `yaml:"redisPassword"`
}

func NewConfig() *Configuration {
	return &Configuration{}
}

var once sync.Once
var deFaultConfigPath string
var rootConfig *Configuration

func InitConfig() *Configuration {
	once.Do(func() {
		if _, err := initConfig(); err != nil {
			panic(err)
		}
	})
	return Config()
}

func Config() *Configuration {
	if rootConfig == nil {
		panic("init config error")
	}
	return rootConfig
}

func getDefaultConfigPath() string {
	if deFaultConfigPath != "" {
		return deFaultConfigPath
	}
	configFile := ""
	if runtime.GOOS == "windows" {
		configFile = os.Getenv("GOPATH") + "\\src\\goFix\\logger\\goFix.log"
	} else {
		configFile = os.Getenv("GOPATH") + "/goFix/logger/goFix.log"
	}

	configFile = os.Getenv("GOPATH") + "/goFix/config/config.yaml"
	return configFile
}

func initConfig() (config *Configuration, err error) {
	conf := NewConfig()
	configPath := getDefaultConfigPath()
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	rootConfig = conf
	return
}
