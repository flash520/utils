/**
 * @Author: koulei
 * @Description:
 * @File: cale
 * @Version: 1.0.0
 * @Date: 2021/7/19 16:46
 */

package pack

/**
 * @Description: 判断数字是否为偶数
 * @return bool
 */
func cale(number int) bool {
	return number%2 == 0
}

/**
 * @Description:
 * @param number
 * @return bool
 */
func cale1(number int) bool {

	return number^1 == 0
}
