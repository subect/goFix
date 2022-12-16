package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type Configuration struct {
	DevMode     bool        `yaml:"devMode"`
	Service     Service     `yaml:"service"`
	MysqlServer MysqlServer `yaml:"mysqlServer"`
	RedisServer RedisServer `yaml:"redisServer"`
	EsServer    EsServer    `yaml:"esServer"`
}

type EsServer struct {
	EsHost   string `yaml:"hosts"`
	EsUser   string `yaml:"esUser"`
	EsPasswd string `yaml:"esPassword"`
}

type Service struct {
	ServerPort string `yaml:"serverPort"`
}

type MysqlServer struct {
	Address       string `yaml:"address"`
	UserName      string `yaml:"userName"`
	PassWord      string `yaml:"passWord"`
	MysqlMaxDBs   int    `yaml:"mysqlMaxDBs"`
	DefaultDbName string `yaml:"defaultDbName"`
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
		if err := initConfig(); err != nil {
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
	//if runtime.GOOS == "windows" {
	//	//configFile = os.Getenv("GOPATH") + "\\src\\goFix\\logger\\goFix.log"
	//	configFile = "./logger/goFix.log"
	//} else {
	//	configFile = os.Getenv("GOPATH") + "/goFix/logger/goFix.log"
	//}

	//configFile = os.Getenv("GOPATH") + "/goFix/config/config.yaml"
	configFile = "./config/config.yaml"
	return configFile
}

func initConfig() error {
	conf := NewConfig()
	//configPath := getDefaultConfigPath()
	//yamlFile, err := os.ReadFile(configPath)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//err = yaml.Unmarshal(yamlFile, conf)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("./config/")
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s\n", err))
	}
	err = vp.Unmarshal(conf)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
	rootConfig = conf
	return err
}
