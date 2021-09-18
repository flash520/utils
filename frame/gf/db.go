/**
 * @Author: koulei
 * @Description:
 * @File: db
 * @Version: 1.0.0
 * @Date: 2021/8/26 17:49
 */

package gf

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"gitee.com/flash520/utils/container"
	"gitee.com/flash520/utils/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const (
	mysqlClient string = "MysqlClient"
)

var (
	mutex sync.Mutex
)

func DB(loglevel string, config ...string) *gorm.DB {
	mutex.Lock()
	defer func() { mutex.Unlock() }()
	c := container.CreateContainersFactory()
	r := c.Get(mysqlClient)
	if r == nil {
		log.Debug("新建数据库对象")
		if len(config) > 0 {
			if rds, err := newDB(strings.ToUpper(loglevel), config[0]); err != nil {
				return nil
			} else {
				c.Set(mysqlClient, rds)
				return rds
			}
		} else {
			if rds, err := newDB(strings.ToUpper(loglevel)); err != nil {
				return nil
			} else {
				c.Set(mysqlClient, rds)
				return rds
			}
		}

	}
	log.Debug("缓存获取数据库对象")
	rds := r.(*gorm.DB)
	return rds
}

func newDB(loglevel string, config ...string) (*gorm.DB, error) {
	v := viper.New()
	if len(config) > 0 {
		v.SetConfigFile(config[0])
	} else {
		v.SetConfigFile("./config/config.yml")
		getwd, _ := os.Getwd()
		fmt.Println("current path: ", getwd)
	}
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	r := mysql.CreateMysql(
		v.GetString("mysql.hosts"),
		v.GetString("mysql.username"),
		v.GetString("mysql.password"),
		v.GetString("mysql.dbname"), loglevel).GetConn()
	if r == nil {
		return nil, errors.New("数据库初始化失败")
	}
	return r, nil
}
