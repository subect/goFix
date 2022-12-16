package model

import (
	"github.com/olivere/elastic/v7"
	"goFix/config"
	"log"
	"os"
	"sync"
)

var (
	once     sync.Once
	esClient *elastic.Client
)

func GetEsOpt() *elastic.Client {
	once.Do(func() {
		initEs()
	})
	return esClient
}

func initEs() {
	hosts := config.Config().EsServer.EsHost
	log.Println(hosts)
	client, err := elastic.NewClient(
		elastic.SetURL(config.Config().EsServer.EsHost), //用来设置ES服务地址
		elastic.SetSniff(false),                         //允许指定弹性是否应该定期检查集群（默认为true）
		//elastic.SetBasicAuth(config.Config().EsServer.EsUser, config.Config().EsServer.EsPasswd), //基于http base auth 验证机制的账号密码
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)), //设置错误日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),          //设置info日志输出
	)
	if err != nil {
		log.Fatalln("Failed to create elastic client")
	}
	esClient = client
}
