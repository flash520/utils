/**
 * @Author: koulei
 * @Description:
 * @File: cale_test.go
 * @Version: 1.0.0
 * @Date: 2021/7/19 16:58
 */

package pack

import "testing"

func BenchmarkPackage_Pack(b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cale1(i)
	}

}

func BenchmarkCacel1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cale(i)
	}
}
