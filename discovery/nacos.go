/**
 * @Author: koulei
 * @Description: TODO
 * @File:  nacos
 * @Version: 1.0.0
 * @Date: 2021/5/11 22:43
 */

package discovery

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	nacosClient naming_client.INamingClient
	err         error
)

type Nacos struct {
}

// CreateNacos 初始化nacos 并将应用注册到nacos
func CreateNacos(node, contextPath, nameSpaceID, appName string, Port, appPort uint64) *Nacos {
	var once sync.Once
	once.Do(func() {
		sc := []constant.ServerConfig{
			{
				IpAddr:      node,
				Port:        Port,
				ContextPath: contextPath,
			},
		}
		cc := constant.ClientConfig{
			NamespaceId:         nameSpaceID,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			RotateTime:          "1h",
			MaxAge:              3,
			LogLevel:            "info",
		}

		nacosClient, err = clients.CreateNamingClient(map[string]interface{}{
			"serverConfigs": sc,
			"clientConfig":  cc,
		})
		if err != nil {
			panic(err)
		}

		RegisterServiceInstance(vo.RegisterInstanceParam{
			Ip:          util.LocalIP(),
			Port:        appPort,
			ServiceName: appName,
			Weight:      10,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			Metadata:    map[string]string{"idc": "chengdu"},
		})

		signalNotify(appName, appPort)
	})

	return &Nacos{}
}

func signalNotify(appName string, appPort uint64) {
	signs := make(chan os.Signal, 1)
	signal.Notify(signs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		switch <-signs {
		case syscall.SIGHUP:
			fallthrough
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGKILL:
			fallthrough
		case syscall.SIGTERM:
			DeRegisterServiceInstance(vo.DeregisterInstanceParam{
				Ip:          util.LocalIP(),
				Port:        appPort,
				Cluster:     "",
				ServiceName: appName,
				GroupName:   "",
				Ephemeral:   true,
			})
			os.Exit(0)
		}
	}()
}
