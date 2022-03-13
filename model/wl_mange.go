package model

import (
	"mi/global"
)

type HandleWlMange struct{}

type WlMange struct {
	Basics
	MangeId int64  `gorm:"column:mange_id"`
	Number  string `gorm:"number"`
}

func (h *HandleWlMange) Query(query, args interface{}) WlMange {
	db := global.GVA_DB.Table("wl_mange")
	var wlMange WlMange
	if err := db.Where(query, args).Find(&wlMange).Error; err != nil {
		return wlMange
	}
	return wlMange
}
