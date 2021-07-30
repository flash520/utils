/**
 * @Author: koulei
 * @Description:
 * @File: sumx
 * @Version: 1.0.0
 * @Date: 2021/7/21 12:51
 */

package main

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	sumx "github.com/xtaci/smux"
)

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Errorf("服务监听失败, err: %s\n", err)
	}

	log.Infof("MultiComm 多路复用服务器 Starting...")
	for i := 0; i < 100; i++ {
		conn, err := l.Accept()
		if err != nil {
			log.Errorf("服务接收请求连接失败, err: %s\n", err)
			continue
		}

		go sessionHandler(conn)
	}
	log.Info()
}

func sessionHandler(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	session, err := sumx.Server(conn, nil)
	if err != nil {
		log.Errorf("创建 Session 连接失败, err: %s\n", err)
		return
	}

	log.Infof("收到客户端连接，创建新会话，对端地址：%s", session.RemoteAddr().String())
	for !session.IsClosed() {
		s, err := session.AcceptStream()
		if err != nil {
			log.Errorf("创建 Stream 连接失败, err: %s\n", err)
			break
		}

		go streamHandler(s)
	}
	log.Errorf("客户端连接断开，销毁会话，对端地址：%s", session.RemoteAddr().String())
}

func streamHandler(s *sumx.Stream) {
	var (
		n   int
		err error
	)

	defer func() { _ = s.Close() }()

	_ = s.SetReadDeadline(time.Now().Add(time.Second * 3))
	buf := make([]byte, 4096)
	var d int
	for {

		n, err = s.Read(buf)
		d += n
		s.SetReadDeadline(time.Now().Add(time.Second * 3))
		if err != nil {
			log.Errorf("streamID: %v 读取数据失败, err: %s", s.ID(), err)
			return
		}
		log.Infof("streamID: %v 成功读取数据: %d", s.ID(), d)
	}
}
