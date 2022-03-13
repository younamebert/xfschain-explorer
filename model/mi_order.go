package model

import (
	"mi/global"
	"time"

	"github.com/shopspring/decimal"
)

type HandleMiOrder struct{}

type MiOrder struct {
	Basics
	Id            int64           `gorm:"column:id"`
	Iccid         string          `gorm:"iccid"`
	PayType       int             `gorm:"pay_type"`
	PayCode       string          `gorm:"pay_Price"`
	OrderNumber   string          `gorm:"order_number"`
	PaymentAmount decimal.Decimal `gorm:"payment_amount"`
	Number        string          `gorm:"number"`
}

func (h *HandleMiOrder) Insert(mie *MiOrder) error {
	mie.CreateTime = time.Now()
	mie.UpdateTime = time.Now()
	db := global.GVA_DB.Table("mi_order")
	if err := db.Create(&mie).Error; err != nil {
		return err
	}
	return nil
}

func (h *HandleMiOrder) Count(query, args interface{}) int64 {
	db := global.GVA_DB.Table("mi_order")
	var count int64 = 0
	if query != nil && args != nil {
		db.Where(query, args)
	}

	if err := db.Count(&count).Error; err != nil {
		global.GVA_LOG.Warn(err.Error())
		return count
	}
	return count
}

func (h *HandleMiOrder) Querys(query, args interface{}, page, pageSize int) []*MiEquipment {
	db := global.GVA_DB.Table("mi_order")
	miEquipments := make([]*MiEquipment, pageSize)
	if query != nil && args != nil {
		db.Where(query, args)
	}

	if err := db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&miEquipments).Error; err != nil {
		global.GVA_LOG.Warn(err.Error())
		return nil
	}

	return miEquipments
}

func (h *HandleMiOrder) SetSwitchad(iccid string, i int) error {
	db := global.GVA_DB.Table("mi_order")

	err := db.Where("iccid = ?", iccid).Update("switchad", i).Error
	if err != nil {
		return err
	}
	return err
}

func (h *HandleMiOrder) SetSwitchadLed(iccid string, i int) error {
	db := global.GVA_DB.Table("mi_order")

	err := db.Where("iccid = ?", iccid).Update("switch_lamp", i).Error
	if err != nil {
		return err
	}
	return err
}
