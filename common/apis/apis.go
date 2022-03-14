package apis

import (
	"mi/events"
	"mi/model"
)

type LinkApi struct {
	HandleMiEquipment *model.HandleMiEquipment
	HandleMiWarehouse *model.HandleMiWarehouse
	HandleMiOrder     *model.HandleMiOrder
	EventsBus         *events.EventBus
}

func NewLinkApi() *LinkApi {
	return &LinkApi{
		HandleMiEquipment: new(model.HandleMiEquipment),
		HandleMiWarehouse: new(model.HandleMiWarehouse),
		HandleMiOrder:     new(model.HandleMiOrder),
		EventsBus:         events.EventBusExample,
	}
}

var ApiResource = NewLinkApi()
