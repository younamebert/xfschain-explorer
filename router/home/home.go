package home

import (
	api "xfschainbrowser/api/index"
	"xfschainbrowser/common/apis"

	"github.com/gin-gonic/gin"
)

type HomeRouterGroup struct{}

func (r *HomeRouterGroup) HomeRouters(Router *gin.RouterGroup) {
	group := Router.Group("/index")

	resources := api.IndexLinkApi{
		Handle: apis.NewLinkApi(),
	}

	group.GET("/status", resources.Status)
	group.GET("/search", resources.Search)
	group.GET("/lastest", resources.LatestBlocksAndTxs)
	group.GET("/txcountbyday", resources.TxCountByDay)
}
