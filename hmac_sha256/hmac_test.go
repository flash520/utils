/**
 * @Author: koulei
 * @Description:
 * @File: hmac_test
 * @Version: 1.0.0
 * @Date: 2021/8/14 15:06
 */

package hmac_sha256

import (
	"fmt"
	"testing"
)

const (
	secret = "123456789adfafa"
)

func TestComputeHmacSha256(t *testing.T) {
	sha256 := ComputeHmacSha256("hello world", secret)
	fmt.Println(sha256)
}
