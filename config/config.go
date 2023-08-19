package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type Configuration struct {
	DevMode      bool        `yaml:"devMode" mapstructure:"devMode"`
	Service      Service     `yaml:"service" mapstructure:"service"`
	MysqlServer  MysqlServer `yaml:"mysqlServer" mapstructure:"mysqlServer"`
	RedisServer  RedisServer `yaml:"redisServer" mapstructure:"redisServer"`
	EsServer     EsServer    `yaml:"esServer" mapstructure:"esServer"`
	KafkaConfig  KafkaConf   `json:"kafka" yaml:"kafka" mapstructure:"kafka"`
	LoggerConfig LoggerConf  `json:"logger" yaml:"logger" mapstructure:"logger"`
}

type LoggerConf struct {
	LogLevel       int    `json:"logLevel" yaml:"logLevel" mapstructure:"logLevel"`
	LogDir         string `json:"logDir" yaml:"logDir" mapstructure:"logDir"`
	MaxSizeM       int    `json:"maxSizeM" yaml:"maxSizeM" mapstructure:"maxSizeM"`
	RotateMaxAge   int    `json:"rotateMaxAge" yaml:"rotateMaxAge" mapstructure:"rotateMaxAge"`
	FileTimeFormat string `json:"fileTimeFormat" yaml:"fileTimeFormat" mapstructure:"fileTimeFormat"`
}

type EsServer struct {
	EsHost   string `yaml:"esHosts" mapstructure:"esHosts"`
	EsUser   string `yaml:"esUser" mapstructure:"esUser"`
	EsPasswd string `yaml:"esPassword" mapstructure:"esPassword"`
}

type Service struct {
	ServerPort string `yaml:"serverPort" mapstructure:"serverPort"`
}

type MysqlServer struct {
	Address       string `yaml:"address" mapstructure:"address"`
	UserName      string `yaml:"userName" mapstructure:"userName"`
	PassWord      string `yaml:"passWord" mapstructure:"passWord"`
	MysqlMaxDBs   int    `yaml:"mysqlMaxDBs" mapstructure:"mysqlMaxDBs"`
	DefaultDbName string `yaml:"defaultDbName" mapstructure:"defaultDbName"`
}

type RedisServer struct {
	Address       string `yaml:"address" mapstructure:"address"`
	RedisPassword string `yaml:"passWord" mapstructure:"passWord"`
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

// KafkaConf kafka连接配置
type KafkaConf struct {
	Brokers         string `json:"brokers" yaml:"brokers" mapstructure:"brokers"`
	ConsumerGroupID string `json:"consumerGroupID" yaml:"consumerGroupID" mapstructure:"consumerGroupID"`
	GoroutineCount  int    `json:"goroutineCount" yaml:"goroutineCount" mapstructure:"goroutineCount"`
	//Topic           string `json:"topic" yaml:"topic" mapstructure:"topic"`
	Producer ProducerConf `json:"producer" yaml:"producer" mapstructure:"producer"`
	Consumer ConsumerConf `json:"consumer" yaml:"consumer" mapstructure:"consumer"`
}

type ProducerConf struct {
	Topic string `json:"topic" yaml:"topic" mapstructure:"topic"`
}

type ConsumerConf struct {
	Topic string `json:"topic" yaml:"topic" mapstructure:"topic"`
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
