package main

import (
	"fmt"
	"goFix/api"
	"goFix/config"
	"goFix/dao"
	router2 "goFix/router"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	conf := config.InitConfig()
	basicLog := conf.Logger()
	dao.InitDb()
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
	os.Exit(0)

}
