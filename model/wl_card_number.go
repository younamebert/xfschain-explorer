package model

import (
	"mi/global"
	"time"
)

type HandleWlCardNumber struct {
}

type WlCardNumber struct {
	CardNumberId int     `gorm:"column:card_number_id"`
	Number       string  `gorm:"column:number"`
	Money        float64 `gorm:"column:money"`
	MemberId     int     `gorm:"column:member_id"`
	//
	Basics
}

//查
func (h *HandleWlCardNumber) Query(query, args interface{}) *WlCardNumber {
	db := global.GVA_DB.Table("wl_card_number")

	// wlCardNumber := make([]*WlCardNumber, 0)

	wlCardNumber := new(WlCardNumber)

	if err := db.Where(query, args).Take(&wlCardNumber).Error; err != nil {
		return wlCardNumber
	}
	return wlCardNumber
}

func (h *HandleWlCardNumber) Update(query, args interface{}, mie *WlCardNumber) error {
	db := global.GVA_DB.Table("wl_card_number")
	mie.UpdateTime = time.Now()
	if err := db.Where(query, args).Update(&mie).Error; err != nil {
		return err
	}
	return nil
}
