/**
 * @Author: koulei
 * @Description:
 * @File: randstring_test.go
 * @Version: 1.0.0
 * @Date: 2021/7/26 23:54
 */

package randomstring

import "testing"

func benchmarkRandString(b *testing.B, function func(n int) string) {
	for i := 0; i < b.N; i++ {
		function(6)
	}
}

func BenchmarkRandStringRunes(b *testing.B) {
	benchmarkRandString(b, RandStringRunes)
}

func BenchmarkRandStringBytes(b *testing.B) {
	benchmarkRandString(b, RandStringBytes)
}

func BenchmarkRandStringBytesMask(b *testing.B) {
	benchmarkRandString(b, RandStringBytesMask)
}

func BenchmarkRandStringBytesRmndr(b *testing.B) {
	benchmarkRandString(b, RandStringBytesRmndr)
}

func BenchmarkRandStringBytesMaskImpr(b *testing.B) {
	benchmarkRandString(b, RandStringBytesMaskImpr)
}

func BenchmarkRandStringBytesMaskImprSrc(b *testing.B) {
	benchmarkRandString(b, RandStringBytesMaskImprSrc)
}

func BenchmarkRandStringBytesMaskImprSrcSB(b *testing.B) {
	benchmarkRandString(b, RandStringBytesMaskImprSrcSB)
}

func BenchmarkRandStringBytesMaskImprSrcUnsafe(b *testing.B) {
	benchmarkRandString(b, RandStringBytesMaskImprSrcUnsafe)
}
