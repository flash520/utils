package main

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"gitee.com/flash520/utils/kafka"
	"github.com/Shopify/sarama"
)

func main() {
	url := "localhost:9092"
	c, err := kafka.NewKafkaConsumer(url, "topic")
	if err != nil {
		fmt.Print(err)
		return
	}

	go c.Receive("topic", handler, map[int]int64{0: 10000000})

	//	Producer
	p, err := kafka.NewKafkaSyncProducer(url)
	//p, err := kafka.NewKafkaAsyncProducer(url, false)
	if err != nil {
		fmt.Print(err)
		return
	}
	//p.LogHandler(nil)

	started := time.Now()
	for i := 0; i < 100000; i++ {
		go func(i int) {
			_, _, err = p.Send("topic", "body: "+strconv.Itoa(i))
			//p.Input("topic", "body: "+strconv.Itoa(i))
			if err != nil {
				log.Error(err.Error())
				return
			}
			//fmt.Println("success")
		}(i)
	}

	fmt.Println("time: ", time.Since(started))
	select {}
}

func handler(msg *sarama.ConsumerMessage) {
	fmt.Printf("\rTopic: %s Partition: %d Offset: %d Key: %v Value: %v", msg.Topic, msg.Partition, msg.Offset, msg.Key, string(msg.Value))
}

func logger(logType string, data interface{}) {
	switch logType {
	case "SUCCESS":
		fmt.Println(data.(*sarama.ProducerMessage))
	case "ERROR":
		fmt.Println(data.(*sarama.ProducerError))
	}
}
