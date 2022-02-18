package model

type ChainAddress struct {
	Basics
	Id                    int64   `gorm:"column:block_number"`
	Address               string  `gorm:"column:address"`
	Balance               float64 `gorm:"column:balance"`
	Nonce                 float64 `gorm:"column:nonce"`
	Extra                 string  `gorm:"column:extra"`
	Code                  string  `gorm:"column:code"`
	StateRoot             string  `gorm:"column:state_root"`
	Alias                 string  `gorm:"column:alias"`
	Type                  int     `gorm:"column:type"`
	Display               int     `gorm:"column:display"`
	FromStateRoot         string  `gorm:"column:from_state_root"`
	FromBlockHeight       int64   `gorm:"column:from_block_height"`
	FromBlockHash         string  `gorm:"column:from_block_hash"`
	CreateFromAddress     string  `gorm:"column:create_from_address"`
	CreateFromBlockHeight int64   `gorm:"column:create_from_block_height"`
	CreateFromBlockHash   string  `gorm:"column:create_from_block_hash"`
	CreateFromStateRoot   string  `gorm:"column:create_from_state_root"`
	CreateFromTxHash      string  `gorm:"column:create_from_tx_hash"`
}
