package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"goFix/config"
	"sync"
	"time"
)

type DbPool struct {
	maxDbs int
	Dbs    map[string]*gorm.DB
	mux    *sync.RWMutex
}

func (dp *DbPool) Get(dbName string) (*gorm.DB, error) {
	basicLog.Debugf("Get db:%v", dbName)
	dp.mux.Lock()
	defer dp.mux.Unlock()
	if db, ok := dp.Dbs[dbName]; ok {
		return db, nil
	}
	if len(dp.Dbs) > dp.maxDbs {
		for k, db := range dp.Dbs {
			err := db.Close()
			if err != nil {
				return nil, err
			}
			delete(dp.Dbs, k)
			break
		}
	}
	newDb, err := createNewDBConn(dbName)
	if err != nil {
		return nil, err
	}
	return newDb, nil
}

func createNewDBConn(dbName string) (*gorm.DB, error) {
	mysqlAddr := config.Config().MysqlServer.Address
	mysqlUserName := config.Config().MysqlServer.UserName
	mysqlPassWord := config.Config().MysqlServer.PassWord
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=30s&charset=utf8mb4&collation=utf8mb4_general_ci&loc=Local&parseTime=true", mysqlUserName, mysqlPassWord, mysqlAddr, dbName)
	basicLog.Debugf("mysql open addr:{dsn:%v", dsn)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		basicLog.Errorf("coon faul:%v", err.Error())
		return nil, err
	}
	db.DB().SetMaxIdleConns(30)
	db.DB().SetMaxOpenConns(10)
	db.DB().SetConnMaxLifetime(1 * time.Minute)

	err = db.DB().Ping()
	if err != nil {
		basicLog.Errorf("ping error:%v", err.Error())
		return nil, err
	}

	db.LogMode(config.Config().DevMode)

	return db, nil
}

var ClientPool *DbPool

func NewDbPool() *DbPool {
	ClientPool = &DbPool{
		maxDbs: config.Config().MysqlServer.MysqlMaxDBs,
		Dbs:    make(map[string]*gorm.DB, 2),
		mux:    new(sync.RWMutex),
	}
	return ClientPool
}

func GetMysqlPool() (*gorm.DB, error) {
	dbName := config.Config().MysqlServer.DefaultDbName
	DbInstance, err := ClientPool.Get(dbName)
	return DbInstance, err
}

func (dp *DbPool) CloseDb() {
	if dp == nil {
		return
	}
	for k, db := range dp.Dbs {
		basicLog.Debugf("exit closing db coon:%v", k)
		err := db.Close()
		if err != nil {
			return
		}
	}
}
