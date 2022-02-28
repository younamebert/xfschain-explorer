package model

import (
	"time"
	"xfschainbrowser/global"
)

// type  HandleChainBlockTxExternal interface{
// 	Insert(data *ChainBlockTx) error
// }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       HandlerChainBlockTx struct{}
type HandleChainBlockTx struct {
}

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
		global.GVA_LOG.Error(err.Error())
		return err
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

func (hanle *HandleChainBlockTx) QueryLastBlockTxs(limit int64) []*ChainBlockTx {
	db := global.GVA_DB.Table("chain_block_tx")

	ChainBlockTxs := make([]*ChainBlockTx, limit)
	if err := db.Limit(limit).Order("block_height desc,nonce desc").Find(&ChainBlockTxs).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}

	return ChainBlockTxs
}

func (hanle *HandleChainBlockTx) GetTxs(query, args interface{}, page, pageSize int) []*ChainBlockTx {
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

func (hanl *HandleChainBlockTx) QueryLikeTx(query interface{}, where []interface{}) []*ChainBlockTx {
	ChainBlockTxs := make([]*ChainBlockTx, 0)
	db := global.GVA_DB.Table("chain_block_tx")

	if err := db.Where(query, where...).Order("block_height desc").Find(&ChainBlockTxs).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}
	return ChainBlockTxs
}

func (hanle *HandleChainBlockTx) Count() int64 {
	db := global.GVA_DB.Table("chain_block_tx")
	var count int64
	if err := db.Count(&count).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return 0
	}
	return count
}
