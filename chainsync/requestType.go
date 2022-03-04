package chainsync

type getBlockByHashArgs struct {
	Hash string `json:"hash"`
}

type GetBlockHeaderByNumberArgs struct {
	Number string `json:"number"`
	//Count  string `json:"count"`
}

type getAccountByAddrArgs struct {
	RootHash string `json:"root_hash"`
	Address  string `json:"address"`
}

type getTxByHashArgs struct {
	TxHash string `json:"hash"`
}
