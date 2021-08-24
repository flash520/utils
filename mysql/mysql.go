/**
 * @Author: koulei
 * @Description: TODO
 * @File:  mysql
 * @Version: 1.0.0
 * @Date: 2021/5/12 18:14
 */

package mysql

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db  *gorm.DB
	err error
)

type Mysql struct {
}

// CreateMysql 初始化 MySQL 数据库连接
func CreateMysql(addr, username, password, dbname string, loglevel interface{}) *Mysql {
	var logLevel logger.LogLevel
	switch loglevel.(string) {
	default:
		logLevel = logger.Warn
	case "INFO":
		logLevel = logger.Info
	}
	var once sync.Once
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			username, password,
			addr,
			dbname)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger:      logger.Default.LogMode(logLevel),
			PrepareStmt: true,
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
		})
		if err != nil {
			fmt.Println(err)
		}

		if rawDB, err := db.DB(); err != nil {
			fmt.Println(err.Error())
		} else {
			rawDB.SetMaxIdleConns(10)
			rawDB.SetMaxOpenConns(100)
			rawDB.SetConnMaxLifetime(time.Hour)
		}
	})

	return &Mysql{}
}

func (m *Mysql) GetConn() *gorm.DB {
	return db
}
