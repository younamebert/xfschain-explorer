package api

type StatusResp struct {
	LatestHeight int64  `json:"latestHeight"`
	Accounts     int64  `json:"accounts"`
	BlockRewards string `json:"blockRewards"`
	BlockTime    int64  `json:"blockTime"`
	Difficulty   int64  `json:"difficulty"`
	Power        int64  `json:"power"`
	Tps          int64  `json:"tps"`
	Transactions int64  `json:"transactions"`
	TxsInBlock   int64  `json:"txsInBlock"`
}
