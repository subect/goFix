package main

import (
	"fmt"
	"goFix/api"
	"goFix/config"

	"goFix/model"
	router2 "goFix/router"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	//exit := make(chan struct{})
	var wg sync.WaitGroup
	conf := config.InitConfig()
	basicLog := conf.Logger()
	model.Init()
	api.Init()
	router2.Init()
	router := router2.InitRouter()
	basicLog.Debugf("conf.Service.ServerPort:%s", conf.Service.ServerPort)
	port := ":" + conf.Service.ServerPort
	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//kafkaOpt := kafka.InitKafka(conf.KafkaConfig, exit)
	//err := kafkaOpt.Setup(conf.KafkaConfig.Consumer.Topic, handler.ConsumeFunc)
	//if err != nil {
	//	fmt.Printf("init kafka text opt err:%s\n", err.Error())
	//	return
	//}
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	if err := kafkaOpt.Run(); err != nil {
	//		fmt.Printf("run kafka server err:%s\n", err.Error())
	//		return
	//	}
	//}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				fmt.Println("closing gin server!")
			} else {
				panic(err)
			}
		}
	}()

	//fmt.Printf("gin run on port:%v ", port)
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	fmt.Println("use c-c to exit")
	<-sigChan
	// 在程序退出之前关闭数据库连接
	model.CloseDb()
	os.Exit(0)
}
