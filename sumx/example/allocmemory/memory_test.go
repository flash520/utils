/**
 * @Author: koulei
 * @Description:
 * @File: memory_test.go
 * @Version: 1.0.0
 * @Date: 2021/7/22 20:13
 */

package allocmemory

import "testing"

func BenchmarkPreAlloc(b *testing.B) {
	bufChan := make(chan []byte, 1000000)
	for i := 0; i < 1000000; i++ {
		bufChan <- make([]byte, 1000)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prefixAllocMemory(bufChan)
	}
}

func BenchmarkAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		varAllocMemory()
	}
}
