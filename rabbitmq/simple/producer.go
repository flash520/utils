package simple

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type producer struct {
	connect           *amqp.Connection
	channel           *amqp.Channel
	exchangeType      string           // 交换机类型
	exchangeName      string           // 交换机名称
	queueName         string           // 队列名称
	durable           bool             // 是否持久化
	conErr            chan *amqp.Error // 错误通道
	reconnectInterval int              // 断线重连间隔(秒)
	reconnectCount    int              // 断线重连次数
	mqURL             string           // 服务器地址
}

// NewSimpleProducer 创建一个新的 Producer 实例
func NewSimpleProducer(url, user, password, vhost, exchangeType, exchangeName, queueName string, reconnectInterval int) (*producer, error) {
	addr := fmt.Sprintf("amqp://%s:%s@%s", user, password, url)
	conn, err := amqp.DialConfig(addr, amqp.Config{Vhost: vhost})
	if err != nil {
		log.Error("connect: ", err.Error())
		return nil, err
	}

	Producer := &producer{
		connect:           conn,
		exchangeType:      exchangeType,
		exchangeName:      exchangeName,
		queueName:         queueName,
		durable:           true,
		conErr:            conn.NotifyClose(make(chan *amqp.Error)),
		reconnectInterval: reconnectInterval,
		reconnectCount:    0,
		mqURL:             addr,
	}

	// 初始化 channel
	for Producer.channel == nil {
		err = Producer.newChannel()
		if err != nil {
			log.Error(err.Error())
			continue
		}

		if Producer.channel != nil {
			break
		}
		time.Sleep(time.Second)
	}

	// 开启协程监听连接状态，如果断开，则尝试重新连接并输出日志
	go Producer.OnConnectionErrorReConnection()
	return Producer, nil
}

func (p *producer) newChannel() error {
	var err error
	// 如果通道为空则建立通道，后续复用该通道
	if p.channel == nil {
		p.channel, err = p.connect.Channel()
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	// 声明交换机
	err = p.channel.ExchangeDeclare(
		p.exchangeName,
		p.exchangeType,
		p.durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// Send 发送数据，channel 复用
func (p *producer) Send(routingKey string, data string) error {
	var err error
	if p.channel == nil {
		err = p.newChannel()
		if err != nil {
			return err
		}
	}
	// 发送数据到交换机
	err = p.channel.Publish(
		p.exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	if err != nil {
		return err
	}

	return nil
}

// Close 手动关闭连接
func (p *producer) Close() {
	_ = p.connect.Close()
}

// OnConnectionErrorReConnection 监听连接错误，自动重连
func (p *producer) OnConnectionErrorReConnection() {
	select {
	case e := <-p.conErr:
		log.Errorf("RabbitMQ Producer 连接错误: %s\n", e)
		go func() {
			for {
				time.Sleep(time.Duration(p.reconnectInterval) * time.Second)
				p.reconnectCount++
				conn, err := amqp.Dial(p.mqURL)
				if err != nil {
					log.Errorf("RabbitMQ Producer 重连 %d 次失败: %s\n", p.reconnectCount, err)
					continue
				}
				p.connect = conn
				p.conErr = p.connect.NotifyClose(make(chan *amqp.Error))
				go p.OnConnectionErrorReConnection()
				log.Infof("RabbitMQ Producer 重连 %d 次成功, 实例类型: %T, 实例地址: %p \n", p.reconnectCount, p, p)
				break
			}
		}()
	}
}
