package blocks

import (
	"github.com/gin-gonic/gin"
)

type BlocksRouterGroup struct{}

func (r *BlocksRouterGroup) BlocksRouters(Router *gin.RouterGroup) {
	// group := Router.Group("/blocks")

	// // resources := api.BlocksLinkApi{
	// // 	Handle: apis.NewLinkApi(),
	// // }
	// resources := api.BlocksLinkApi{
	// 	Handle: apis.NewLinkApi(),
	// }
	// group.GET("/getblocks", resources.GetBlocks)
	// group.GET("/detailed", resources.Detailed)
	// group.GET("/detailedtx", resources.DetailedTx)
}
