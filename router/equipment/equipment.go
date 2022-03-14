package equipment

import (
	api "mi/api/equipment"
	"mi/common/apis"

	"github.com/gin-gonic/gin"
)

type EquipmentRouterGroup struct{}

func (r *EquipmentRouterGroup) EquipmentRouters(Router *gin.RouterGroup) {
	group := Router.Group("/equipment")

	// resources := api.BlocksLinkApi{
	// Handle: apis.NewLinkApi(),
	// }
	// resources := api.AccountLinkApi{
	// Handle: apis.NewLinkApi(),
	// }

	resources := api.EquipmentLinkApi{
		Handle: apis.NewLinkApi(),
	}

	group.GET("/list", resources.EquipmentList)

	group.GET("/switchAdvertising", resources.SwitchAdvertising)

	group.GET("/switchLed", resources.SwitchLed)

	// group.GET("/getaccounts", resources.GetAccounts)
	// group.GET("/detailed", resources.Detailed)
	// group.GET("/detailedtxs", resources.DetailedTxs)
}
