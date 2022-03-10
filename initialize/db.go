package initialize

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"xfschainbrowser/common"
	"xfschainbrowser/conf"
	"xfschainbrowser/global"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//MysqlDb mysql结构体
func Gorm() *gorm.DB {
	db, err := gorm.Open("mysql", "root:1263701671@(127.0.0.1:3306)/xfschain?charset=utf8mb4&parseTime=True&loc=Local")
	// db.SetLogger(true)
	db.LogMode(false)
	if err != nil {
		fmt.Printf("gorm err:%v\n", err)
		os.Exit(1)
	}
	if err := installMysql(db); err != nil {
		fmt.Printf("installMysql err:%v\n", err)
		os.Exit(1)
	}
	db.DB().SetMaxIdleConns(50)
	db.DB().SetMaxOpenConns(50)
	db.DB().SetConnMaxLifetime(time.Minute)
	// defer db.Close()
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
