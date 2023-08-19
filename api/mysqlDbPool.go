package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"goFix/model"
)

func MysqlDbPool(c *gin.Context) {

	var userInfo model.User

	dbClient := model.GetDb()
	dbClient.Where("name = ?", "huxiang").Find(&userInfo)

	fmt.Println("userInfo:", userInfo)
	c.JSON(200, gin.H{"code": 200, "msg": "success", "data": userInfo})

	//fmt.Println("start MysqlDbPool")

	// 初始化数据库连接池
	//initDB()

	//// 创建表
	//if err := model.CreateTable() ; err != nil {
	//	log.Fatalf("failed to create table: %v", err)
	//}
	//
	//// 插入数据
	//lastInsertID, err := model.InsertData("John", 30)
	//if err != nil {
	//	log.Fatalf("failed to insert data: %v", err)
	//}
	//fmt.Printf("Inserted ID: %d\n", lastInsertID)

	// 查询数据
	//users, err := model.QueryData()
	//if err != nil {
	//	log.Fatalf("failed to query data: %v", err)
	//}
	//fmt.Println("Users:")
	//for _, user := range users {
	//	fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	//}
	//
	////respJson, _ := json.Marshal(users)
	//c.JSON(200, gin.H{"code": 200, "msg": "success", "data": users})

	//
	//// 更新数据
	//if err := model.UpdateData(lastInsertID, "John Doe", 31); err != nil {
	//	log.Fatalf("failed to update data: %v", err)
	//}
	//
	//// 查询更新后的数据
	//users, err = model.QueryData()
	//if err != nil {
	//	log.Fatalf("failed to query data: %v", err)
	//}
	//fmt.Println("Updated Users:")
	//for _, user := range users {
	//	fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	//}
	//
	//// 删除数据
	//if err := model.DeleteData(lastInsertID); err != nil {
	//	log.Fatalf("failed to delete data: %v", err)
	//}
	//
	//// 查询删除后的数据
	//users, err = model.QueryData()
	//if err != nil {
	//	log.Fatalf("failed to query data: %v", err)
	//}
	//fmt.Println("Remaining Users:")
	//for _, user := range users {
	//	fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	//}
}
