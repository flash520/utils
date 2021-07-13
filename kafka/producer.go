package kafka

import (
	"strings"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type producer struct {
	brokers []string
	config  *sarama.Config
	client  sarama.SyncProducer
}

// NewKafkaProducer 创建一个新的 kafka Producer 实例
func NewKafkaProducer(urls string) (*producer, error) {
	var err error
	brokers := strings.Split(urls, ",")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &producer{
		brokers: brokers,
		config:  config,
		client:  p,
	}, nil
}

// connect 建立 producer 连接
func (p *producer) connect() error {
	var err error
	if p.client == nil {
		p.client, err = sarama.NewSyncProducer(p.brokers, p.config)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

// Send 以同步模式发送消息
func (p *producer) Send(topic string, body string) (int32, int64, error) {
	var err error
	if p.client == nil {
		err = p.connect()
		if err != nil {
			return 0, 0, err
		}
	}

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(body)

	partition, offset, err := p.client.SendMessage(msg)
	if err != nil {
		return 0, 0, err
	}

	return partition, offset, nil
}
