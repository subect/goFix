package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"goFix/config"
)

var conf *config.Configuration
var basicLog *zap.SugaredLogger
var redisClient *redis.Client
var esClient *elastic.Client

func Init() {
	conf = config.InitConfig()
	basicLog = conf.Logger()
	//redisClient = model.InitRedis()
	//esClient = model.GetEsOpt()
}
