package model

import (
	"xfschainbrowser/global"
)

// type HandleChainBlockHeaderExternal interface {
// 	Insert(data *ChainBlockHeader) error
// }

type HandleChainBlockHeader struct{}

type ChainBlockHeader struct {
	Id               int64   `gorm:"column:id"`
	Height           int64   `gorm:"column:height"`
	Hash             string  `gorm:"column:hash"`
	Version          int64   `gorm:"column:version"`
	HashPrevBlock    string  `gorm:"column:hash_prev_block"`
	Timestamp        int64   `gorm:"column:timestamp"`
	Coinbase         string  `gorm:"column:coinbase"`
	StateRoot        string  `gorm:"column:state_root"`
	TransactionsRoot string  `gorm:"column:transactions_root"`
	ReceiptsRoot     string  `gorm:"column:receipts_root"`
	GasLimit         int64   `gorm:"column:gas_limit"`
	GasUsed          int64   `gorm:"column:gas_used"`
	Bits             int64   `gorm:"column:bits"`
	Nonce            int64   `gorm:"column:nonce"`
	ExtraNonce       float64 `gorm:"extra_nonce"`
	TxCount          int     `gorm:"column:tx_count"`
	Rewards          float64 `gorm:"column:rewards"`
}

func (handle *HandleChainBlockHeader) Insert(data *ChainBlockHeader) error {
	db := global.GVA_DB.Table("chain_block_header")

	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}
