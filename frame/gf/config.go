/**
 * @Author: koulei
 * @Description:
 * @File: config
 * @Version: 1.0.0
 * @Date: 2021/8/26 17:35
 */

package gf

import "github.com/spf13/viper"

func ParseConfig() interface{} {
	var config interface{}
	v := viper.New()
	v.SetConfigFile("config/config.yml")
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	return config
}
