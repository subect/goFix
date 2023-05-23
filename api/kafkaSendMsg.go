package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"goFix/kafka"
)

func SendMsg(c *gin.Context) {
	msg := "这是一条kafka消息"
	kafkaClient := kafka.GetKafka()
	if kafkaClient == nil {
		basicLog.Errorln("kafka client is nil")
		return
	}
	msgByte, err := json.Marshal(msg)
	if err != nil {
		basicLog.Errorf("json marshal err:%s", err.Error())
		return
	}
	err = kafkaClient.PublishMessage(msgByte)
	if err != nil {
		basicLog.Errorf("publish message err:%s", err.Error())
		return
	}
	basicLog.Debugln("publish message success")
}
