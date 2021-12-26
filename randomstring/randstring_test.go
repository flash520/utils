/**
 * @Author: koulei
 * @Description:
 * @File: randstring_test.go
 * @Version: 1.0.0
 * @Date: 2021/7/26 23:54
 */

package randomstring

import (
	"fmt"
	"testing"
	"time"
)

func testRandStringBytes(t *testing.T, function func(n int) string) {
	start := time.Now()
	fmt.Printf("%-15s | time: %-10v\n", function(15), time.Since(start))
}

func TestRandStringBytes(t *testing.T) {
	rs := NewRandomString(1)
	testRandStringBytes(t, rs.RandStringBytes)
	testRandStringBytes(t, rs.RandStringBytesMask)
	rs = NewRandomString(2)
	testRandStringBytes(t, rs.RandStringBytesMaskImpr)
	testRandStringBytes(t, rs.RandStringBytesRmndr)
	testRandStringBytes(t, rs.RandStringRunes)
	rs = NewRandomString(3)
	testRandStringBytes(t, rs.RandStringBytesMaskImprSrc)
	testRandStringBytes(t, rs.RandStringBytesMaskImprSrcSB)
	testRandStringBytes(t, rs.RandStringBytesMaskImprSrcUnsafe)
}

func benchmarkRandString(b *testing.B, function func(n int) string) {
	for i := 0; i < b.N; i++ {
		function(6)
	}
}

func BenchmarkRandStringRunes(b *testing.B) {
	rs := NewRandomString(1)
	benchmarkRandString(b, rs.RandStringRunes)
}

func BenchmarkRandStringBytes(b *testing.B) {
	rs := NewRandomString(1)
	benchmarkRandString(b, rs.RandStringBytes)
}

func BenchmarkRandStringBytesMask(b *testing.B) {
	rs := NewRandomString(1)
	benchmarkRandString(b, rs.RandStringBytesMask)
}

func BenchmarkRandStringBytesRmndr(b *testing.B) {
	rs := NewRandomString(1)
	benchmarkRandString(b, rs.RandStringBytesRmndr)
}

func BenchmarkRandStringBytesMaskImpr(b *testing.B) {
	rs := NewRandomString(1)
	benchmarkRandString(b, rs.RandStringBytesMaskImpr)
}

func BenchmarkRandStringBytesMaskImprSrc(b *testing.B) {
	rs := NewRandomString(1)
	benchmarkRandString(b, rs.RandStringBytesMaskImprSrc)
}

func BenchmarkRandStringBytesMaskImprSrcSB(b *testing.B) {
	rs := NewRandomString(1)
	benchmarkRandString(b, rs.RandStringBytesMaskImprSrcSB)
}

func BenchmarkRandStringBytesMaskImprSrcUnsafe(b *testing.B) {
	rs := NewRandomString(1)
	benchmarkRandString(b, rs.RandStringBytesMaskImprSrcUnsafe)
}
