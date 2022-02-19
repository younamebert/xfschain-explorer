package home

import (
	api "xfschainbrowser/api/index"

	"github.com/gin-gonic/gin"
)

type HomeRouterGroup struct{}

func (r *HomeRouterGroup) HomeRouters(Router *gin.RouterGroup) {
	group := Router.Group("/index")

	resources := api.IndexLinkApi{}
	group.GET("/status", resources.Status)
}
