package model

import (
	"time"
	"xfschainbrowser/common"
	"xfschainbrowser/global"

	"github.com/shopspring/decimal"
)

// type  HandleChainBlockTxExternal interface{
// 	Insert(data *ChainBlockTx) error
// }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       HandlerChainBlockTx struct{}
type HandleChainBlockTx struct {
}

type ChainBlockTx struct {
	Basics
	Id          int64           `gorm:"column:id"`
	BlockHash   string          `gorm:"column:block_hash"`
	BlockHeight int64           `gorm:"column:block_height"`
	BlockTime   int64           `gorm:"column:block_time"`
	Version     int             `gorm:"column:version"`
	TxFrom      string          `gorm:"column:tx_from"`
	TxTo        string          `gorm:"column:tx_to"`
	GasPrice    decimal.Decimal `gorm:"column:gas_price"`
	GasLimit    decimal.Decimal `gorm:"column:gas_limit"`
	GasUsed     decimal.Decimal `gorm:"column:gas_used"`
	GasFee      string          `gorm:"column:gas_fee"`
	Data        string          `gorm:"column:data"`
	Nonce       int64           `gorm:"column:nonce"`
	Value       string          `gorm:"column:value"`
	Signature   string          `gorm:"column:signature"`
	Hash        string          `gorm:"column:hash"`
	Status      int             `gorm:"column:status"`
	Type        int             `gorm:"column:type"`
}

func (handle *HandleChainBlockTx) Insert(data *ChainBlockTx) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	db := global.GVA_DB.Table("chain_block_tx")
	err := db.Create(&data).Error

	if err != nil {
		if common.ContainsErr(err.Error(), "Duplicate") {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (handle *HandleChainBlockTx) Query(query, args interface{}) []*ChainBlockTx {
	db := global.GVA_DB.Table("chain_block_tx")

	ChainBlockTxs := make([]*ChainBlockTx, 0)
	if err := db.Where(query, args).Find(&ChainBlockTxs).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}
	return ChainBlockTxs
}

func (handle *HandleChainBlockTx) QueryLastBlockTxs(limit int64) []*ChainBlockTx {
	db := global.GVA_DB.Table("chain_block_tx")

	ChainBlockTxs := make([]*ChainBlockTx, limit)
	if err := db.Limit(limit).Order("block_height desc,nonce desc").Find(&ChainBlockTxs).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}

	return ChainBlockTxs
}

func (handle *HandleChainBlockTx) GetTxs(query, args interface{}, page, pageSize int) []*ChainBlockTx {
	db := global.GVA_DB.Table("chain_block_tx")

	ChainBlockTxs := make([]*ChainBlockTx, pageSize)
	if query != nil && args != nil {
		db = db.Where(query, args)
	}
	if err := db.Limit(pageSize).Offset((page - 1) * pageSize).Order("block_height desc,nonce desc").Find(&ChainBlockTxs).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}

	return ChainBlockTxs
}

func (handle *HandleChainBlockTx) QueryLikeTx(query interface{}, where interface{}) []*ChainBlockTx {
	ChainBlockTxs := make([]*ChainBlockTx, 0)
	db := global.GVA_DB.Table("chain_block_tx")

	if err := db.Where(query, where).Order("block_height desc").Find(&ChainBlockTxs).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}
	return ChainBlockTxs
}

func (handle *HandleChainBlockTx) Count(query, args interface{}) int64 {
	db := global.GVA_DB.Table("chain_block_tx")
	var count int64
	if query != nil && args != nil {
		db = db.Where(query, args)
	}
	if err := db.Count(&count).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return 0
	}
	return count
}
