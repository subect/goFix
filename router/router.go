package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goFix/api"
	"io/ioutil"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(Mid())

	router.GET("/", Hello)
	router.POST("/", Hello)
	router.GET("/trainDb", api.Mysqltd)

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
