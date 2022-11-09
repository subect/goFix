package router

import (
	"go.uber.org/zap"
	"goFix/config"
)

var basicLog *zap.SugaredLogger
var conf *config.Configuration

func Init() {
	basicLog = conf.Logger()
}
