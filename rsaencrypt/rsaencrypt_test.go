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
	GenerateRSAKey(1024)
	// 加密
	data := []byte("hello world")
	start := time.Now()
	encrypt, _ := RSAEncrypt(data, "public.pem")
	fmt.Println(string(encrypt), time.Since(start))

	// 解密
	start = time.Now()
	decrypt, _ := RSADecrypt(encrypt, "private.pem")
	fmt.Println(string(decrypt), time.Since(start))
}
