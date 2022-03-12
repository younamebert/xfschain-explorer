package model

import "github.com/shopspring/decimal"

type HandleMiOrder struct{}

type MiOrder struct {
	Basics
	Id            int64           `gorm:"column:id"`
	Iccid         string          `gorm:"iccid"`
	PayType       int             `gorm:"pay_type"`
	PayCode       string          `gorm:"pay_Price"`
	OrderNumber   string          `gorm:"order_number"`
	PaymentAmount decimal.Decimal `gorm:"payment_amount"`
}
