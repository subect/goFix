package model

import (
	"github.com/jinzhu/gorm"
	"go-component-library/mysql"
)

var dbClient *gorm.DB

// Initialize 在应用程序启动时初始化数据库连接
func Initialize() error {
	// 获取数据库连接
	var err error
	dbClient, err = GetMysqlClient()
	if err != nil {
		return err
	}
	return nil
}

func GetDb() *gorm.DB {
	if dbClient == nil {
		dbClient, _ = GetMysqlClient()
	}
	return dbClient
}

func GetMysqlDsn() string {
	mysqlConf := conf.MysqlServer
	dbName := mysqlConf.DefaultDbName
	dbUser := mysqlConf.UserName
	dbPwd := mysqlConf.PassWord
	dbAddr := mysqlConf.Address
	dsn := mysql.NewDbDsn(dbName, dbUser, dbPwd, dbAddr)
	return dsn
}

func GetMysqlClient() (*gorm.DB, error) {
	dsn := GetMysqlDsn()
	devMode := conf.DevMode
	db, err := mysql.NewClient(dsn, devMode)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CloseDb 关闭数据库连接
func CloseDb() {
	if dbClient != nil {
		dbClient.Close()
	}
}
