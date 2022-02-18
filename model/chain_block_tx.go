package model

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
	GasFee      float64 `gorm:"column:gas_fee"`
	Data        string  `gorm:"column:data"`
	Nonce       int64   `gorm:"column:nonce"`
	Value       float64 `gorm:"column:value"`
	Signature   string  `gorm:"column:signature"`
	Hash        string  `gorm:"column:hash"`
	Status      int     `gorm:"column:status"`
	Type        int     `gorm:"column:type"`
}
