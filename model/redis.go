package model

import (
	"context"
	"github.com/go-redis/redis/v8"
	"goFix/config"
	"time"
)

// InitRedis 单机
func InitRedis() *redis.Client {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:        config.Config().RedisServer.Address,
		Password:    config.Config().RedisServer.RedisPassword,
		DB:          0,
		IdleTimeout: 300, // 默认Idle超时时间
		PoolSize:    100,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		basicLog.Errorf("Connect Failed! Err: %v\n", err.Error())
		return nil
	}
	basicLog.Debugf("Connection Success \nPING => %v\n", result)
	return RedisClient

	//TODO 哨兵模式
}
