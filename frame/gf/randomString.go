/**
 * @Author: koulei
 * @Description:
 * @File: randomString
 * @Version: 1.0.0
 * @Date: 2021/12/26 19:13
 */

package gf

import "gitee.com/flash520/utils/randomstring"

var RandomStr = randomStr{}

type randomStr struct {
}

// New 获取随机数对象
// 1、混合型字符串，
// 2、数字型字符串
// 3、字母型字符串（包含大小写）
func (r *randomStr) New(n int) *randomstring.RandomString {
	switch n {
	case 1:
		return randomstring.NewRandomString(n)
	case 2:
		return randomstring.NewRandomString(n)
	case 3:
		return randomstring.NewRandomString(n)
	default:
		return nil
	}
}
