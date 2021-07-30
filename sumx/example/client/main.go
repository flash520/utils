/**
 * @Author: koulei
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2021/7/21 13:36
 */

package main

import (
	"context"
	"net"
	"os"
	"time"

	pool "github.com/jolestar/go-commons-pool"
	log "github.com/sirupsen/logrus"
	sumx "github.com/xtaci/smux"
)

var (
	commonPool *pool.ObjectPool
	ctx        = context.Background()
)

func init() {
	factory := pool.NewPooledObjectFactorySimple(NewSession)
	commonPool = pool.NewObjectPoolWithDefaultConfig(ctx, factory)
	commonPool.Config.MaxTotal = 10
}

func NewSession(context context.Context) (interface{}, error) {
	log.Info("在连接池中生成一个连接")
	var conn net.Conn
	var err error
	conn, err = net.Dial("tcp", ":9000")
	if err != nil {
		for {
			log.Errorf("在连接池中生成连接失败, err: %s\n", err)
			time.Sleep(time.Second)
			conn, err = net.Dial("tcp", ":9000")
			if err == nil {
				break
			}
		}
	} else {
		log.Infof("连接生成完成")
	}

	config := sumx.DefaultConfig()
	session, err := sumx.Client(conn, config)
	if err != nil {
		log.Errorf("打开会话 SmuxSession 失败,err: %s\n", err)
	}

	return session, err
}

func main() {
	for {
		Input()
	}
}

func Input() {
	object, err := commonPool.BorrowObject(ctx)
	if err != nil {
		log.Errorf("从连接池中获取 Session 失败, err: %s\n", err)
		return
	} else {
		log.Info("\n\n")
		log.Info("从连接池获取 Session 成功\n")
	}

	defer func() {
		err := commonPool.ReturnObject(ctx, object)
		if err != nil {
			log.Errorf("资源释放失败,err: %s\n", err)
			return
		}
		log.Info("资源释放成功")
	}()

	client := object.(*sumx.Session)
	s, err := client.OpenStream()
	if err != nil {
		log.Errorf("打开流失败, err: %s\n", err)
		_ = commonPool.InvalidateObject(ctx, object)
		return
	}

	defer func() {
		err = s.Close()
		if err == nil {
			log.Info("Stream 通道关闭成功:", s.ID())
		}
	}()

	f, err := os.Open("/Users/koulei/唐伯虎点秋香.mp4")
	buf := make([]byte, 4096)
	for {
		n, _ := f.Read(buf)
		s.Write(buf[:n])
		time.Sleep(time.Second * 2)
	}
}
