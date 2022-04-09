package main

import (
	"mi/core"
	"mi/global"
	"mi/initialize"
)

func main() {
	global.GVA_LOG = core.Zap()       // 初始化zap日志库
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	initialize.DBList()
	if global.GVA_DB != nil {
		db := global.GVA_DB.DB()
		defer db.Close()
	}
	core.RunServer()
}
