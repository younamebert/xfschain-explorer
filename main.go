package main

import (
	"xfschainbrowser/core"
	"xfschainbrowser/global"
	"xfschainbrowser/initialize"
)

func main() {
	global.GVA_LOG = core.Zap()       // 初始化zap日志库
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	// core.chainSyncCore()
	core.RunServer()
}
