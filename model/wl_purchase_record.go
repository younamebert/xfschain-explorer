package model

import (
	"fmt"
	"mi/global"
	"time"
)

type HandleWlPurchaseRecord struct{}

type WlPurchaseRecord struct {
	Basics
	PurchaseRecordId int64     `gorm:"column:purchase_record_id"`
	Price            float64   `gorm:"price"`
	OrderNo          string    `gorm:"order_no"`
	PayTime          time.Time `gorm:"pay_time"`
	MemberId         int       `gorm:"member_id"`
	Iccid            string    `gorm:"iccid"`
	MangeId          int64     `gorm:"mange_id"`
}

func (h *HandleWlPurchaseRecord) Insert(mie *WlPurchaseRecord) error {
	mie.CreateTime = time.Now()
	mie.UpdateTime = time.Now()
	db := global.GVA_DB.Table("wl_purchase_record")
	if err := db.Create(&mie).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
