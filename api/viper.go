package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"goFix/config"
)

var configuration config.Configuration

func Viper(c *gin.Context) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("./config/")
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s\n", err))
	}
	err = vp.Unmarshal(&configuration)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
	basicLog.Debugf("configuration:%+v", configuration)
	c.JSON(200, configuration)
}
