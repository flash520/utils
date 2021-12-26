/**
 * @Author: koulei
 * @Description: 生成随机字符串
 * @File: randomstring
 * @Version: 1.0.0
 * @Date: 2021/7/26 23:42
 */

package randomstring

import (
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type RandomString struct {
	letterRunes []rune
	letterBytes string
}

// NewRandomString 生成随机字符串对象
// 1、混合型字符串，
// 2、数字型字符串
// 3、字母型字符串（包含大小写）
func NewRandomString(n int) *RandomString {
	switch n {
	case 1:
		return &RandomString{
			letterRunes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789="),
			letterBytes: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789=",
		}
	case 2:
		return &RandomString{
			letterRunes: []rune("0123456789"),
			letterBytes: "0123456789",
		}
	case 3:
		return &RandomString{
			letterRunes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
			letterBytes: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		}
	default:
		return nil
	}
}

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringRunes 使用 rune 存储随机数
func (rs *RandomString) RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rs.letterRunes[rand.Intn(len(rs.letterRunes))]
	}
	return string(b)
}

func (rs *RandomString) RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rs.letterBytes[rand.Intn(len(rs.letterBytes))]
	}
	return string(b)
}

func (rs *RandomString) RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rs.letterBytes[rand.Int63()%int64(len(rs.letterBytes))]
	}
	return string(b)
}

func (rs *RandomString) RandStringBytesMask(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(rs.letterBytes) {
			b[i] = rs.letterBytes[idx]
			i++
		}
	}
	return string(b)
}

func (rs *RandomString) RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(rs.letterBytes) {
			b[i] = rs.letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

type Source interface {
	Int63() int64
	Seed(seed int64)
}

var src = rand.NewSource(time.Now().UnixNano())

func (rs *RandomString) RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(rs.letterBytes) {
			b[i] = rs.letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func (rs *RandomString) RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(rs.letterBytes) {
			sb.WriteByte(rs.letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// String returns the accumulated string.
// func (b *Builder) String() string {
// 	return *(*string)(unsafe.Pointer(&b.buf))
// }

func (rs *RandomString) RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(rs.letterBytes) {
			b[i] = rs.letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
