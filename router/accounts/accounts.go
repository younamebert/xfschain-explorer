package accounts

import (
	api "xfschainbrowser/api/account"
	"xfschainbrowser/common/apis"

	"github.com/gin-gonic/gin"
)

type AccountsRouterGroup struct{}

func (r *AccountsRouterGroup) AccountsRouters(Router *gin.RouterGroup) {
	group := Router.Group("/accounts")

	// resources := api.BlocksLinkApi{
	// 	Handle: apis.NewLinkApi(),
	// }
	resources := api.AccountLinkApi{
		Handle: apis.NewLinkApi(),
	}
	group.GET("/getaccounts", resources.GetAccounts)
	group.GET("/detailed", resources.Detailed)
	group.GET("/detailedtxs", resources.DetailedTxs)
}
