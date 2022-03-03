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
	Rewards          string  `gorm:"column:rewards"`
}

func (handle *HandleChainBlockHeader) Insert(data *ChainBlockHeader) error {
	// data.CreateTime = time.Now()
	db := global.GVA_DB.Table("chain_block_header")

	if err := db.Create(&data).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return err
	}
	return nil
}

func (handle *HandleChainBlockHeader) Query(query, args interface{}) []*ChainBlockHeader {
	db := global.GVA_DB.Table("chain_block_header")

	chainBlockHeaders := make([]*ChainBlockHeader, 0)
	if err := db.Where(query, args).Order("height desc").Find(&chainBlockHeaders).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}
	if len(chainBlockHeaders) == 0 {
		return nil
	}
	return chainBlockHeaders
}

func (handle *HandleChainBlockHeader) QuerySort(limit int64, order string) []*ChainBlockHeader {
	db := global.GVA_DB.Table("chain_block_header")
	chainBlockHeaders := make([]*ChainBlockHeader, limit)
	if err := db.Limit(limit).Order("height desc").Find(&chainBlockHeaders).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}
	return chainBlockHeaders
}

// func (handle *HandleChainBlockHeader) QueryDown(limit int64) []*ChainBlockHeader {
// 	db := global.GVA_DB.Table("chain_block_header")
// 	chainBlockHeaders := make([]*ChainBlockHeader, limit)
// 	if err := db.Limit(limit).Order("height asc").Find(&chainBlockHeaders).Error; err != nil {
// 		global.GVA_LOG.Error(err.Error())
// 		return nil
// 	}
// 	return chainBlockHeaders
// }

// func (handle *HandleChainBlockHeader) QueryBlockHeadersByTime(startTime int64) []*ChainBlockHeader {
// 	db := global.GVA_DB.Table("chain_block_header")

// 	chainBlockHeaders := make([]*ChainBlockHeader, 0)

// 	if err := db.Where("timestamp > ?", startTime).Order("height desc").Find(&chainBlockHeaders).Error; err != nil {
// 		global.GVA_LOG.Error(err.Error())
// 		return nil
// 	}
// 	return chainBlockHeaders
// }

func (handle *HandleChainBlockHeader) QueryTxCountSumByTime(startTime int64) int64 {
	db := global.GVA_DB.Table("chain_block_header")

	var result []int64
	var sum int64
	if err := db.Pluck("tx_count", &result).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return 0
	} else {
		for _, v := range result {
			sum += v
		}
		return sum
	}

}

func (handle *HandleChainBlockHeader) GetBlocks(query, args interface{}, page, pageSize int) []*ChainBlockHeader {
	db := global.GVA_DB.Table("chain_block_header")

	chainBlockHeaders := make([]*ChainBlockHeader, pageSize)
	if query != nil && args != nil {
		db = db.Where(query, args)
	}
	if err := db.Limit(pageSize).Offset((page - 1) * pageSize).Order("height desc").Find(&chainBlockHeaders).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}
	return chainBlockHeaders
}

func (handle *HandleChainBlockHeader) QueryLike(query interface{}, where interface{}) []*ChainBlockHeader {
	chainBlockHeaders := make([]*ChainBlockHeader, 0)
	db := global.GVA_DB.Table("chain_block_header")

	if err := db.Where(query, where).Order("height desc").Find(&chainBlockHeaders).Error; err != nil {
		global.GVA_LOG.Error(err.Error())
		return nil
	}
	return chainBlockHeaders
}

func (handle *HandleChainBlockHeader) Count(query, args interface{}) int64 {
	db := global.GVA_DB.Table("chain_block_header")
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
