/**
 * @Author: koulei
 * @Description:
 * @File: main_test.go
 * @Version: 1.0.0
 * @Date: 2021/7/20 23:27
 */

package main

import (
	"net"
	"testing"
)

func BenchmarkSend(b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			c, err := net.Dial("tcp", ":8000")
			if err != nil {
				return
			}
			sendfile(err, c)
		}()
	}
}
