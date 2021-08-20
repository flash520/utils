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

var private = []byte(`
-----BEGIN RSA Private Key-----
MIICXAIBAAKBgQDcIUYxGQ7t3xDJ3p6/pyajPQTc3r/cy8xUTmLTGKLttkpI+Q+d
I2PA4Mr2VOU3U8rC46cAaO5WAfIJzRsPDjlAHn4l2hiimnQmMaFE80Le9pBcHkLh
mrk7IcttbM3bO+f0rfKhLnLZKlvbrRMQ4k0AWQYxum2CpHeix54Xhs+fOwIDAQAB
AoGAC1AKc2t+QOs9yaIPNno4mhsArPklkwtGLO76VS7m8KB1oNpr2v9+mOL0i0RB
o15DBVD9vB+oX/MppSdNRLDOcEy5Mm3HqBD14Y9KvryCOSnUI9KOT0qXsog/brLB
UqPO0B80eYL8z4D4blUFds49tetzPwjysiM5EBOj5fVJxwECQQDobrQGCg9iIdSf
KHzPK7rNK6/Xv+4G1jAgN4Yh2qjsoYFKMm1wI3yWeVsbD5nFhLi3TUcRIwIeiBWB
9BQ6+2TrAkEA8nM7yY+moay6+/Mnqmb55X5kA3YKnbFmOMNp3OL5+gi8gmnNBPlu
/RY94CqzB8O0/0cEfWDELYV+7cyNPUFa8QJBAMM270Y/PspZxZ1jQOgzPzJA5fBb
x8vAKy1z1NksMEKGJvOtRNMxStuK02T4BlupbAawBeczsvz3qLC70h3ztL8CQBJY
c9qxmvs61b9Ay+yR9DDQWoMEiJMcHE8JQlZrelzYEmOP2+qXcTWHXFE9CeT5gxWZ
6xYNn2qOYmdeBgXvjxECQFgtNCK4TkAWrw5UKQ+RdCZXaXn6VJm4Vwb++bS7meT2
Gtqa/jr6rjuyT4CnBBT3aJiDJeWnLs+lIbA5V1+ZcV0=
-----END RSA Private Key-----
`)
var public = []byte(`
-----BEGIN RSA Public Key-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDcIUYxGQ7t3xDJ3p6/pyajPQTc
3r/cy8xUTmLTGKLttkpI+Q+dI2PA4Mr2VOU3U8rC46cAaO5WAfIJzRsPDjlAHn4l
2hiimnQmMaFE80Le9pBcHkLhmrk7IcttbM3bO+f0rfKhLnLZKlvbrRMQ4k0AWQYx
um2CpHeix54Xhs+fOwIDAQAB
-----END RSA Public Key-----
`)

func TestRSAEncrypt(t *testing.T) {
	// GenerateRSAKey(1024)
	// 分段加密
	data := []byte("hello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello worldhello world哈哈哈")
	start := time.Now()
	encrypt, err := EncryptBlock(data, public)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(encrypt), time.Since(start))

	// 分段解密
	start = time.Now()
	decrypt, err := DecryptBlock(encrypt, private)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(decrypt), time.Since(start))
}
