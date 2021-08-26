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

func DB(loglevel string) *gorm.DB {
	mutex.Lock()
	defer func() { mutex.Unlock() }()
	c := container.CreateContainersFactory()
	r := c.Get(mysqlClient)
	if r == nil {
		log.Debug("新建数据库对象")
		if rds, err := newDB(strings.ToUpper(loglevel)); err != nil {
			return nil
		} else {
			c.Set(mysqlClient, rds)
			return rds
		}
	}
	log.Debug("缓存获取数据库对象")
	rds := r.(*gorm.DB)
	return rds
}

func newDB(loglevel string) (*gorm.DB, error) {
	v := viper.New()
	v.SetConfigFile("config/config.yml")
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	stringMap := v.GetStringMap("mysql")

	r := mysql.CreateMysql(stringMap["hosts"].(string), stringMap["username"].(string), stringMap["password"].(string),
		stringMap["dbname"].(string), loglevel).GetConn()
	if r == nil {
		return nil, errors.New("数据库初始化失败")
	}
	return r, nil
}
