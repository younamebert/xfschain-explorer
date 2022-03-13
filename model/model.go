package model

type RecordHandle struct {
	HandleMiEquipment      *HandleMiEquipment
	HandleMiWarehouse      *HandleMiWarehouse
	HandleMiOrder          *HandleMiOrder
	HandleWlCardNumber     *HandleWlCardNumber
	HandleWlSale           *HandleWlSale
	HandleWlPurchaseRecord *HandleWlPurchaseRecord
	HandleWlMange          *HandleWlMange
	HandleWlEarly          *HandleWlEarly
}

func NewRecordHandle() *RecordHandle {
	return &RecordHandle{
		HandleMiEquipment:      new(HandleMiEquipment),
		HandleMiWarehouse:      new(HandleMiWarehouse),
		HandleMiOrder:          new(HandleMiOrder),
		HandleWlCardNumber:     new(HandleWlCardNumber),
		HandleWlSale:           new(HandleWlSale),
		HandleWlPurchaseRecord: new(HandleWlPurchaseRecord),
		HandleWlMange:          new(HandleWlMange),
		HandleWlEarly:          new(HandleWlEarly),
	}
}
