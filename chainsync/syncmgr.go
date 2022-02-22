package chainsync

import (
	"math/big"
	"strconv"
	"sync"
	"time"
	"xfschainbrowser/common"
	"xfschainbrowser/conf"
	"xfschainbrowser/model"
)

type syncService struct {
	chainMgr     chainMgr
	recordHandle *recordHandle
	mx           sync.Mutex
}

func NewSyncService() *syncService {
	return &syncService{
		chainMgr:     newsyncMgr(),
		recordHandle: newRecordHandle(),
	}
}

func (s *syncService) Start() {
	if err := s.process(); err != nil {
		panic(err)
	}
}

func (s *syncService) process() error {
	timeDur, err := time.ParseDuration(conf.SyncForceSyncPeriod)
	if err != nil {
		return err
	}
	for {
		select {
		case <-time.After(timeDur):
			if err := s.SyncBlocks(); err != nil {
				return err
			}
		}
	}
}

func (s *syncService) Stop() {
}

func (s *syncService) SyncBlocks() error {
	s.mx.Lock()
	defer s.mx.Unlock()

	lastBlock := s.chainMgr.CurrentBHeader()
	if v := s.recordHandle.QueryByHash(lastBlock.Hash); v != nil {
		downBlock := s.recordHandle.QueryDown()
		if err := s.syncBlock(downBlock.HashPrevBlock); err != nil {
			return err
		}
	} else {
		if err := s.syncBlock(lastBlock.Hash); err != nil {
			return err
		}
	}
	return nil
}

func (s *syncService) syncBlock(lastBlockHash string) error {

	header := s.chainMgr.GetBlockHeaderByHash(lastBlockHash)
	txs := s.chainMgr.GetTxsByBlockHash(lastBlockHash)
	if header.Height == 0 {
		return nil
	}
	// sync blockchain header
	if err := s.syncBlockHeader(header, len(txs)); err != nil {
		return err
	}
	// sync blockchain txs
	if err := s.syncTxs(header, txs); err != nil {
		return err
	}
	return nil
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
		receipts := s.chainMgr.GetReceiptByHash(tx.Hash)
		gasuesd, _ := new(big.Float).SetInt64(receipts.GasUsed).Float64()
		carrier := &model.ChainBlockTx{
			BlockHash:   header.Hash,
			BlockHeight: header.Height,
			BlockTime:   header.Timestamp,
			Version:     int(header.Version),
			From:        tx.From,
			To:          tx.To,
			GasPrice:    tx.GasPrice,
			GasLimit:    tx.GasLimit,
			GasUsed:     gasuesd,
			GasFee:      common.CalcGasFee(gasuesd, tx.GasPrice).String(),
			Data:        tx.Data,
			Nonce:       tx.Nonce,
			Value:       strconv.FormatInt(tx.Value, 10),
			// Signature: "",
			Hash:   tx.Hash,
			Status: int(receipts.Status),
			Type:   1,
		}

		if err := s.syncAccouts(header, tx); err != nil {
			return err
		}
		if err := s.recordHandle.writeTxs(carrier); err != nil {
			return err
		}
	}
	return nil
}

func (s *syncService) syncAccouts(header *BlockHeader, tx *Transaction) error {

	if err := s.updateAccount(header, tx, tx.From); err != nil {
		return err
	}
	if err := s.updateAccount(header, tx, tx.To); err != nil {
		return err
	}
	return nil
}

func (s *syncService) updateAccount(header *BlockHeader, tx *Transaction, addr string) error {
	obj := s.recordHandle.QueryAccount(addr)
	objChain := s.chainMgr.GetAccountInfo(addr)

	carrier := new(model.ChainAddress)
	if obj == nil {
		carrier.Address = addr
		carrier.Balance = objChain.Balance
		carrier.Nonce = objChain.Nonce
		carrier.Extra = objChain.Extra
		carrier.Code = objChain.Code
		carrier.StateRoot = objChain.StateRoot
		carrier.Type = 1
		carrier.Display = 1
		carrier.FromStateRoot = header.StateRoot
		carrier.FromBlockHeight = header.Height
		carrier.FromBlockHash = header.Hash
		carrier.CreateFromBlockHash = header.Hash
		carrier.CreateFromBlockHeight = header.Height
		carrier.CreateFromBlockHash = header.Hash
		carrier.CreateFromStateRoot = header.StateRoot
		carrier.CreateFromTxHash = tx.Hash

		return s.recordHandle.writeAccount(carrier)
	} else {
		if header.Height > obj.CreateFromBlockHeight {
			carrier.CreateFromBlockHash = header.Hash
			carrier.CreateFromBlockHeight = header.Height
			carrier.CreateFromBlockHash = header.Hash
			carrier.CreateFromStateRoot = header.StateRoot
			carrier.CreateFromTxHash = tx.Hash
		}
		if header.Height < obj.CreateFromBlockHeight {
			carrier.Address = tx.From
			carrier.Balance = objChain.Balance
			carrier.Nonce = objChain.Nonce
			carrier.Extra = objChain.Extra
			carrier.Code = objChain.Code
			carrier.StateRoot = objChain.StateRoot
			carrier.FromStateRoot = header.StateRoot
			carrier.FromBlockHeight = header.Height
			carrier.FromBlockHash = header.Hash

		}
		return s.recordHandle.updateAccount(carrier)
	}
}
