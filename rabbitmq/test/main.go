package main

import (
	"gitee.com/flash520/utils/rabbitmq/simple"
	"gitee.com/flash520/utils/rabbitmq/test/handler"
	"gitee.com/flash520/utils/response"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var (
	producer, err = simple.NewSimpleProducer(
		"amqp://guest:guest@localhost:5672/", "test", "test",
	)
	consumer, _ = simple.NewSimpleConsumer(
		"amqp://guest:guest@127.0.0.1:5672/", "test", "test", 2,
	)
)

func init() {
	consumer.Declare()
	go receive()
}

func main() {

	r := gin.Default()
	pprof.Register(r)
	r.GET("/send/:id", send)

	_ = r.Run(":80")
}

func send(c *gin.Context) {
	msg := c.Param("id")
	err = producer.Send(msg, "", "")
	if err != nil {
		response.Fail(c, 0, err, nil)
		return
	}
	response.Success(c, 1, "success", nil)
}

func receive() {
	h := handler.Handler{}
	consumer.Received(false, h.Receive)
}
