/**
 * @Author: koulei
 * @Description:
 * @File: config_test.go
 * @Version: 1.0.0
 * @Date: 2021/8/26 17:40
 */

package gf

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestParseConfig(t *testing.T) {
	a := func() interface{} {
		var b interface{}
		v := viper.New()
		v.SetConfigFile("config.yml")
		v.SetConfigType("yml")
		err := v.ReadInConfig()
		if err != nil {
			panic(err)
		}
		err = v.Unmarshal(&b)
		if err != nil {
			panic(err)
		}
		return b
	}

	fmt.Println(a())
}
