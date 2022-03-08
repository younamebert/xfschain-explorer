package chainsync

import (
	"xfschainbrowser/chainsync/httpxfs"
	"xfschainbrowser/common"
	"xfschainbrowser/conf"
	"xfschainbrowser/global"

	"go.uber.org/zap"
)

type chainMgr interface {
	CurrentBHeader() *BlockHeader
	GetReceiptByHash(txhash string) *Receipt
	GetTxsByBlockHash(blockHash string) []*Transaction
	GetBlockHeaderByHash(blockHash string) *BlockHeader
	GetAccountInfo(address string) *AccountState
	GetBlockHeaderByNumber(blocknumber string) *BlockHeader
	CheckResponse() error
}

type syncMgr struct {
	xfsClient *httpxfs.Client
}

func newsyncMgr() *syncMgr {
	cli := httpxfs.NewClient(conf.SyncRquesetURL, conf.SyncTimeoutBlockFetch)
	return &syncMgr{
		xfsClient: cli,
	}
}

func (syncMgr *syncMgr) CheckResponse() error {
	lastBlockHeader := new(BlockHeader)
	return syncMgr.xfsClient.CallMethod(1, "Chain.Head", nil, &lastBlockHeader)
}

// func (syncMgr *syncMgr)
//CurrentBHeader 获取最新的高度
func (syncMgr *syncMgr) CurrentBHeader() *BlockHeader {
	lastBlockHeader := new(BlockHeader)
	if err := syncMgr.xfsClient.CallMethod(1, "Chain.Head", nil, &lastBlockHeader); err != nil {
		global.GVA_LOG.Warn("code:"+common.SystemErr+" err:", zap.Any(" error:", err.Error()))
		return nil
	}
	return lastBlockHeader
}

//GetReceiptByHash 条件区块哈希
func (syncMgr *syncMgr) GetReceiptByHash(txhash string) *Receipt {
	req := &getTxByHashArgs{
		TxHash: txhash,
	}
	recs := new(Receipt)
	if err := syncMgr.xfsClient.CallMethod(1, "Chain.GetReceiptByHash", &req, &recs); err != nil {
		global.GVA_LOG.Warn("code:"+common.SystemErr+" err:", zap.Any(" error:", err.Error()))
		return nil
	}
	return recs
}

//GetTxsByBlockHash 区块哈希获取区块所有交易
func (syncMgr *syncMgr) GetTxsByBlockHash(blockHash string) []*Transaction {
	req := &getBlockByHashArgs{
		Hash: blockHash,
	}
	txs := make([]*Transaction, 0)
	if err := syncMgr.xfsClient.CallMethod(1, "Chain.GetTxsByBlockHash", &req, &txs); err != nil {
		global.GVA_LOG.Warn("code:"+common.SystemErr+" err:", zap.Any(" error:", err.Error()))
		return nil
	}
	return txs
}

//GetBlockHeaderByNumber 区块number获取区块头部详情
func (syncMgr *syncMgr) GetBlockHeaderByNumber(blocknumber string) *BlockHeader {
	req := &GetBlockHeaderByNumberArgs{
		Number: blocknumber,
	}
	rets := new(BlockHeader)
	if err := syncMgr.xfsClient.CallMethod(1, "Chain.GetBlockHeaderByNumber", &req, &rets); err != nil {
		global.GVA_LOG.Warn("code:"+common.SystemErr+" err:", zap.Any(" error:", err.Error()))
		return nil
	}
	return rets
}

//GetBlockHeaderByHash 区块哈希获取区块头部详情
func (syncMgr *syncMgr) GetBlockHeaderByHash(blockHash string) *BlockHeader {
	req := &getBlockByHashArgs{
		Hash: blockHash,
	}
	rets := new(BlockHeader)
	if err := syncMgr.xfsClient.CallMethod(1, "Chain.GetBlockHeaderByHash", &req, &rets); err != nil {
		global.GVA_LOG.Warn("code:"+common.SystemErr+" err:", zap.Any(" error:", err.Error()))
		return nil
	}
	return rets
}

func (syncMgr *syncMgr) GetAccountInfo(addr string) *AccountState {
	req := &getAccountByAddrArgs{
		Address: addr,
	}
	rets := new(AccountState)
	if err := syncMgr.xfsClient.CallMethod(1, "State.GetAccount", &req, &rets); err != nil {
		return nil
	}
	return rets
}
