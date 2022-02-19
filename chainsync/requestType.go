package chainsync

type getBlockByHashArgs struct {
	Hash string `json:"hash"`
}

type getAccountByAddrArgs struct {
	RootHash string `json:"root_hash"`
	Address  string `json:"address"`
}

type getTxByHashArgs struct {
	TxHash string `json:"hash"`
}
