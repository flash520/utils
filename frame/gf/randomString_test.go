/**
 * @Author: koulei
 * @Description:
 * @File: randomString_test.go
 * @Version: 1.0.0
 * @Date: 2021/12/26 19:37
 */

package gf

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(RandomStr.New(1).RandStringBytesMaskImprSrcUnsafe(10))
}
