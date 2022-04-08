package payment

import (
	"fmt"
	"mi/common/apis"
	"mi/pay"

	"github.com/gin-gonic/gin"
)

type PayLinkApi struct {
	Handle *apis.LinkApi
}

func (ac *PayLinkApi) Payment(c *gin.Context) {
	pay.Payment()
	fmt.Println("支付页面")
}
