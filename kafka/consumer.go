package kafka

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type consumer struct {
	client sarama.Consumer
	config *sarama.Config
	url    []string
}

// NewKafkaConsumer 创建新的 kafka Consumer 实例
func NewKafkaConsumer(url string, topic string) (*consumer, error) {
	addr := strings.Split(url, ",")
	config := sarama.NewConfig()
	client, err := sarama.NewConsumer(addr, config)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &consumer{
		client: client,
		config: config,
		url:    addr,
	}, nil
}

// connect 建立连接
func (c *consumer) connect() error {
	var err error
	if c.client == nil {
		c.client, err = sarama.NewConsumer(c.url, c.config)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// Receive 消息处理
func (c *consumer) Receive(topic string, handler func(msg *sarama.ConsumerMessage)) {
	var err error
	if c.client == nil {
		c.client, err = sarama.NewConsumer(c.url, c.config)
		if err != nil {
			log.Error(err)
			return
		}
	}
	// 获取指定 topic 下的全部分区列表
	partitionList, err := c.client.Partitions(topic)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("分区列表: %v\n", partitionList)

	// 遍历分区列表
	for partition := range partitionList {
		go func(partition int) {
			// 针对每个分区创建一个对应的分区消费者
			pc, err := c.client.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
				return
			}

			// 异步从每个分区消费信息
			go func(c sarama.PartitionConsumer) {
				defer pc.AsyncClose()
				for msg := range c.Messages() {
					handler(msg)
				}
			}(pc)
		}(partition)
	}
}
