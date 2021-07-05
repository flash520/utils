/**
 * @Author: koulei
 * @Description: TODO
 * @File:  md5_encrypt
 * @Version: 1.0.0
 * @Date: 2021/7/5 12:18
 */

package md5

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func MD5(param string) string {
	md5ctx := md5.New()
	md5ctx.Write([]byte(param))
	return hex.EncodeToString(md5ctx.Sum(nil))
}

// Base64MD5 先 base64 ,然后 md5
func Base64MD5(param string) string {
	return MD5(base64.StdEncoding.EncodeToString([]byte(param)))
}
