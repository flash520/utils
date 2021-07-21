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

const MaxFilesize = 2 * 1024 * 1024 * 1024

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Receiver is Starting...")
	for {
		time.Sleep(time.Second * 10)
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go receiver(conn)
	}
}

func receiver(conn net.Conn) {
	var err error
	var n int
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error("conn close failed, ", err.Error())
			return
		} else {
			log.Info("conn close success")
			return
		}
	}()

	var i int

	buf := make([]byte, 4096)
	var length int64
	n, err = conn.Read(buf)
	if err != nil {
		log.Error(err.Error())
		return
	}
	i++

	pack := &Pack{}
	reader := bytes.NewReader(buf[:n])

	headerSize, err := pack.Unpack(reader)
	if err != nil {
		log.Error("包头解析错误,err: ", err.Error())
		log.Warnf("包异常, data: %v, count: %d\n", len(buf[:n]), i)
		_, _ = conn.Write([]byte("包头解析错误,err: " + err.Error()))
		return
	}
	log.Infof("filenamelength: %d, filename: %v, total: %v, filesize: %d\n",
		pack.FilenameLength, string(pack.Filename), pack.Length,
		pack.Length-16-pack.FilenameLength)

	f, err := os.Create(string(pack.Filename))
	if err != nil {
		log.Error("Create file failed, err: ", err.Error())
		return
	}
	defer func() { _ = f.Close() }()

	_, err = f.Write(buf[headerSize:n])
	if err != nil {
		log.Error("first pack write to file failed, err: ", err)
		return
	}

	defer func(s *Pack) {
		filesize := s.Length - 16 - s.FilenameLength
		if length != s.Length-16-s.FilenameLength {
			log.Errorf("file size: %d, receive size: %d", filesize, length)
		} else {
			log.Infof("file size: %d, receive size: %d", filesize, length)
		}

	}(pack)
	var fileSize int64
	length = int64(len(buf[:n]) - headerSize)
	fileSize = pack.Length - 8 - 8 - pack.FilenameLength
	if len(buf[:n])-headerSize == int(fileSize) {
		log.Info("文件一次性完成了")
		goto Confirm
	}

	for {
		n, err = conn.Read(buf)
		if n != len(buf) && n != 0 {
			log.Warnf("offset: %d\n", n)
		}
		if err != nil {
			if err == io.EOF {
				log.Warnf("流已经EOF,offset: %d\n", n)
				break
			} else {
				log.Error("conn is broken,close: ", err.Error())
				return
			}
		}
		n, _ = f.Write(buf[:n])

		length += int64(n)
		i++
		if length == fileSize {
			log.Info("File Write success")
			break
		}
	}

Confirm:
	_, _ = conn.Write([]byte("Transaction Success"))
}

func (p *Pack) Unpack(reader io.Reader) (int, error) {
	var err error
	err = binary.Read(reader, binary.BigEndian, &p.FilenameLength)
	err = binary.Read(reader, binary.BigEndian, &p.Length)
	p.Filename = make([]byte, p.FilenameLength)
	err = binary.Read(reader, binary.BigEndian, &p.Filename)

	//p.Data = make([]byte, p.Length-8-8-p.FilenameLength)
	//err = binary.Read(reader, binary.BigEndian, &p.Data)

	n := 8 + 8 + int(p.FilenameLength)
	return n, err
}

func Receive(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF { // 由于我们定义的数据包头最开始为两个字节的版本号，所以只有以V开头的数据包才处理
			if len(data) > 16 { // 如果收到的数据>4个字节(2字节版本号+2字节数据包长度)
				length := int64(0)
				binary.Read(bytes.NewReader(data[8:16]), binary.BigEndian, &length) // 读取数据包第3-4字节(int16)=>数据部分长度

				if int(length) <= len(data) { // 如果读取到的数据正文长度+2字节版本号+2字节数据长度不超过读到的数据(实际上就是成功完整的解析出了一个包)
					return int(length), data[:int(length)], nil
				}
			}
		}
		return
	})
	buf := make([]byte, 4096)
	scanner.Buffer(buf, MaxFilesize)

	// 打印接收到的数据包
	for scanner.Scan() {
		scannedPack := new(Pack)
		_, err := scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
		if err != nil {
			log.Error("scanner err: ", err)
			return
		}
		f, _ := os.Create(string(scannedPack.Filename))
		write, err := f.Write(scannedPack.Data)
		if err != nil {
			log.Error("write file err: ", err)
			return
		}
		log.Infof("file touch success, size: %d\n", write)
		//log.Printf("文件长度: %d ,文件名: %s ,文件内容: \n%v\n", scannedPack.Length, string(scannedPack.Filename), string(scannedPack.Data))
	}
	if err := scanner.Err(); err != nil {
		log.Error(err.Error())
	}
}
