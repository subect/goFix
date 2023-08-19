package model

import (
	"go.uber.org/zap"
	"goFix/config"
)

var conf *config.Configuration
var basicLog *zap.SugaredLogger

func Init() {
	conf = config.InitConfig()
	basicLog = conf.Logger()
	NewDbPool()
	Initialize()
	initMysqlDB()

}
