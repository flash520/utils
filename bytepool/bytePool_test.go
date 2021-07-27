/**
 * @Author: koulei
 * @Description:
 * @File: bytePool_test.go
 * @Version: 1.0.0
 * @Date: 2021/7/27 12:43
 */

package bytepool

import "testing"

func BenchmarkNewBytePool(b *testing.B) {
	bp := NewBytePool(10000, 4096, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := bp.GET()
		d = append(d, 0)
		bp.PUT(d)
	}
}
