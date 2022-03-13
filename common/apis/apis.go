package apis

import "mi/model"

type LinkApi struct {
	HandleMiEquipment *model.HandleMiEquipment
	HandleMiWarehouse *model.HandleMiWarehouse
	HandleMiOrder     *model.HandleMiOrder
}

func NewLinkApi() *LinkApi {
	return &LinkApi{
		HandleMiEquipment: new(model.HandleMiEquipment),
		HandleMiWarehouse: new(model.HandleMiWarehouse),
		HandleMiOrder:     new(model.HandleMiOrder),
	}
}

var ApiResource = NewLinkApi()
