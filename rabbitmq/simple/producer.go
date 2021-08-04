package simple

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Producer struct {
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
func NewSimpleProducer(url, exchangeName, queueName string) (*Producer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Error("connect: ", err.Error())
		return nil, err
	}

	producer := &Producer{
		connect:           conn,
		exchangeType:      exchangeType,
		exchangeName:      exchangeName,
		queueName:         queueName,
		durable:           true,
		conErr:            conn.NotifyClose(make(chan *amqp.Error)),
		reconnectInterval: reconnectInterval,
		reconnectCount:    0,
		mqURL:             url,
	}

	// 初始化 channel
	producer.newChannel()

	// 开启协程监听连接状态，如果断开，则尝试重新连接并输出日志
	go producer.OnConnectionErrorReConnection()
	return producer, nil
}

func (p *Producer) newChannel() {
	var err error

	// 如果通道为空则建立通道，后续复用该通道
	p.channel, err = p.connect.Channel()
	if err != nil {
		log.Error(err.Error())
		return
	}
}

// Send 发送数据，channel 复用
func (p *Producer) Send(data, routeKey, expire, objectType string) error {
	var err error
	if routeKey == "" {
		routeKey = "normal"
	}

	// 发送数据到交换机
	err = p.channel.Publish(
		p.exchangeName,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
			Expiration:  expire,
			Headers: map[string]interface{}{
				// 声明发送的数据的原对象类型名称，便于消费者实例化数据时使用
				// 该声明可以为空
				"object": objectType,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// Close 手动关闭连接
func (p *Producer) Close() {
	_ = p.channel.Close()
	_ = p.connect.Close()
}

// OnConnectionErrorReConnection 监听连接错误，自动重连
func (p *Producer) OnConnectionErrorReConnection() {
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
				p.newChannel()
				go p.OnConnectionErrorReConnection()
				log.Infof("RabbitMQ Producer 重连 %d 次成功, 实例类型: %T, 实例地址: %p \n", p.reconnectCount, p, p)
				break
			}
		}()
	}
}
