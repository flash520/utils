/**
 * @Author: koulei
 * @Description:
 * @File: container
 * @Version: 1.0.0
 * @Date: 2021/8/3 10:44
 */

package container

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var sMap sync.Map

// CreateContainersFactory 创建一个容器工厂
func CreateContainersFactory() *Containers {
	return &Containers{}
}

// Containers 定义一个容器结构体
type Containers struct {
}

// Set 1.以键值对的形式将代码注册到容器
func (c *Containers) Set(key string, value interface{}) (res bool) {

	if _, exists := c.KeyIsExists(key); exists == false {
		sMap.Store(key, value)
		res = true
	} else {
		// 程序启动阶段，zaplog 未初始化，使用系统log打印启动时候发生的异常日志
		log.Warnf("请解决键名重复问题,相关键：%s\n", key)
		res = false
	}
	return
}

// Delete 2.删除
func (c *Containers) Delete(key string) {
	sMap.Delete(key)
}

// Get 3.传递键，从容器获取值
func (c *Containers) Get(key string) interface{} {
	if value, exists := c.KeyIsExists(key); exists {
		return value
	}
	return nil
}

// KeyIsExists 4. 判断键是否被注册
func (c *Containers) KeyIsExists(key string) (interface{}, bool) {
	return sMap.Load(key)
}

// FuzzyDelete 按照键的前缀模糊删除容器中注册的内容
func (c *Containers) FuzzyDelete(keyPre string) {
	sMap.Range(func(key, value interface{}) bool {
		if keyname, ok := key.(string); ok {
			if strings.HasPrefix(keyname, keyPre) {
				sMap.Delete(keyname)
			}
		}
		return true
	})
}
