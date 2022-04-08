package pays

import (
	api "mi/api/payment"
	"mi/common/apis"

	"github.com/gin-gonic/gin"
)

type PayRouterGroup struct{}

func (r *PayRouterGroup) PayRouters(Router *gin.RouterGroup) {
	group := Router.Group("/pay")

	resources := api.PayLinkApi{
		Handle: apis.NewLinkApi(),
	}

	group.GET("/payment", resources.Payment)
}
