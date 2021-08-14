/**
 * @Author: koulei
 * @Description:
 * @File: hmac_sha256
 * @Version: 1.0.0
 * @Date: 2021/8/14 14:57
 */

package hmac_sha256

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func ComputeHmacSha256(str, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}
