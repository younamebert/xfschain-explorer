package model

import "xfschainbrowser/global"

type HandleChainAddress struct{}

type ChainAddress struct {
	Basics
	Id                    int64  `gorm:"column:block_number"`
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
	db := global.GVA_DB.Table("chain_block_header")

	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (hanle *HandleChainAddress) Select(addr string) *ChainAddress {
	db := global.GVA_DB.Table("chain_block_header")

	addrChain := new(ChainAddress)
	db.Where("address = ?", addr).First(&addrChain)
	return addrChain
}
