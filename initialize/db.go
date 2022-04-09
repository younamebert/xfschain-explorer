package initialize

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mi/common"
	"mi/conf"
	"mi/global"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//MysqlDb mysql结构体
func Gorm() *gorm.DB {
	db, err := gorm.Open("mysql", "wl:123456@(120.26.217.42:3306)/wl?charset=utf8mb4&parseTime=True&loc=Local")
	// defer db.Close()
	// db.SetLogger(true)
	db.LogMode(false)
	if err != nil {
		fmt.Printf("gorm err:%v\n", err)
		os.Exit(1)
	}

	db.DB().SetMaxIdleConns(50)
	db.DB().SetMaxOpenConns(50)
	db.DB().SetConnMaxLifetime(time.Minute)
	// if err := installMysql(db); err != nil {
	// 	fmt.Printf("installMysql err:%v\n", err)
	// 	os.Exit(1)
	// }

	// http://java.sun.com/xml/ns/javaee
	// http://java.sun.com/xml/ns/javaee/web-app_2_5.xsd
	// http://java.sun.com/xml/ns/javaee/web-app_3_0.xsd

	// fmt.Println(12)
	db.SingularTable(true)
	return db
}

func installMysql(db *gorm.DB) error {
	dataExit, _ := common.IsFileExist(conf.SqlFile)
	if !dataExit {
		return errors.New("file xfschain.sql non-existent")
	}
	sqls, err := ioutil.ReadFile(conf.SqlFile)
	if err != nil {
		return err
	}
	sqlArr := strings.Split(string(sqls), "|")
	for _, sql := range sqlArr {
		if sql == "" {
			continue
		}
		if err := db.Exec(sql).Error; err != nil {
			global.GVA_LOG.Error("table exists in database")
			continue
		}
	}
	return nil
}

func DBList() {
	dbMap := make(map[string]*gorm.DB)
	for _, dbname := range conf.DBlist {
		dbMap[dbname] = global.GVA_DB.Table(dbname)
	}
	// 做特殊判断,是否有迁移
	// 适配低版本迁移多数据库版本
	// if sysDB, ok := dbMap[sys]; ok {
	// 	global.GVA_DB = sysDB
	// }
	global.GVA_DBList = dbMap
}
