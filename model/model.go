package model

type RecordHandle struct {
	HandleMiEquipment *HandleMiEquipment
	HandleMiWarehouse *HandleMiWarehouse
	HandleMiOrder     *HandleMiOrder
}

func NewRecordHandle() *RecordHandle {
	return &RecordHandle{
		HandleMiEquipment: new(HandleMiEquipment),
		HandleMiWarehouse: new(HandleMiWarehouse),
		HandleMiOrder:     new(HandleMiOrder),
	}
}
