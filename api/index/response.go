package api

import "xfschainbrowser/model"

type StatusResp struct {
	LatestHeight int64   `json:"latestHeight"`
	Accounts     int64   `json:"accounts"`
	BlockRewards string  `json:"blockRewards"`
	BlockTime    float64 `json:"blockTime"`
	Difficulty   int64   `json:"difficulty"`
	Power        int64   `json:"power"`
	Tps          float64 `json:"tps"`
	Transactions int64   `json:"transactions"`
	TxsInBlock   int64   `json:"txsInBlock"`
}

type LatestResp struct {
	Blocks []*model.ChainBlockHeader `json:"blocks"`
	Txs    []*model.ChainBlockTx     `json:"txs"`
}

type TxCountByDayResp struct {
	Timestamp int64 `json:"timestamp"`
	TxCount   int64 `json:"txcount"`
}

type SearchResp struct {
	Type      int    `json:"type"`
	PathValue string `json:"pathValue"`
}

type TxCountByDayResps []*TxCountByDayResp

func (s TxCountByDayResps) Len() int           { return len(s) }
func (s TxCountByDayResps) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TxCountByDayResps) Less(i, j int) bool { return s[i].Timestamp < s[j].Timestamp }
