package model

import (
	"mi/global"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type HandleWlCardNumber struct {
}

type WlCardNumber struct {
	CardNumberId int             `gorm:"column:card_number_id"`
	Number       string          `gorm:"column:number"`
	Balance      decimal.Decimal `gorm:"column:balance"`
	MemberId     int             `gorm:"column:member_id"`
	Basics
}

//查
func (h *HandleWlCardNumber) Query(query, args interface{}) *WlCardNumber {
	db := global.GVA_DBList["wl_card_number"]
	wlCardNumber := new(WlCardNumber)
	if err := db.Where(query, args).Take(&wlCardNumber).Error; err != nil {
		return nil
	}
	return wlCardNumber
}

func (h *HandleWlCardNumber) Update(query, args interface{}, mie *WlCardNumber) error {
	//开启事物在修改
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txsubmit := tx.Table("wl_card_number")
		if query != nil {
			txsubmit.Where(query, args)
		}
		if err := txsubmit.Update(&mie).Error; err != nil {
			return err
		}
		return nil
	})
}
