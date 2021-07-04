/**
 * @Author: koulei
 * @Description: TODO
 * @File:  eureka
 * @Version: 1.0.0
 * @Date: 2021/7/3 18:28
 */

package discovery

import (
	"errors"
	"strings"
	"sync"

	eureka "github.com/xuanbo/eureka-client"
)

type Eureka struct {
	eurekaClient *eureka.Client
}

func CreateEureka(username, password, node, app string, port int) *Eureka {
	e := Eureka{}
	var once sync.Once
	once.Do(func() {
		e.eurekaClient = eureka.NewClient(&eureka.Config{
			DefaultZone:           "http://" + username + ":" + password + "@" + node,
			App:                   app,
			Port:                  port,
			RenewalIntervalInSecs: 15,
			DurationInSecs:        30,
			Metadata: map[string]interface{}{
				"VERSION":              "0.1.0",
				"NODE_GROUP_ID":        0,
				"PRODUCT_CODE":         "DEFAULT",
				"PRODUCT_VERSION_CODE": "DEFAULT",
				"PRODUCT_ENV_CODE":     "DEFAULT",
				"SERVICE_VERSION_CODE": "DEFAULT",
			},
		})
	})
	go e.eurekaClient.Start()
	return &Eureka{}
}

func (e *Eureka) GetInstance(app string) (string, error) {
	apps := e.eurekaClient.Applications
	var instance string
	for _, v := range apps.Applications {
		if v.Name == strings.ToUpper(app) {
			for i := 0; i < len(v.Instances); i++ {
				instance = v.Instances[i].HomePageURL
			}
			break
		}
		return "", errors.New(strings.ToUpper(instance + " " + "找不到该服务"))
	}

	return instance, err
}
