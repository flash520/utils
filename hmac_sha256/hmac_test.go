/**
 * @Author: koulei
 * @Description:
 * @File: hmac_test
 * @Version: 1.0.0
 * @Date: 2021/8/14 15:06
 */

package hmac_sha256

import (
	"encoding/json"
	"fmt"
	"testing"
)

const (
	secret = "123456789adfafa"
)

func TestComputeHmacSha256(t *testing.T) {
	var a = struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}{
		Id:   "abc",
		Name: "fly",
	}

	marshal, _ := json.Marshal(a)
	fmt.Println(string(marshal))
	sha256 := ComputeHmacSha256(string(marshal), secret)
	fmt.Println(sha256)

	b := `{"id":"abc","name":"fly"}`

	hmacSha256 := ComputeHmacSha256(b, secret)
	fmt.Println(hmacSha256)
}
