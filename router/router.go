package router

import (
	"mi/router/equipment"
)

type RouterGroup struct {
	//
	EquipmentRouter equipment.EquipmentRouterGroup
}

var RouterGroupApp = new(RouterGroup)
