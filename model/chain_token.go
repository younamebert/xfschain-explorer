package model

type ChainToken struct {
	Basics
	Id              int64   `gorm:"column:id"`
	Name            string  `gorm:"column:name"`
	Symbol          string  `gorm:"column:symbol"`
	TotalSypply     float64 `gorm:"column:total_sypply"`
	Decimals        int     `gorm:"column:decimals"`
	Address         string  `gorm:"column:address"`
	Creator         string  `gorm:"column:creator"`
	TxCount         int     `gorm:"column:tx_count"`
	FromTxHash      string  `gorm:"column:from_tx_hash"`
	FromBlockHeight int     `gorm:"column:from_block_height"`
	FromBlockHash   string  `gorm:"column:from_block_hash"`
	FromStateRoot   string  `gorm:"column:from_state_root"`
}
