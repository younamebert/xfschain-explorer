package model

import (
	"mi/global"
	"time"
)

type HandleWlSale struct{}

type WlSale struct {
	Basics
	SaleId  int64   `gorm:"column:sale_id"`
	MangeId int64   `gorm:"mange_id"`
	Number  int     `gorm:"number"`
	Money   float64 `gorm:"money"`
}

func (h *HandleWlSale) Insert(mie *WlSale) error {
	mie.CreateTime = time.Now()
	mie.UpdateTime = time.Now()
	db := global.GVA_DB.Table("wl_sale")
	if err := db.Create(&mie).Error; err != nil {
		return err
	}
	return nil
}
