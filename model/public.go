package model

import "time"

type Basics struct {
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func NewBasics() *Basics {
	return &Basics{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}
