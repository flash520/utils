package simple

import (
	"time"

	"gitee.com/flash520/utils/randomstring"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Consumer struct {
	connect           *amqp.Connection
	channel           *amqp.Channel
	exchangeType      string                                 // 交换机类型
	exchangeName      string                                 // 交换机名称
	queueName         string                                 // 队列名称
	durable           bool                                   // 是否持久化
	conErr            chan *amqp.Error                       // 错误通道
	reconnectInterval int                                    // 断线重连间隔(秒)
	reconnectCount    int                                    // 断线重连次数
	mqURL             string                                 // 服务器地址
	routeKey          string                                 // 路由标识
	callbackHandler   func(receive, objectType string) error // 用于断线重连(消息处理器)
	chanNumber        int                                    // 并发消息者数量
}

const (
	reconnectInterval = 5
	exchangeType      = "direct"
	routeKey          = "normal"
	deadRouteKey      = "dead"
	deadSuffix        = ".dlx"
)

// NewSimpleConsumer 创建一个新的 Consumer 实例
// url 示例: amqp://user:password@addr:5672/
func NewSimpleConsumer(url, exchangeName, queueName string, chanNumber int) (*Consumer, error) {
	var err error
	conn, err := amqp.DialConfig(url, amqp.Config{
		Properties: map[string]interface{}{
			"connection_name": "Go#Consumer#Simple:" + queueName,
		},
	})
	if err != nil {
		log.Error("connect: ", err.Error())
		return nil, err
	}

	consumer := &Consumer{
		connect:           conn,
		exchangeType:      exchangeType,
		exchangeName:      exchangeName,
		queueName:         queueName,
		durable:           true,
		conErr:            conn.NotifyClose(make(chan *amqp.Error)),
		reconnectInterval: reconnectInterval,
		reconnectCount:    0,
		mqURL:             url,
		routeKey:          routeKey,
		chanNumber:        chanNumber,
	}

	// 开启协程监听连接状态，如果断开，则尝试重新连接并输出日志
	go consumer.OnConnectionErrorReConnection()

	return consumer, nil
}

func (c *Consumer) Destroy() {
	err := c.channel.Close()
	if err != nil {
		return
	}
	err = c.connect.Close()
	if err != nil {
		return
	}
}

// OnConnectionErrorReConnection 连接错误，自动重连
func (c *Consumer) OnConnectionErrorReConnection() {
	select {
	case e := <-c.conErr:
		log.Errorf("RabbitMQ Consumer 连接错误: %s\n", e)
		go func() {
			for {
				time.Sleep(time.Duration(c.reconnectInterval) * time.Second)
				c.reconnectCount++
				conn, err := amqp.Dial(c.mqURL)
				if err != nil {
					log.Errorf("RabbitMQ Consumer 重连 %d 次失败: %s\n", c.reconnectCount, err)
					continue
				}
				c.connect = conn
				c.conErr = c.connect.NotifyClose(make(chan *amqp.Error))
				go c.OnConnectionErrorReConnection()
				go c.Received(false, c.callbackHandler)
				log.Infof("RabbitMQ Consumer 重连 %d 次成功, 实例类型: %T, 实例地址: %p \n", c.reconnectCount, c, c)
				break
			}
		}()
	}
}

// Declare 队列绑定
func (c *Consumer) Declare() {
	var err error
	c.channel, err = c.connect.Channel()
	if err != nil {
		log.Error("channel err: ", err)
		return
	}
	defer func() { _ = c.channel.Close() }()

	// 声明死信交换机
	err = c.channel.ExchangeDeclare(
		c.exchangeName+deadSuffix,
		c.exchangeType,
		c.durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error(err.Error())
	}

	// 声明死信队列
	_, err = c.channel.QueueDeclare(
		c.queueName+deadSuffix,
		c.durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error(err.Error())
	}

	// 绑定死信交换机和队列
	err = c.channel.QueueBind(
		c.queueName+deadSuffix,
		deadRouteKey,
		c.exchangeName+deadSuffix,
		false,
		nil,
	)
	if err != nil {
		log.Error(err.Error())
	}

	// 声明常规交换机
	err = c.channel.ExchangeDeclare(
		c.exchangeName,
		exchangeType,
		c.durable,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		log.Errorf(err.Error())
	}

	// 声明常规队列
	_, err = c.channel.QueueDeclare(
		c.queueName,
		c.durable,
		false,
		false,
		false,
		amqp.Table{
			// 添加死信队列配置
			// 在常队列中声明 nack/reject 消息转发目的交换机名称
			"x-dead-letter-exchange":    c.exchangeName + ".dlx",
			"x-dead-letter-routing-key": "dead",
		},
	)
	if err != nil {
		log.Errorf(err.Error())
	}

	// 绑定常规队列
	err = c.channel.QueueBind(
		c.queueName,
		routeKey,
		c.exchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Errorf(err.Error())
	}
}

// Received 通过回调函数处理接收到的消息
func (c *Consumer) Received(autoAck bool, handler func(receiveData, objectType string) error) {
	defer func() { c.Destroy() }()

	c.callbackHandler = handler

	blocking := make(chan bool)

	for i := 1; i <= c.chanNumber; i++ {
		go func(chanNumber int) {
			channel, err := c.connect.Channel()
			if err != nil {
				log.Error("channel err: ", err)
				return
			}
			// 自动生成 ConsumerID 作为消费者标记，并确保在一个消息频道唯一
			messages, err := channel.Consume(
				c.queueName, // 队列名称
				randomstring.RandStringBytesMaskImprSrcUnsafe(12), // 消费者标记，请确保在一个消息频道唯一
				autoAck, // 是否自动响应确认，这里设置为false，手动确认
				false,   // 是否私有队列，false标识允许多个 Consumer 向该队列投递消息，true 表示独占
				false,   // RabbitMQ不支持noLocal标志。
				false,   // 队列如果已经在服务器声明，设置为 true ，否则设置为 false；
				nil,
			)
			if err != nil {
				log.Errorf(err.Error())
			}

			for msg := range messages {
				// 通过回调处理消息
				if autoAck {
					_ = handler(string(msg.Body), msg.Type)
				} else {
					err = handler(string(msg.Body), msg.Type)
					if err != nil {
						// 启用死信交换机后，此处 requeue 一定要设为 false
						// 当消息被 Nack 和 Reject 否认确认时，该消息将发送到 死信队列
						// 死信消息可以发送到队列声明时定义的死信交换机或者死信路由对应的队列上
						_ = msg.Nack(false, false)
					} else {
						_ = msg.Ack(false)
					}
				}
			}
			log.Errorf("RabbitMQ Consumer %d quit", chanNumber)
		}(i)
	}
	<-blocking
}
