/**
 * @Author: koulei
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2021/7/23 00:06
 */

package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {

	// ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	do := make(chan struct{})

	var a [10]int
	go func(ctx context.Context) {
		defer func() {
			if e := recover(); e != nil {
				log.Warn("PANIC: ", e)
			}
		}()
		c, cancel := context.WithCancel(ctx)
		defer func() { cancel() }()
		select {
		case <-c.Done():
			log.Info("sub end")
			close(do)
		default:
			for i := 0; i < 20; i++ {
				a[i] = 0
			}
		}
	}(ctx)

	// time.Sleep(time.Second * 10)
	log.Info("main end")
	// panic("err is error")
	<-do
}
