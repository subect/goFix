package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goFix/api"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 在 /debug/pprof/ 路由下注册性能分析处理器
	router.GET("/debug/pprof/", gin.WrapF(http.HandlerFunc(pprof.Index)))
	router.GET("/debug/pprof/profile", gin.WrapF(http.HandlerFunc(pprof.Profile)))
	router.GET("/debug/pprof/heap", gin.WrapF(http.HandlerFunc(pprof.Handler("heap").ServeHTTP)))
	router.GET("/debug/pprof/block", gin.WrapF(http.HandlerFunc(pprof.Handler("block").ServeHTTP)))
	router.GET("/debug/pprof/goroutine", gin.WrapF(http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP)))
	router.GET("/debug/pprof/threadcreate", gin.WrapF(http.HandlerFunc(pprof.Handler("threadcreate").ServeHTTP)))
	router.GET("/debug/pprof/cmdline", gin.WrapF(http.HandlerFunc(pprof.Cmdline)))
	router.GET("/debug/pprof/symbol", gin.WrapF(http.HandlerFunc(pprof.Symbol)))
	router.POST("/debug/pprof/symbol", gin.WrapF(http.HandlerFunc(pprof.Symbol)))

	router.Use(Mid())

	router.GET("/", Hello)
	router.POST("/", Hello)
	router.GET("/trainDb", api.Mysqltd)
	router.POST("/trainRedis", api.Redistd)
	router.GET("/receivePb", api.ReceivePb)
	router.GET("/sedPb", api.SendPb)
	router.POST("/esTrain", api.EsTd)
	router.POST("/validator", api.Validator)
	router.GET("/pinyin", api.Pinyin)
	router.POST("/filerKeyWords", api.FilerKeyWords)
	router.GET("viper", api.Viper)
	router.POST("/sendMsg", api.SendMsg) //kafka 发送消息

	return router
}

func Hello(c *gin.Context) {
	if c.Request.Method == "GET" {
		basicLog.Debugln("go into hello!")
		c.JSON(200, "Success")
	}
	if c.Request.Method == "POST" {
		bs, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		basicLog.Debugf("resp:%+v", string(bs))
		c.JSON(200, string(bs))
	}
	//c.JSON(200, "Good")
	return
}

func Mid() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
