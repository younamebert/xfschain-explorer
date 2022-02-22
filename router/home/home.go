package home

import (
	api "xfschainbrowser/api/index"
	"xfschainbrowser/model"

	"github.com/gin-gonic/gin"
)

type HomeRouterGroup struct{}

func (r *HomeRouterGroup) HomeRouters(Router *gin.RouterGroup) {
	group := Router.Group("/index")

	resources := api.IndexLinkApi{
		HandleBlockHeader:  new(model.HandleChainBlockHeader),
		HandleBlockTxs:     new(model.HandleChainBlockTx),
		HandleChainAddress: new(model.HandleChainAddress),
	}
	group.GET("/status", resources.Status)
}
