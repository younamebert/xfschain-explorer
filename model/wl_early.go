package model

import (
	"mi/global"
	"time"
)

type HandleWlEarly struct{}

type WlEarly struct {
	Basics
	EarlyId   int64  `gorm:"column:early_id"`
	MangeId   int    `gorm:"column:mange_id"`
	Number    string `gorm:"column:number"`
	Type      int    `gorm:"column:type"`
	CStatus   int    `gorm:"c_status"`
	Warehouse int    `gorm:"column:warehouse"`
	Iccid     string `gorm:"column:iccid"`
}

func (h *HandleWlEarly) Insert(mie *WlEarly) error {
	mie.CreateTime = time.Now()
	mie.UpdateTime = time.Now()
	db := global.GVA_DB.Table("wl_early")
	if err := db.Create(&mie).Error; err != nil {
		return err
	}
	return nil
}

//单条数据查询
func (h *HandleWlEarly) QueryOne(query, args, querys, args2 interface{}) *WlEarly {
	db := global.GVA_DB.Table("wl_early")

	wlEarly := new(WlEarly)

	if err := db.Where(query, args).Where(querys, args2).Find(&wlEarly); err != nil {
		return wlEarly
	}
	return nil
}

//修改
func (h *HandleWlEarly) SaveEarly(query, args, querys, args2 interface{}, mie *WlEarly) error {
	db := global.GVA_DB.Table("wl_early")
	mie.UpdateTime = time.Now()
	if err := db.Where(query, args).Where(querys, args2).Update(&mie).Error; err != nil {
		return err
	}
	return nil
}
