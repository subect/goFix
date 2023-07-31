package model

import (
	"database/sql"
	"fmt"
	"goFix/config"
	"log"
)

var db *sql.DB

//初始化 数据库连接信息
func initMysqlDB() {
	// ... 在此处配置数据库连接信息 ...
	conf := config.Config().MysqlServer
	dbUser, dbPwd, dbAddr, dbName := conf.UserName, conf.PassWord, conf.Address, conf.DefaultDbName
	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=30s&charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=Local&parseTime=true", dbUser, dbPwd, dbAddr, dbName)

	// 使用数据库驱动打开数据库连接
	var err error
	db, err = sql.Open("mysql", uri)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 设置连接池的最大空闲连接数和最大打开连接数
	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(10)
}

// CreateTable 创建表
func CreateTable() error {
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			age INT NOT NULL
		)
	`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

// InsertData 插入数据
func InsertData(name string, age int) (int64, error) {
	result, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", name, age)
	if err != nil {
		return 0, fmt.Errorf("failed to insert data: %v", err)
	}

	// 获取插入的自增ID
	lastInsertID, _ := result.LastInsertId()
	return lastInsertID, nil
}

// QueryData 查询数据
func QueryData() ([]User, error) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateData 更新数据
func UpdateData(id int64, name string, age int) error {
	_, err := db.Exec("UPDATE users SET name=?, age=? WHERE id=?", name, age, id)
	if err != nil {
		return fmt.Errorf("failed to update data: %v", err)
	}

	return nil
}

// DeleteData 删除数据
func DeleteData(id int64) error {
	_, err := db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("failed to delete data: %v", err)
	}

	return nil
}

// User 结构体用于保存从数据库中查询出来的数据
type User struct {
	ID   int64
	Name string
	Age  int
}
