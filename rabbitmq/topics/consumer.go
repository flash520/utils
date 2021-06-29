/**
 * @Author: koulei
 * @Description: TODO
 * @File:  consumer
 * @Version: 1.0.0
 * @Date: 2021/5/11 15:05
 */

package topics

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type consumer struct {
	connect                     *amqp.Connection
	exchangeType                string
	exchangeName                string
	queueName                   string
	durable                     bool
	occurError                  error // 记录初始化过程中的错误
	connErr                     chan *amqp.Error
	routeKey                    string                    //  断线重连，结构体内部使用
	callbackForReceived         func(receivedData string) //   断线重连，结构体内部使用
	offLineReconnectIntervalSec time.Duration
	retryTimes                  int
	callbackOffLine             func(err *amqp.Error) //   断线重连，结构体内部使用
	addr                        string
}

func CreateConsumer(addr, exchangeType, exchangeName, queueName string, durable bool, reconnectInterval, retryTimes int) (*consumer, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Errorf("error: %v\n", err)
		return nil, err
	}

	Consumer := &consumer{
		connect:                     conn,
		exchangeType:                exchangeType,
		exchangeName:                exchangeName,
		queueName:                   queueName,
		durable:                     durable,
		connErr:                     conn.NotifyClose(make(chan *amqp.Error, 1)),
		offLineReconnectIntervalSec: time.Duration(reconnectInterval),
		retryTimes:                  retryTimes,
		addr:                        addr,
	}
	return Consumer, nil
}

func (c *consumer) Received(routeKey string, callbackFunDealSmg func(receivedData string)) {
	defer func() {
		_ = c.connect.Close()
	}()
	// 将回调函数地址赋值给结构体变量，用于掉线重连使用
	c.routeKey = routeKey
	c.callbackForReceived = callbackFunDealSmg

	blocking := make(chan bool)

	go func(key string) {

		ch, err := c.connect.Channel()
		c.occurError = errorDeal(err)
		defer func() {
			_ = ch.Close()
		}()

		// 声明exchange交换机
		err = ch.ExchangeDeclare(
			c.exchangeName, //exchange name
			c.exchangeType, //exchange kind
			c.durable,      //数据是否持久化
			!c.durable,     //所有连接断开时，交换机是否删除
			false,
			false,
			nil,
		)
		// 声明队列
		queue, err := ch.QueueDeclare(
			c.queueName,
			c.durable,
			!c.durable,
			false,
			false,
			nil,
		)
		c.occurError = errorDeal(err)

		//队列绑定
		err = ch.QueueBind(
			queue.Name,
			key, //  Topics 模式,生产者会将消息投递至交换机的route_key， 消费者匹配不同的key获取消息、处理
			c.exchangeName,
			false,
			nil,
		)
		c.occurError = errorDeal(err)

		msgs, err := ch.Consume(
			queue.Name, // 队列名称
			"",         //  消费者标记，请确保在一个消息频道唯一
			true,       //是否自动响应确认，这里设置为false，手动确认
			false,      //是否私有队列，false标识允许多个 consumer 向该队列投递消息，true 表示独占
			false,      //RabbitMQ不支持noLocal标志。
			false,      // 队列如果已经在服务器声明，设置为 true ，否则设置为 false；
			nil,
		)
		c.occurError = errorDeal(err)

		for msg := range msgs {
			// 通过回调处理消息
			callbackFunDealSmg(string(msg.Body))
		}

	}(routeKey)

	<-blocking

}

// OnConnectionError 消费者端，掉线重连失败后的错误回调
func (c *consumer) OnConnectionError(callbackOfflineErr func(err *amqp.Error)) {
	c.callbackOffLine = callbackOfflineErr
	go func() {
		select {
		case err := <-c.connErr:
			var i = 1
			for i = 1; i <= c.retryTimes; i++ {
				// 自动重连机制
				time.Sleep(c.offLineReconnectIntervalSec * time.Second)
				conn, err := CreateConsumer(
					c.addr,
					c.exchangeType,
					c.exchangeName,
					c.queueName,
					c.durable,
					int(c.offLineReconnectIntervalSec),
					c.retryTimes,
				)
				if err != nil {
					continue
				} else {
					go func() {
						c.connErr = conn.connect.NotifyClose(make(chan *amqp.Error, 1))
						go conn.OnConnectionError(c.callbackOffLine)
						conn.Received(c.routeKey, c.callbackForReceived)
					}()
					break
				}
			}
			if i > c.retryTimes {
				callbackOfflineErr(err)
			}
		}
	}()
}
