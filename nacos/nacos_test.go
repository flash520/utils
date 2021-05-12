/**
 * @Author: koulei
 * @Description: TODO
 * @File:  nacos_test
 * @Version: 1.0.0
 * @Date: 2021/5/11 23:11
 */

package nacos

import (
	"testing"
)

func TestNacos(t *testing.T) {
	CreateNacos("nacos.dev.cdqidi.cn", "/nacos", "public", "fly", 80, 80)
}
