package model

import (
	"mi/global"
	"time"
)

// type  HandleChainBlockTxExternal interface{
// 	Insert(data *ChainBlockTx) error
// }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       HandlerChainBlockTx struct{}
type HandleMiWarehouse struct {
}

type MiWarehouse struct {
	Basics
	Id             int64   `gorm:"column:id"`
	Iccid          string  `gorm:"column:iccid"`
	WarehouseType  int64   `gorm:"column:warehouse_type"`
	WarehousePrice float64 `gorm:"column:warehouse_price"`
}

func (h *HandleMiWarehouse) Insert(mie *MiWarehouse) error {
	mie.CreateTime = time.Now()
	mie.UpdateTime = time.Now()
	db := global.GVA_DB.Table("mi_warehouse")
	if err := db.Create(&mie).Error; err != nil {
		return err
	}
	return nil
}

func (h *HandleMiWarehouse) Update(mie *MiWarehouse) error {
	mie.CreateTime = time.Now()
	mie.UpdateTime = time.Now()
	db := global.GVA_DB.Table("mi_warehouse")
	if err := db.Update(&mie).Error; err != nil {
		return err
	}
	return nil
}

func (h *HandleMiWarehouse) Query(query, args interface{}) []*MiWarehouse {
	db := global.GVA_DB.Table("mi_warehouse")
	miWarehouses := make([]*MiWarehouse, 0)
	if err := db.Where(query, args).Find(&miWarehouses).Error; err != nil {
		return miWarehouses
	}
	return miWarehouses
}

func (h *HandleMiWarehouse) Querys(query, args interface{}, page, pageSize int) []*MiWarehouse {
	db := global.GVA_DB.Table("mi_warehouse")
	miWarehouses := make([]*MiWarehouse, pageSize)
	if query != nil && args != nil {
		db.Where(query, args)
	}

	if err := db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&miWarehouses).Error; err != nil {
		global.GVA_LOG.Warn(err.Error())
		return nil
	}
	return miWarehouses
}
