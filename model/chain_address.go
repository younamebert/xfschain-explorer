package model

import (
	"time"
	"xfschainbrowser/global"
)

// type HandleChainAddressExternal interface {
// 	Insert(data *ChainAddress) error
// 	Query(addr string) *ChainAddress
// 	Update(target *ChainAddress) error
// }

type HandleChainAddress struct{}

type ChainAddress struct {
	Basics
	Id                    int64  `gorm:"column:id"`
	Address               string `gorm:"column:address"`
	Balance               string `gorm:"column:balance"`
	Nonce                 int64  `gorm:"column:nonce"`
	Extra                 string `gorm:"column:extra"`
	Code                  string `gorm:"column:code"`
	StateRoot             string `gorm:"column:state_root"`
	Alias                 string `gorm:"column:alias"`
	Type                  int    `gorm:"column:type"`
	Display               int    `gorm:"column:display"`
	FromStateRoot         string `gorm:"column:from_state_root"`
	FromBlockHeight       int64  `gorm:"column:from_block_height"`
	FromBlockHash         string `gorm:"column:from_block_hash"`
	CreateFromAddress     string `gorm:"column:create_from_address"`
	CreateFromBlockHeight int64  `gorm:"column:create_from_block_height"`
	CreateFromBlockHash   string `gorm:"column:create_from_block_hash"`
	CreateFromStateRoot   string `gorm:"column:create_from_state_root"`
	CreateFromTxHash      string `gorm:"column:create_from_tx_hash"`
}

func (handle *HandleChainAddress) Insert(data *ChainAddress) error {
	// db := global.GVA_DB.Table("chain_address")
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if err := global.GVA_DB.Create(&data).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return err
	}
	return nil
}

func (hanle *HandleChainAddress) Query(addr string) *ChainAddress {
	db := global.GVA_DB.Table("chain_address")

	addrChain := new(ChainAddress)
	if err := db.Where("address = ?", addr).First(&addrChain).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return nil
	}
	return addrChain
}

func (handle *HandleChainAddress) Update(target *ChainAddress) error {
	db := global.GVA_DB.Table("chain_address")
	target.UpdateTime = time.Now()
	if err := db.Where("address = ?", target.Address).Updates(&target).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return err
	}
	return nil
}

func (hanle *HandleChainAddress) Count() int64 {
	db := global.GVA_DB.Table("chain_address")
	var count int64
	if err := db.Count(&count).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return 0
	}
	return count
}
