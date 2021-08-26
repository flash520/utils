package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type Pack struct {
	FilenameLength int64
	Length         int64
	Filename       []byte
	Data           []byte
}

func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		log.Error("net dail failed, err: ", err.Error())
		return
	}

	// sendfile(err, conn)
	Input(conn)
}

func Input(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	_, err := conn.Write([]byte(input))
	if err != nil {
		log.Error(err)
	} else {
		log.Info("success")
	}
}
func sendfile(err error, conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error("conn close: ", err.Error())
			return
		} else {
			log.Info("conn close success")
		}
	}()
	//frame, err := os.Open("/Users/koulei/唐伯虎点秋香.mp4")
	f, err := os.Open("/Users/koulei/msg.json")
	if err != nil {
		log.Error("read file failed, err: ", err.Error())
		return
	}
	defer func() { _ = f.Close() }()

	info, _ := f.Stat()

	pack := &Pack{
		FilenameLength: int64(len(info.Name())),
		Length:         int64(0),
		Filename:       []byte(info.Name()),
		//Data:           []byte("abcd"),
	}
	pack.Length = 8 + 8 + pack.FilenameLength + info.Size()

	log.Infof("文件名长度: %v, 文件名: %s,  数据长度: %d\n", pack.FilenameLength, pack.Filename, pack.Length)
	log.Infof("连接端口: %v\n", conn.LocalAddr())

	writer := new(bytes.Buffer)
	buf := make([]byte, 4096)
	err = pack.Pack(writer)
	if err != nil {
		log.Error("pack failed,err: ", err.Error())
		return
	}
	_, err = conn.Write(writer.Bytes())
	if err != nil {
		log.Errorf("向网络写入数据失败: %v\n", err)
		return
	}

	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {

				break
			} else {
				log.Error("read stream failed, err: ", err)
				return
			}
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Error("sendfile err: ", err)
			return
		}
	}

	// 在规定的时间内如果没有收到服务器返回的结果或者返回值不为 receive ok，则退出并将退出错误码设为 -1
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	n, _ := conn.Read(buf)
	if string(buf[:n]) == "Transaction Success" {
		log.Info("file send success")
		return
	}
	log.Error("failed")
	return
}

func (p *Pack) Pack(writer io.Writer) error {
	var err error
	err = binary.Write(writer, binary.BigEndian, &p.FilenameLength)
	err = binary.Write(writer, binary.BigEndian, &p.Length)
	err = binary.Write(writer, binary.BigEndian, &p.Filename)
	err = binary.Write(writer, binary.BigEndian, &p.Data)
	return err
}
