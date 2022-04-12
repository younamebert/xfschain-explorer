package router

import (
	"mi/router/equipment"
	"mi/router/payment"
)

type RouterGroup struct {
	EquipmentRouter equipment.EquipmentRouterGroup
	PaymentRouter   payment.PaymentRouterGroup
}

var RouterGroupApp = new(RouterGroup)
