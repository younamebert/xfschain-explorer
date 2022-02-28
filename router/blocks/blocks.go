package blocks

import (
	api "xfschainbrowser/api/blockchain"
	"xfschainbrowser/common/apis"

	"github.com/gin-gonic/gin"
)

type BlocksRouterGroup struct{}

func (r *BlocksRouterGroup) BlocksRouters(Router *gin.RouterGroup) {
	group := Router.Group("/blocks")

	// resources := api.BlocksLinkApi{
	// 	Handle: apis.NewLinkApi(),
	// }
	resources := api.BlocksLinkApi{
		Handle: apis.NewLinkApi(),
	}
	group.GET("/getblocks", resources.GetBlocks)
}
