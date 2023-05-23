package kafka

import (
	"github.com/Shopify/sarama"
	"log"
	"time"
)

type ConsumerGroupHandle struct {
	Broker string    `json:"broker"`
	CbM    []topicCb `json:"cbm"` // consumer name - conf mapping
}

func (c ConsumerGroupHandle) Setup(s sarama.ConsumerGroupSession) error {
	log.Printf("Setup Claim:%v\n", s.Claims())
	return nil
}

func (c ConsumerGroupHandle) Cleanup(s sarama.ConsumerGroupSession) error {
	log.Printf("CleanUpClaims:%v\n", s.Claims())
	return nil
}

// ConsumeClaim 消费消息
func (c ConsumerGroupHandle) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		//将消息标记为使用
		s.MarkMessage(msg, "")
		for _, v := range c.CbM {
			if msg.Topic == v.Topic && v.Cb != nil {
				time.Sleep(20 * time.Millisecond)
				go func(value []byte) {
					if err := v.Cb(value); err != nil {
						log.Printf("run callback err: %s\n", err.Error())
					}
				}(msg.Value)
			}
		}
	}
	return nil
}
