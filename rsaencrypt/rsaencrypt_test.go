/**
 * @Author: koulei
 * @Description:
 * @File: rsaencrypt_test.go
 * @Version: 1.0.0
 * @Date: 2021/8/20 13:02
 */

package rsaencrypt

import (
	"fmt"
	"testing"
	"time"
)

func TestRSAEncrypt(t *testing.T) {
	// GenerateRSAKey(1024)
	// 分段加密
	data := []byte("hello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello world哈哈哈")
	start := time.Now()
	encrypt, err := EncryptBlock(data, "public.pem")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(encrypt), time.Since(start))

	// 分段解密
	start = time.Now()
	decrypt, err := DecryptBlock(encrypt, "private.pem")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(decrypt), time.Since(start))
}
