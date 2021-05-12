/**
 * @Author: koulei
 * @Description: TODO
 * @File:  producer
 * @Version: 1.0.0
 * @Date: 2021/5/11 12:58
 */

package topics

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type producer struct {
	connect      *amqp.Connection
	exchangeType string
	exchangeName string
	queueName    string
	durable      bool
	occurError   error
}

func CreateProducer(Addr, exchangeType, exchangeName, queueName string, durable bool) (*producer, error) {
	conn, err := amqp.Dial(Addr)
	if err != nil {
		log.Errorf("error: %v\n", err)
		return nil, err
	}

	Producer := &producer{
		connect:      conn,
		exchangeType: exchangeType,
		exchangeName: exchangeName,
		queueName:    queueName,
		durable:      durable,
	}
	return Producer, nil
}

func (p *producer) Send(routeKey string, data string) bool {
	ch, err := p.connect.Channel()
	p.occurError = errorDeal(err)
	defer func() { _ = ch.Close() }()

	err = ch.ExchangeDeclare(
		p.exchangeType,
		p.exchangeName,
		p.durable,
		!p.durable,
		false,
		false,
		nil,
	)
	p.occurError = errorDeal(err)

	err = ch.Publish(
		p.exchangeName,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})

	if p.occurError != nil {
		return false
	} else {
		return true
	}
}

func (p *producer) Close() {
	_ = p.connect.Close()
}

func errorDeal(err error) error {
	if err != nil {
		log.Errorf("error: %v\n", err.Error())
	}
	return err
}
