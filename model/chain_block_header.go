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
	// data.CreateTime = time.Now()
	db := global.GVA_DB.Table("chain_block_header")

	if err := db.Create(&data).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return err

	}
	return nil
}

func (hande *HandleChainBlockHeader) QueryByHash(hash string) *ChainBlockHeader {
	db := global.GVA_DB.Table("chain_block_header")

	chainBlockHeader := new(ChainBlockHeader)
	if err := db.Where("hash = ?", hash).First(&chainBlockHeader).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return nil
	}
	return chainBlockHeader
}

func (hanle *HandleChainBlockHeader) QueryUp(limit int64) []*ChainBlockHeader {
	db := global.GVA_DB.Table("chain_block_header")
	chainBlockHeaders := make([]*ChainBlockHeader, limit)
	if err := db.Limit(limit).Order("height desc").First(&chainBlockHeaders).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return nil
	}
	return chainBlockHeaders
}

func (hanle *HandleChainBlockHeader) QueryDown(limit int64) []*ChainBlockHeader {
	db := global.GVA_DB.Table("chain_block_header")
	chainBlockHeaders := make([]*ChainBlockHeader, limit)
	if err := db.Limit(limit).Order("height asc").First(&chainBlockHeaders).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return nil
	}
	return chainBlockHeaders
}

func (handle *HandleChainBlockHeader) QueryBlockHeadersByTime(startTime int64) []*ChainBlockHeader {
	db := global.GVA_DB.Table("chain_block_header")

	chainBlockHeaders := make([]*ChainBlockHeader, 0)

	if err := db.Where("timestamp > ?", startTime).Order("height desc").Find(&chainBlockHeaders).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return nil
	}
	return chainBlockHeaders
}

func (handle *HandleChainBlockHeader) QueryTxCountSumByTime(startTime int64) int64 {
	db := global.GVA_DB.Table("chain_block_header")

	var result []int64
	var sum int64
	if err := db.Pluck("tx_count", &result).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return 0
	} else {
		for _, v := range result {
			sum += v
		}
		return sum
	}

}

func (hanle *HandleChainBlockHeader) GetBlocks(page, pageSize, limit int) []*ChainBlockHeader {
	db := global.GVA_DB.Table("chain_block_header")

	chainBlockHeaders := make([]*ChainBlockHeader, limit)
	if err := db.Limit(limit).Offset((page - 1) * pageSize).Order("block_height desc,nonce desc").First(&chainBlockHeaders).Error; err != nil {
		global.GVA_LOG.Fatal(err.Error())
		return nil
	}
	return chainBlockHeaders
}
