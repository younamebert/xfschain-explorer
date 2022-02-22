package chainsync

//有hash
type BlockHeader struct {
	Height        int64  `json:"height"`
	Version       int64  `json:"version"`
	HashPrevBlock string `json:"hash_prev_block"`
	Timestamp     int64  `json:"timestamp"`
	Coinbase      string `json:"coinbase"`
	// merkle tree root hash
	StateRoot        string `json:"state_root"`
	TransactionsRoot string `json:"transactions_root"`
	ReceiptsRoot     string `json:"receipts_root"`
	GasLimit         int64  `json:"gas_limit"`
	GasUsed          int64  `json:"gas_used"`
	// pow
	Bits       int64   `json:"bits"`
	Nonce      int64   `json:"nonce"`
	ExtraNonce float64 `json:"extranonce"`
	Hash       string  `json:"hash"`
}

type AccountState struct {
	Balance   string `json:"balance"`
	Nonce     int64  `json:"nonce"`
	Extra     string `json:"extra"`
	Code      string `json:"code"`
	StateRoot string `json:"state_root"`
}

type Block struct {
	Header       *BlockHeader   `json:"header"`
	Transactions []*Transaction `json:"transactions"`
	Receipts     []*Receipt     `json:"receipts"`
}

type Receipt struct {
	Version int64  `json:"version"`
	Status  int64  `json:"status"`
	TxHash  string `json:"tx_hash"`
	GasUsed int64  `json:"gas_used"`
}

type Transaction struct {
	Version   int64   `json:"version"`
	To        string  `json:"to"`
	GasPrice  float64 `json:"gas_price"`
	GasLimit  float64 `json:"gas_limit"`
	Data      string  `json:"data"`
	Nonce     int64   `json:"nonce"`
	Value     int64   `json:"value"`
	Signature string  `json:"signature"`
	From      string  `json:"from"`
	Hash      string  `json:"hash"`
}
