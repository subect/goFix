package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
)

func Pinyin(c *gin.Context) {
	text := "我是刘德华"
	a := pinyin.NewArgs()
	a.Style = pinyin.Tone2
	c.JSON(200, pinyin.Pinyin(text, a))
}
