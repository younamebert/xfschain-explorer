package api

import "xfschainbrowser/model"

type StatusResp struct {
	LatestHeight int64  `json:"latestHeight"`
	Accounts     int64  `json:"accounts"`
	BlockRewards string `json:"blockRewards"`
	BlockTime    int64  `json:"blockTime"`
	Difficulty   int64  `json:"difficulty"`
	Power        int64  `json:"power"`
	Tps          string `json:"tps"`
	Transactions int64  `json:"transactions"`
	TxsInBlock   int64  `json:"txsInBlock"`
}

type LatestResp struct {
	Blocks []*model.ChainBlockHeader
	Txs    []*model.ChainBlockTx
}

type TxCountByDayResp struct {
	Timestamp int64 `json:"timestamp"`
	TxCount   int64 `json:"txcount"`
}
