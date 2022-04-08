package router

import (
	"mi/router/equipment"
	"mi/router/pays"
)

type RouterGroup struct {
	EquipmentRouter equipment.EquipmentRouterGroup
}

type RoutersGroup struct {
	PayRouterGroup pays.PayRouterGroup
}

var RouterGroupApp = new(RouterGroup)
var RouterGroupApps = new(RoutersGroup)
