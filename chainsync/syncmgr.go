package chainsync

import (
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"
	"xfschainbrowser/common"
	"xfschainbrowser/common/progressbar"
	"xfschainbrowser/conf"
	"xfschainbrowser/global"
	"xfschainbrowser/model"
)

type syncService struct {
	chainMgr     chainMgr
	recordHandle *recordHandle
	mx           sync.Mutex
	// wg           sync.WaitGroup
	bar progressbar.Bar
}

func NewSyncService() *syncService {
	return &syncService{
		chainMgr:     newsyncMgr(),
		recordHandle: newRecordHandle(),
	}
}

func (s *syncService) Start() {
	s.UpdateBar()
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
			go s.SyncBlocks()
			s.barShow()
		}
	}
}

func (s *syncService) Stop() {
}

func (s *syncService) UpdateBar() {
	lastHeight := s.recordHandle.handleBlockHeader.QuerySort(1, "height desc")
	if len(lastHeight) > 0 {
		s.bar.NewOptionWithGraph(0, lastHeight[0].Height, "#")
	}
}

func (s *syncService) barShow() {
	count := s.recordHandle.handleBlockHeader.Count(nil, nil)
	s.bar.Play(count)
	s.bar.Finish()

}
func (s *syncService) SyncBlocks() {
	for i := 0; i < conf.SyncMaxBlocksFetch; i++ {
		go func() {
			if err := s.syncBlocks(); err != nil {
				panic(err)
			}
		}()
	}
}

func (s *syncService) checkIntervalBlockTxs() string {

	//链上的最高块
	lastBlockChain := s.chainMgr.CurrentBHeader()

	lastBlocks := s.recordHandle.handleBlockHeader.QuerySort(conf.CheckIntervalBlock, "height desc")
	//数据库没有当前最高高度
	if len(lastBlocks) < 1 {
		return lastBlockChain.Hash
	}

	//链上的最高块和是否存在数据库(同步)
	if lastBlocks[0].Height < lastBlockChain.Height {
		disparity := lastBlockChain.Height - lastBlocks[0].Height
		go s.handleMissBlock(disparity, lastBlocks[0])
	}

	// 验证最高块的交易数据是否全部同步完成.
	go s.handleMissTx(lastBlocks)

	//从高往小同步
	nextBlock := s.recordHandle.handleBlockHeader.QuerySort(1, "height asc")

	return nextBlock[0].HashPrevBlock
}

// 验证数据库的最高块是否有连续
func (s *syncService) handleMissBlock(disparity int64, block *model.ChainBlockHeader) {
	for i := 0; i < int(disparity); i++ {
		nextSyncBlockNumber := strconv.FormatInt(block.Height+int64(i), 10)
		pending := s.chainMgr.GetBlockHeaderByNumber(nextSyncBlockNumber)
		if pending != nil {
			s.syncBlock(pending.Hash)
		} else {
			global.GVA_LOG.Error(fmt.Sprintf("Block bifurcation height:%v", nextSyncBlockNumber))
		}
	}

}

// 验证最区块的交易数据是否全部同步完成
func (s *syncService) handleMissTx(lastBlocks []*model.ChainBlockHeader) {
	for _, bks := range lastBlocks {
		txsLen := len(s.chainMgr.GetTxsByBlockHash(bks.Hash))
		dbTxsLen := s.recordHandle.handleBlockTxs.Count("block_hash =? ", bks.Hash)
		if txsLen != int(dbTxsLen) {
			if err := s.syncBlock(bks.Hash); err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("blockHash:%v blockHeight:%v err func:handleMissTx error:%v", bks.Hash, bks.Height, err.Error()))
				continue
			}
		}
	}
}

func (s *syncService) syncBlocks() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	lastBlockHash := s.checkIntervalBlockTxs()
	if err := s.syncBlock(lastBlockHash); err != nil {
		return err
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
	rewards, _ := common.BaseCoin2Atto("14")
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
		Rewards:          rewards.String(),
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
			TxFrom:      tx.From,
			TxTo:        tx.To,
			GasPrice:    tx.GasPrice,
			GasLimit:    tx.GasLimit,
			GasUsed:     gasuesd,
			GasFee:      common.CalcGasFee(gasuesd, tx.GasPrice).String(),
			Data:        tx.Data,
			Nonce:       tx.Nonce,
			Value:       tx.Value.String(),
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
		carrier.TxCount = 0
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
			carrier.TxCount = obj.TxCount + 1
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
