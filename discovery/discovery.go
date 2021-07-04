/**
 * @Author: koulei
 * @Description: TODO
 * @File:  discovery
 * @Version: 1.0.0
 * @Date: 2021/7/3 16:56
 */

package discovery

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Discovery interface {
	// GetInstance 从服务发现获取应用地址
	GetInstance(app string) (string, error)
}

type Registry struct {
	DiscoveryType Discovery
}

// InitDiscovery 初始化注册中心
//
// discoveryType 为 eureka 时，config 参数共 4 个，依次为:"用户名", "密码", "eureka地址:7000/eureka/", 应用端口
//
// discoveryType 为 nacos时，config 参数共 5个，依次为:nacos地址, contextPath, nameSpaceID, nacosPort, appPort
func InitDiscovery(discoveryType string, app string, config ...interface{}) *Registry {
	if len(config) < 1 {
		log.Errorf("注册中心参数不足")
		os.Exit(-1)
	}

	r := &Registry{}

	switch discoveryType {
	case "eureka":
		if len(config) != 4 {
			log.Errorf("eureka 注册中心初始化失败, 配置参数数量不足")
			os.Exit(-1)
		}
		r.DiscoveryType = CreateEureka(config[0].(string), config[1].(string), config[2].(string),
			app, config[3].(int))
	case "nacos":
		if len(config) != 5 {
			log.Errorf("nacos 注册中心初始化失败, 配置参数数量不足")
			os.Exit(-1)
		}
		r.DiscoveryType = CreateNacos(config[0].(string), config[1].(string),
			config[2].(string), app, config[3].(uint64), config[4].(uint64))
	default:
		log.Errorf("注册中心类型错误")
		os.Exit(-1)
	}

	return r
}

func (r *Registry) GetInstance(app string) string {
	result, err := r.DiscoveryType.GetInstance(app)
	if err != nil {
		log.Errorf(err.Error())
	}
	return result
}
