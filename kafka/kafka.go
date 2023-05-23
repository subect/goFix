package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"goFix/config"
	"log"
	"strings"
	"sync"
	"time"
)

var kafkaOpt *KafkaOpt

// KafkaOpt 操作实例
type KafkaOpt struct {
	once         sync.Once
	Conf         config.KafkaConf  `json:"conf"`
	ProducerConf ProducerConfig    `json:"producerConf"`
	ConsumerConf ConsumerConfigure `json:"consumerConf"`

	Producer sarama.AsyncProducer
	Consumer sarama.ConsumerGroup

	Exit chan struct{}
}

func (o *KafkaOpt) Setup(topic string, callBackFunc func(in []byte) error) error {
	o.once.Do(func() {
		o.ProducerConf.Broker = o.Conf.Brokers
		o.ConsumerConf.Broker = o.Conf.Brokers
		o.ConsumerConf.GroupID = o.Conf.ConsumerGroupID

		if err := o.setupProducer(); err != nil {
			log.Fatalf("setup producer failed, err: %v", err)
			return
		}
		if err := o.setupConsumer(); err != nil {
			log.Fatalf("setup consumer failed, err: %v", err)
			return
		}
	})

	if err := o.setConsumerHandler(topic, callBackFunc); err != nil {
		log.Fatalf("set consumer handler failed, err: %v", err)
		return err
	}

	return nil
}

func (o *KafkaOpt) setupProducer() error {
	config := sarama.NewConfig()
	cfg := o.ProducerConf
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = time.Duration(cfg.Frequency) * time.Millisecond
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.MaxMessageBytes = 1000000
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 500 * time.Millisecond
	config.Net.DialTimeout = 30 * time.Second
	config.Net.ReadTimeout = 30 * time.Second
	config.Net.WriteTimeout = 30 * time.Second
	config.Net.KeepAlive = 30 * time.Second
	config.Net.MaxOpenRequests = 5

	producer, err := sarama.NewAsyncProducer(strings.Split(cfg.Broker, ","), config)
	if err != nil {
		return err
	}
	o.Producer = producer
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				log.Printf("producer error: %v", err)
			case res := <-success:
				data, _ := res.Value.Encode()
				log.Printf("producer success send:%v\n !", string(data))
			case <-o.Exit:
				return
			}
		}
	}(producer)
	return nil
}

func (o *KafkaOpt) setupConsumer() error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Return.Errors = false
	config.Consumer.Fetch.Max = 1000000
	config.Consumer.Retry.Backoff = 500 * time.Millisecond
	//config.ClientID = "IPADDR"
	//config.Consumer.Offsets.Initial = sarama.OffsetNewest
	//config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Net.DialTimeout = 30 * time.Second
	config.Net.ReadTimeout = 30 * time.Second
	config.Net.WriteTimeout = 30 * time.Second
	config.Net.KeepAlive = 30 * time.Second
	config.Net.MaxOpenRequests = 5
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = time.Second

	consumer, err := sarama.NewConsumerGroup(strings.Split(o.ConsumerConf.Broker, ","), o.ConsumerConf.GroupID, config)
	if err != nil {
		return err
	}
	o.Consumer = consumer
	go func() {
		for err := range consumer.Errors() {
			log.Printf("consumer error: %v", err)
		}
	}()
	return nil
}

func (o *KafkaOpt) setConsumerHandler(topic string, callBackFunc func(in []byte) error) error {
	if topic == "" || callBackFunc == nil {
		return nil
	}
	o.ConsumerConf.CbM = append(o.ConsumerConf.CbM, topicCb{
		Topic: topic,
		Cb:    callBackFunc,
	})
	return nil
}

func (o *KafkaOpt) Run() error {
	if o.Consumer == nil {
		return nil
	}
	ctx := context.Background()
	handler := ConsumerGroupHandle{
		Broker: o.ConsumerConf.Broker,
		CbM:    o.ConsumerConf.CbM,
	}

	topics := make([]string, len(o.ConsumerConf.CbM))
	for _, v := range o.ConsumerConf.CbM {
		topics = append(topics, v.Topic)
	}

	for {
		if err := o.Consumer.Consume(ctx, topics, handler); err != nil {
			log.Printf("consumer consume error: %v", err)
		}
		if o.Exit != nil {
			return nil
		}
	}
}

// ConsumerConfigure 消费者配置
type ConsumerConfigure struct {
	Broker  string    `json:"broker"`
	GroupID string    `json:"groupID"`
	CbM     []topicCb `json:"cbm"` // consumer name - conf mapping
}

// topic - 消息处理方法映射
type topicCb struct {
	Topic string `json:"topic"`
	Cb    func(body []byte) (err error)
}

// ProducerConfig 生产者发布配置
type ProducerConfig struct {
	Broker string
	// DefaultTopic string
	Frequency  int
	MaxMessage int
	// Exit         chan struct{}
}

func NewKafkaOpt() *KafkaOpt {
	return &KafkaOpt{
		ProducerConf: ProducerConfig{
			Frequency: 500,
		},
		ConsumerConf: ConsumerConfigure{
			CbM: make([]topicCb, 0),
		},
	}
}

func InitKafka(kafka config.KafkaConf, exit chan struct{}) *KafkaOpt {
	opt := NewKafkaOpt()
	opt.Conf = kafka
	opt.Exit = exit
	kafkaOpt = opt
	return opt
}

func GetKafka() *KafkaOpt {
	return kafkaOpt
}

// PublishMessage kafka异步发送消息
func (o *KafkaOpt) PublishMessage(value []byte) error {
	if o == nil || o.Producer == nil {
		return nil
	}
	if len(value) <= 0 {
		log.Printf("get an empty dataJson to kafka\n")
		return nil
	}

	topic := config.Config().KafkaConfig.Producer.Topic
	data := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(string(value)),
	}
	o.Producer.Input() <- data
	return nil
}
