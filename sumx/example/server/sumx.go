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

	log "github.com/sirupsen/logrus"
	sumx "github.com/xtaci/smux"
)

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Errorf("服务监听失败, err: %s\n", err)
	}

	log.Infof("MultiComm 多路复用服务器 Starting...")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Errorf("服务接收请求连接失败, err: %s\n", err)
			continue
		}

		go sessionHandler(conn)
	}
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

	buf := make([]byte, 1024)
	n, err = s.Read(buf)
	if err != nil {
		log.Errorf("streamID: %v 读取数据失败, err: %s", s.ID(), err)
		return
	}
	log.Infof("streamID: %v 成功读取数据: %s", s.ID(), string(buf[:n]))
}
