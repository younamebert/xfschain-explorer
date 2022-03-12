package model

import (
	"mi/global"
	"time"
)

type HandleMiEquipment struct{}

type MiEquipment struct {
	Basics
	Id           int64  `gorm:"column:id"`
	Switchad     int    `gorm:"column:switchad"`
	SwitchadLamp int    `gorm:"column:switch_lamp"`
	Iccid        string `gorm:"column:iccid"`
}

func (h *HandleMiEquipment) Insert(mie *MiEquipment) error {
	mie.CreateTime = time.Now()
	mie.UpdateTime = time.Now()
	db := global.GVA_DB.Table("mi_equipment")
	if err := db.Create(&mie).Error; err != nil {
		return err
	}
	return nil
}

func (h *HandleMiEquipment) Count(query, args interface{}) int64 {
	db := global.GVA_DB.Table("mi_equipment")
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

func (h *HandleMiEquipment) Querys(query, args interface{}, page, pageSize int) []*MiEquipment {
	db := global.GVA_DB.Table("mi_equipment")
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

func (h *HandleMiEquipment) SetSwitchad(iccid string, i int) error {
	db := global.GVA_DB.Table("mi_equipment")

	err := db.Where("iccid = ?", iccid).Update("switchad", i).Error
	if err != nil {
		return err
	}
	return err
}

func (h *HandleMiEquipment) SetSwitchadLed(iccid string, i int) error {
	db := global.GVA_DB.Table("mi_equipment")

	err := db.Where("iccid = ?", iccid).Update("switch_lamp", i).Error
	if err != nil {
		return err
	}
	return err
}
