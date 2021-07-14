package kafka

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type consumer struct {
	consumer sarama.Consumer
	config   *sarama.Config
	brokers  []string
}

// NewKafkaConsumer 创建新的 kafka Consumer 实例
func NewKafkaConsumer(url string, topic string) (*consumer, error) {
	addr := strings.Split(url, ",")
	config := sarama.NewConfig()
	c, err := sarama.NewConsumer(addr, config)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &consumer{
		consumer: c,
		config:   config,
		brokers:  addr,
	}, nil
}

// connect 建立连接
func (c *consumer) connect() error {
	var err error
	if c.consumer == nil {
		c.consumer, err = sarama.NewConsumer(c.brokers, c.config)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// Receive 消息处理
// offset 是一个 map[int]int64, key 为分区号, value 为对应分区下已经消费的 offset 位置
// offset 如果取最新的消息 value 设为 -1，建议设为上次消费的 offset 位置, 消费位置由业务模块自行处理，如：存储到数据库或者 redis
// offset 可以设为 nil ,意味着将全部从新消费
func (c *consumer) Receive(topic string, handler func(msg *sarama.ConsumerMessage), offset map[int]int64) {
	var err error
	if c.consumer == nil {
		c.consumer, err = sarama.NewConsumer(c.brokers, c.config)
		if err != nil {
			log.Error(err)
			return
		}
	}
	// 获取指定 topic 下的全部分区列表
	partitionList, err := c.consumer.Partitions(topic)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("分区列表: %v\n", partitionList)

	var Offset = make(map[int]int64)
	if offset == nil {
	} else {
		Offset = offset
	}

	// 遍历分区列表
	for partition := range partitionList {
		go func(partition int) {
			// 针对每个分区创建一个对应的分区消费者
			if _, ok := Offset[partition]; !ok {
				Offset[partition] = 0
			}
			pc, err := c.consumer.ConsumePartition(topic, int32(partition), Offset[partition])
			if err != nil {
				fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
				return
			}

			// 异步从每个分区消费信息
			go func(pc sarama.PartitionConsumer) {
				defer pc.AsyncClose()
				for msg := range pc.Messages() {
					handler(msg)
				}
			}(pc)
		}(partition)
	}
}
