package payment

import (
	api "mi/api/paynotify"
	"mi/common/apis"

	"github.com/gin-gonic/gin"
)

type PaymentRouterGroup struct{}

func (r *PaymentRouterGroup) PaymentRouters(Router *gin.RouterGroup) {
	group := Router.Group("/pay")

	resources := api.PayLinkApi{
		Handle: apis.NewLinkApi(),
	}

	group.POST("/wechatnotify", resources.WeChatNotify)
	// group.GET("/detailed", resources.Detailed)
}
