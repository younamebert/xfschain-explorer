package dao

import "github.com/jinzhu/gorm"

//MysqlDb mysql结构体
func Gorm() *gorm.DB {
	db, err := gorm.Open("mysql", "root:1263701671@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return db
}
