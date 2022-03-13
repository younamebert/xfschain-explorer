package core

import (
	"mi/conf"
	"mi/global"
	"mi/initialize"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RunServer() {
	gin.SetMode(gin.ReleaseMode)
	Router := initialize.Routers()
	// address := fmt.Sprintf(":%d", conf.Addr)

	s := initServer(conf.Addr, Router)
	// time.Sleep(10 * time.Microsecond)
	// go chainSyncCore()
	go tcpServer()

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
