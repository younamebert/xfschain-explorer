package model

import (
	"time"
	"xfschainbrowser/global"
)

// type  HandleChainBlockTxExternal interface{
// 	Insert(data *ChainBlockTx) error
// }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       HandlerChainBlockTx struct{}
type HandleChainBlockTx struct{}

type ChainBlockTx struct {
	Basics
	Id          int64   `gorm:"column:id"`
	BlockHash   string  `gorm:"column:block_hash"`
	BlockHeight int64   `gorm:"column:block_height"`
	BlockTime   int64   `gorm:"column:block_time"`
	Version     int     `gorm:"column:version"`
	From        string  `gorm:"column:from"`
	To          string  `gorm:"column:to"`
	GasPrice    float64 `gorm:"column:gas_price"`
	GasLimit    float64 `gorm:"column:gas_limit"`
	GasUsed     float64 `gorm:"column:gas_used"`
	GasFee      string  `gorm:"column:gas_fee"`
	Data        string  `gorm:"column:data"`
	Nonce       int64   `gorm:"column:nonce"`
	Value       string  `gorm:"column:value"`
	Signature   string  `gorm:"column:signature"`
	Hash        string  `gorm:"column:hash"`
	Status      int     `gorm:"column:status"`
	Type        int     `gorm:"column:type"`
}

func (handle *HandleChainBlockTx) Insert(data *ChainBlockTx) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	db := global.GVA_DB.Table("chain_block_tx")
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (handle *HandleChainBlockTx) QueryByBlockHash(blockHash string) *ChainBlockTx {
	db := global.GVA_DB.Table("chain_block_tx")

	ChainBlockTx := new(ChainBlockTx)
	if err := db.Where("block_hash = ?", blockHash).First(&ChainBlockTx).Error; err != nil {
		return nil
	}
	return ChainBlockTx
}
