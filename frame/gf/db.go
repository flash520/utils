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
	"io/fs"
	"os"
	"path/filepath"
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
	var config string
	err := filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(path, "config.yml") {
			config = path
		}
		return nil
	})
	if err != nil {
		os.Exit(1)
	}
	v := viper.New()
	v.SetConfigFile(config)
	err = v.ReadInConfig()
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
