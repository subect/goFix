package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"goFix/dao"
	"log"
	"strconv"
	"time"
)

type Animal struct {
	gorm.Model
	AnimalId int64     `gorm:"column:beast_id"`
	Birthday time.Time `gorm:"column:day_of_the_beast"`
	Age      int64     `gorm:"column:age_of_the_beast"`
}

func Mysqltd(c *gin.Context) {
	db, err := dao.GetMysqlPool()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	//创建表
	//err = db.Table("animal").CreateTable(&Animal{}).Error
	//if err != nil {
	//	basicLog.Errorf("error creating Animal:%v", err)
	//}

	//插入
	//animalIn := &Animal{AnimalId: 2, Birthday: time.Now(), Age: 24}
	//err = db.Table("animal").Create(animalIn).Error
	//if err != nil {
	//	basicLog.Errorf("error creating Animal:%v", err)
	//}

	//查询表
	animalId, e := c.GetQuery("beast_id")
	if !e {
		basicLog.Errorf("this data is not exit")
		c.JSON(200, "this data is not exit")
	}

	animalIdq, err := strconv.ParseInt(animalId, 10, 64)
	if err != nil {
		return
	}
	var result Animal
	db.Table("animal").Where(&Animal{AnimalId: animalIdq}).Find(&result)
	//basicLog.Debugln("success create table")
	c.JSON(200, result)

}
