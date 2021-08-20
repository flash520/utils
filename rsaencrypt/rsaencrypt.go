/**
 * @Author: koulei
 * @Description:
 * @File: rsaEncrypt
 * @Version: 1.0.0
 * @Date: 2021/8/20 12:04
 */

package rsaencrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// GenerateRSAKey 生成RSA私钥和公钥，保存到文件中
// bits 证书大小
func GenerateRSAKey(bits int) {
	// GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	// Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	// 保存私钥
	// 通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	// 使用pem格式对x509输出的内容进行编码
	// 创建文件保存私钥
	privateFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer func() { _ = privateFile.Close() }()
	// 构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	// 将数据保存到文件
	_ = pem.Encode(privateFile, &privateBlock)

	// 保存公钥
	// 获取公钥的数据
	publicKey := privateKey.PublicKey
	// X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	// pem格式编码
	// 创建用于保存公钥的文件
	publicFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	defer func() { _ = publicFile.Close() }()
	// 创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	// 保存到文件
	_ = pem.Encode(publicFile, &publicBlock)
}

// Encrypt 加密
// plainText 加密内容
// path 公钥文件地址
func Encrypt(plainText []byte, path string) ([]byte, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	// 读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, _ = file.Read(buf)
	// pem解码
	block, _ := pem.Decode(buf)
	// x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	// 对明文进行加密
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
}

// EncryptBlock 分段加密
// plainText 加密内容
// path 公钥文件地址
func EncryptBlock(plainText []byte, path string) ([]byte, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	// 读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, _ = file.Read(buf)
	// pem解码
	block, _ := pem.Decode(buf)
	// x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	// 分段加密
	offset, once, srcSize := 0, publicKey.Size()-11, len(plainText)
	buffer := bytes.Buffer{}
	for offset < srcSize {
		endIndex := offset + once
		if endIndex > len(plainText) {
			endIndex = srcSize
		}
		// 加密一部分
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText[offset:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offset = endIndex
	}
	return buffer.Bytes(), nil
}

// Decrypt 解密
// cipherText 需要解密的byte数据
// path 私钥文件路径
func Decrypt(cipherText []byte, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	// 获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, _ = file.Read(buf)
	// pem解码
	block, _ := pem.Decode(buf)
	// X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 对密文进行解密
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
}

// DecryptBlock 分段解密
// cipherText 需要解密的byte数据
// path 私钥文件路径
func DecryptBlock(cipherText []byte, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	// 获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, _ = file.Read(buf)
	// pem解码
	block, _ := pem.Decode(buf)
	// X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	keySize := privateKey.Size()
	srcSize := len(cipherText)
	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + keySize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return buffer.Bytes(), nil
}
