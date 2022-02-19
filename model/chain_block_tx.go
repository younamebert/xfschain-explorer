package model

import "xfschainbrowser/global"

type HandlerChainBlockTx struct{}

type ChainBlockTx struct {
	Basics
	Id          int64  `gorm:"column:id"`
	BlockHash   string `gorm:"column:block_hash"`
	BlockHeight int64  `gorm:"column:block_height"`
	BlockTime   int64  `gorm:"column:block_time"`
	Version     int    `gorm:"column:version"`
	From        string `gorm:"column:from"`
	To          string `gorm:"column:to"`
	GasPrice    string `gorm:"column:gas_price"`
	GasLimit    string `gorm:"column:gas_limit"`
	GasUsed     string `gorm:"column:gas_used"`
	GasFee      string `gorm:"column:gas_fee"`
	Data        string `gorm:"column:data"`
	Nonce       int64  `gorm:"column:nonce"`
	Value       string `gorm:"column:value"`
	Signature   string `gorm:"column:signature"`
	Hash        string `gorm:"column:hash"`
	Status      int    `gorm:"column:status"`
	Type        int    `gorm:"column:type"`
}

func (handle *HandlerChainBlockTx) Insert(data *ChainBlockTx) error {
	db := global.GVA_DB.Table("chain_block_tx")
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}
