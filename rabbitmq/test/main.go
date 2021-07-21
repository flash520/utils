package main

import (
	"strconv"

	"gitee.com/flash520/utils/rabbitmq/simple"
	"gitee.com/flash520/utils/rabbitmq/test/handler"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Consumer
	c, err := simple.NewSimpleConsumer(
		"localhost", "guest", "guest", "mm-im", "direct",
		"ex1", "qu1", 2)
	if err != nil {
		log.Error(err.Error())
		return
	}

	h := &handler.Handler{}

	go c.Received("failover", true, h.Receive)

	// Producer
	p, err := simple.NewSimpleProducer(
		"localhost", "guest", "guest", "mm-im", "direct",
		"ex1", "qu1", 2)
	if err != nil {
		log.Error(err.Error())
		return
	}

	go func() {

		for i := 0; i < 1000; i++ {
			err = p.Send("failover", "message body: "+strconv.Itoa(i))
			if err != nil {
				log.Error(err.Error())
			}
		}
	}()

	select {}
}
