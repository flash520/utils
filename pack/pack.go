package pack

import (
	"encoding/binary"
	"io"
)

type Package struct {
	Version        [2]byte // 协议版本，暂定V1
	Length         int16   // 数据内容长度
	Timestamp      int64   // 时间戳
	HostnameLength int16   // 主机名长度
	Hostname       []byte  // 主机名
	TagLength      int16   // 标签长度
	Tag            []byte  // 标签
	Msg            []byte  // 数据内容
}

func (p *Package) Pack(write io.Writer) error {
	var err error
	err = binary.Write(write, binary.BigEndian, &p.Version)
	err = binary.Write(write, binary.BigEndian, &p.Length)
	err = binary.Write(write, binary.BigEndian, &p.Timestamp)
	err = binary.Write(write, binary.BigEndian, &p.HostnameLength)
	err = binary.Write(write, binary.BigEndian, &p.Hostname)
	err = binary.Write(write, binary.BigEndian, &p.TagLength)
	err = binary.Write(write, binary.BigEndian, &p.Tag)
	err = binary.Write(write, binary.BigEndian, &p.Msg)

	return err
}

func (p *Package) Unpack(reader io.Reader) error {
	var err error
	err = binary.Read(reader, binary.BigEndian, &p.Version)
	err = binary.Read(reader, binary.BigEndian, &p.Length)
	err = binary.Read(reader, binary.BigEndian, &p.Timestamp)
	err = binary.Read(reader, binary.BigEndian, &p.HostnameLength)
	p.Hostname = make([]byte, p.HostnameLength)
	err = binary.Read(reader, binary.BigEndian, &p.Hostname)
	err = binary.Read(reader, binary.BigEndian, &p.TagLength)
	p.Tag = make([]byte, p.TagLength)
	err = binary.Read(reader, binary.BigEndian, &p.Tag)
	p.Msg = make([]byte, p.Length-8-2-p.HostnameLength-2-p.TagLength)
	err = binary.Read(reader, binary.BigEndian, &p.Msg)

	return err
}
