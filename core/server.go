package core

import (
	"net/http"
	"time"
	"xfschainbrowser/conf"
	"xfschainbrowser/global"
	"xfschainbrowser/initialize"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RunServer() {
	Router := initialize.Routers()
	// address := fmt.Sprintf(":%d", conf.Addr)
	s := initServer(conf.Addr, Router)
	time.Sleep(10 * time.Microsecond)
	// chainSyncCore()
	global.GVA_LOG.Info("server run success on ", zap.String("address", conf.Addr))
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}

func initServer(address string, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
