package kafka

import (
	"strings"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type producer struct {
	brokers       []string
	syncConfig    *sarama.Config
	asyncConfig   *sarama.Config
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
}

// NewKafkaAsyncProducer 创建一个新的异步 kafka Producer 实例
func NewKafkaAsyncProducer(urls string, returnValue bool) (*producer, error) {
	var err error
	brokers := strings.Split(urls, ",")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	config.Producer.Return.Successes = returnValue
	config.Producer.Return.Errors = returnValue

	p, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range p.Errors() {
			log.Println("Failed to write access log entry:", err)

		}
	}()

	return &producer{
		brokers:       brokers,
		asyncConfig:   config,
		asyncProducer: p,
	}, nil
}

// NewKafkaSyncProducer 创建一个新的 kafka Producer 实例
func NewKafkaSyncProducer(urls string, returnValue bool) (*producer, error) {
	var err error
	brokers := strings.Split(urls, ",")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = returnValue
	config.Producer.Return.Errors = returnValue

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &producer{
		brokers:      brokers,
		syncConfig:   config,
		syncProducer: p,
	}, nil
}

// connectAsyncProducer 建立 AsyncProducer 建接
func (p *producer) connectAsyncProducer() error {
	var err error
	if p.asyncProducer == nil {
		p.asyncProducer, err = sarama.NewAsyncProducer(p.brokers, p.asyncConfig)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// connectSyncProducer 建立 syncProducer 连接
func (p *producer) connectSyncProducer() error {
	var err error
	if p.syncProducer == nil {
		p.syncProducer, err = sarama.NewSyncProducer(p.brokers, p.syncConfig)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

// LoggerHandler 异步 Producer 如果启用了 return.Successes 参数，则必须调用本方法
func (p *producer) LoggerHandler(handler func(logType string, data interface{})) {
	for {
		select {
		case s := <-p.asyncProducer.Successes():
			handler("SUCCESS", s)
		case e := <-p.asyncProducer.Errors():
			handler("ERROR", e)
		}
	}
}

// Input 异步发送数据，同步是 chan 通道
func (p *producer) Input(topic string, value string) {
	var err error
	if p.asyncProducer == nil {
		err = p.connectAsyncProducer()
		if err != nil {
			log.Error(err)
			return
		}
	}
	p.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Key:       nil,
		Value:     sarama.StringEncoder(value),
		Headers:   nil,
		Metadata:  nil,
		Timestamp: time.Time{},
	}
}

// Send 以同步模式发送消息
func (p *producer) Send(topic string, body string) (int32, int64, error) {
	var err error
	if p.syncProducer == nil {
		err = p.connectSyncProducer()
		if err != nil {
			return 0, 0, err
		}
	}

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(body)

	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return 0, 0, err
	}

	return partition, offset, nil
}
