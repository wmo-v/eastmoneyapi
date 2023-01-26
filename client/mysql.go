package client

import (
	"eastmoneyapi/config"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlClient *gorm.DB
var initMysqlOnce sync.Once

func GetMysqlClient() *gorm.DB {
	if mysqlClient != nil {
		return mysqlClient
	}

	initMysqlOnce.Do(func() {
		conf := config.GetConfg().Mysql
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.User,
			conf.Passwd,
			conf.Host,
			conf.Port,
			conf.DBName,
		)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("数据库初始化失败：" + err.Error())
		}
		mysqlClient = db
	})
	return mysqlClient
}
