package transfer

import (
	api "xfschainbrowser/api/transfer"
	"xfschainbrowser/common/apis"

	"github.com/gin-gonic/gin"
)

type TxsRouterGroup struct{}

func (r *TxsRouterGroup) TxsRouters(Router *gin.RouterGroup) {
	group := Router.Group("/transfer")

	resources := api.TxsLinkApi{
		Handle: apis.NewLinkApi(),
	}

	group.GET("/gettxs", resources.GetTxs)
	group.GET("/detailed", resources.Detailed)
}
