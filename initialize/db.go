package initialize

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//MysqlDb mysql结构体
func Gorm() *gorm.DB {
	db, err := gorm.Open("mysql", "root:1263701671@(127.0.0.1:3306)/xfschain?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("gorm err:%v\n", err)
		panic(err)
	}
	defer db.Close()
	db.SingularTable(true)
	return db
}
