package chainsync

import (
	"sync"
	"xfschainbrowser/model"
)

type syncService struct {
	syncMgr      *syncMgr
	recordHandle *recordHandle
	mx           sync.Mutex
}

func NewSyncService() *syncService {
	return &syncService{
		syncMgr:      newsyncMgr(),
		recordHandle: newRecordHandle(),
	}
}

func (s *syncService) Start() {
	if err := s.SyncBlocks(); err != nil {
		panic(err)
	}
}

func (s *syncService) Stop() {
}

func (s *syncService) SyncBlocks() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	lastBlock := s.syncMgr.CurrentBHeader()

	if err := s.syncBlock(lastBlock.Hash); err != nil {
		return err
	}
	return nil
}

func (s *syncService) syncBlock(lastBlockHash string) error {
	header := s.syncMgr.GetBlockHeaderByHash(lastBlockHash)
	txs := s.syncMgr.GetTxsByBlockHash(lastBlockHash)
	if header.Height == 0 {
		return nil
	}
	// sync blockchain header
	if err := s.syncBlockHeader(header, len(txs)); err != nil {
		return err
	}
	// sync blockchain txs
	if err := s.syncTxs(txs); err != nil {
		return err
	}

	return s.syncBlock(header.HashPrevBlock)
}

func (s *syncService) syncBlockHeader(header *BlockHeader, txCount int) error {
	carrier := &model.ChainBlockHeader{
		Height:           header.Height,
		Hash:             header.Hash,
		Version:          header.Version,
		HashPrevBlock:    header.HashPrevBlock,
		Timestamp:        header.Timestamp,
		Coinbase:         header.Coinbase,
		StateRoot:        header.StateRoot,
		TransactionsRoot: header.TransactionsRoot,
		ReceiptsRoot:     header.ReceiptsRoot,
		GasLimit:         header.GasLimit,
		GasUsed:          header.GasUsed,
		Bits:             header.Bits,
		Nonce:            header.Nonce,
		ExtraNonce:       header.ExtraNonce,
		TxCount:          txCount,
		Rewards:          float64(14),
	}
	if err := s.recordHandle.writeChainHeader(carrier); err != nil {
		return err
	}
	return nil
}

func (s *syncService) syncTxs(header *BlockHeader, txs []*Transaction) error {

	for _, tx := range txs {
		receipts := s.syncMgr.GetReceiptByHash(tx.Hash)
		carrier := &model.ChainBlockTx{
			BlockHash:   header.Hash,
			BlockHeight: header.Height,
			BlockTime:   header.Timestamp,
			Version:     int(header.Version),
			From:        tx.From,
			To:          tx.To,
			GasPrice:    tx.GasPrice,
			GasLimit:    tx.GasLimit,
			GasUsed:     receipts.GasUsed,
			// GasFee:      "",
			Data:  tx.Data,
			Nonce: tx.Nonce,
			Value: tx.Value,
			// Signature: "",
			Hash:   tx.Hash,
			Status: int(receipts.Status),
			Type:   1,
		}

		if err := s.syncAccouts(tx); err != nil {
			return err
		}
		if err := s.recordHandle.writeTxs(carrier); err != nil {
			return err
		}
	}
	return nil
}

func (s *syncService) syncAccouts(tx *Transaction) error {
	_ = s.syncMgr.GetAccountInfo(tx.From)
	_ = s.syncMgr.GetAccountInfo(tx.To)
	// carrier := &model.ChainAddress{
	// 	Address: tx.From,
	// }
	return nil
}
