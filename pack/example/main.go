package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type pack struct {
	Version []byte
	Length  int64
	Data    []byte
}

func main() {
	var err error
	Pack := pack{
		Version: []byte("v1"),
		Length:  int64(4),
		Data:    []byte("data"),
	}

	writer := new(bytes.Buffer)
	err = Pack.Pack(writer)
	if err != nil {
		fmt.Println("打包错误,err: ", err)
		return
	}

	Result := pack{}
	err = Result.Unpack(writer)
	if err != nil {
		fmt.Println("解包错误,err: ", err)
		return
	}
	Result.Data = make([]byte, Result.Length)
	err = binary.Read(writer, binary.BigEndian, &Result.Data)
	Result.String()

	fmt.Println()
}

func (p *pack) Pack(buf io.Writer) error {
	var err error
	err = binary.Write(buf, binary.BigEndian, &p.Version)
	err = binary.Write(buf, binary.BigEndian, &p.Length)
	err = binary.Write(buf, binary.BigEndian, &p.Data)
	return err
}

func (p *pack) Unpack(r io.Reader) error {
	var err error
	p.Version = make([]byte, 2)
	err = binary.Read(r, binary.BigEndian, &p.Version)
	err = binary.Read(r, binary.BigEndian, &p.Length)
	//p.Data = make([]byte, p.Length)
	//err = binary.Read(r, binary.BigEndian, &p.Data)
	return err
}

func (p *pack) String() {
	fmt.Printf("Version: %s, Length: %d, Data: %s\n",
		p.Version,
		p.Length,
		string(p.Data),
	)
}
