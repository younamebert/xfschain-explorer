package chainsync

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"xfschainbrowser/common"
	"xfschainbrowser/common/progressbar"
	"xfschainbrowser/conf"
	"xfschainbrowser/global"
	"xfschainbrowser/model"

	"github.com/shopspring/decimal"
)

type syncService struct {
	chainMgr     chainMgr
	recordHandle *recordHandle
	mx           sync.Mutex
	wg           sync.WaitGroup
	bar          progressbar.Bar
}

func NewSyncService() *syncService {
	newsyncService := &syncService{
		chainMgr:     newsyncMgr(),
		recordHandle: newRecordHandle(),
	}
	newsyncService.UpdateBar()
	return newsyncService
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
			err := s.checkResponse()
			if err != nil {
				global.GVA_LOG.Error("chain service not found")
				continue
			} else {
				go s.SyncBlocks()
				s.barShow()
			}

		}
	}
}

func (s *syncService) Stop() {

}

func (s *syncService) checkResponse() error {
	return s.chainMgr.CheckResponse()
}

func (s *syncService) UpdateBar() {
	lastHeight := s.recordHandle.handleBlockHeader.QuerySort(1, "height desc")
	if len(lastHeight) > 0 {
		s.bar.NewOptionWithGraph(0, lastHeight[0].Height, "#")
	}
}

func (s *syncService) barShow() {
	var (
		count         int64 = 0
		lastHeight    int64 = 0
		currentHeight int64 = 0
	)
	count = s.recordHandle.handleBlockHeader.Count(nil, nil)

	currentBlock := s.recordHandle.handleBlockHeader.QuerySort(1, "height asc")

	if len(currentBlock) > 0 {

		currentHeight = currentBlock[0].Height
		// 向下同步完成不在需要打印同步进度
		if currentHeight == 1 {
			return
		}
	}

	lastBlock := s.recordHandle.handleBlockHeader.QuerySort(1, "height desc")
	if len(lastBlock) > 0 {
		lastHeight = lastBlock[0].Height
	}

	s.bar.Play(count, lastHeight, currentHeight)
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
	var (
		lastBlockChain *BlockHeader              = new(BlockHeader)
		lastBlocks     []*model.ChainBlockHeader = make([]*model.ChainBlockHeader, 0)
		nextBlock      []*model.ChainBlockHeader = make([]*model.ChainBlockHeader, 0)
	)
	lastBlockChain = s.chainMgr.CurrentBHeader()

	//低到高
	lastBlocks = s.recordHandle.handleBlockHeader.QuerySort(conf.CheckIntervalBlock, "height desc")

	//从高往小同步
	nextBlock = s.recordHandle.handleBlockHeader.QuerySort(1, "height asc")

	if lastBlockChain == nil {
		return nextBlock[0].HashPrevBlock
	}

	//数据库没有当前最高高度
	if len(lastBlocks) < 1 {
		return lastBlockChain.Hash
	}

	//链上的最高块和是否存在数据库(同步)
	if lastBlocks[0].Height < lastBlockChain.Height {
		// 链最高和本地最高
		disparity := lastBlockChain.Height - lastBlocks[0].Height
		go s.handleMissBlock(disparity, lastBlocks[0])
		// 需要更新bar
		s.UpdateBar()
	}

	// 验证最高块的交易数据是否全部同步完成.
	go s.handleMissTx(lastBlocks)

	return nextBlock[0].HashPrevBlock
}

// 验证数据库的最高块是否有连续
func (s *syncService) handleMissBlock(disparity int64, block *model.ChainBlockHeader) {
	var handleMissBlockFetchNumber int = 0

	//区块差小于设置同步跨度差，直接disparity number
	if disparity > int64(conf.HandleMissBlockFetch) {
		handleMissBlockFetchNumber = conf.HandleMissBlockFetch
	} else {
		handleMissBlockFetchNumber = int(disparity)
	}

	s.wg.Add(handleMissBlockFetchNumber)
	for i := 1; i <= handleMissBlockFetchNumber; i++ {
		nextSyncBlockNumber := strconv.FormatInt(block.Height+int64(i), 10)
		addBlock := s.chainMgr.GetBlockHeaderByNumber(nextSyncBlockNumber)
		if addBlock != nil && (addBlock.HashPrevBlock == block.Hash) {
			go s.syncBlock(addBlock.Hash)
			global.GVA_LOG.Info(fmt.Sprintf("sync mode:asc order fetch block targetHeight:%v currentHeight:%v", nextSyncBlockNumber, block.Height))
		} else {
			continue
		}
		s.wg.Done()
	}
	s.wg.Wait()
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
	// s.barShow()
	return nil
}

func (s *syncService) syncBlock(lastBlockHash string) error {
	header := s.chainMgr.GetBlockHeaderByHash(lastBlockHash)
	txs := s.chainMgr.GetTxsByBlockHash(lastBlockHash)

	if header == nil {
		return nil
	}

	if header.Height == 0 {
		return nil
	}
	// sync blockchain header
	if err := s.syncBlockHeader(header, len(txs)); err != nil {
		global.GVA_LOG.Warn(err.Error())
		return err
	}
	// sync blockchain txs
	go s.syncTxs(header, txs)
	// if err := go s.syncTxs(header, txs); err != nil {
	// 	return err
	// }
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
		global.GVA_LOG.Warn(err.Error())
		return err
	}
	return nil
}

func (s *syncService) syncTxs(header *BlockHeader, txs []*Transaction) {

	for _, tx := range txs {
		receipts := s.chainMgr.GetReceiptByHash(tx.Hash)
		if receipts == nil {
			continue
		}
		gasuesd := decimal.NewFromInt(receipts.GasUsed)
		gasPice := decimal.NewFromFloat(tx.GasPrice)
		gasLimit := decimal.NewFromFloat(tx.GasLimit)
		gasuesdFloat, ok := gasuesd.Float64()
		if !ok {
			global.GVA_LOG.Warn(fmt.Sprintf("txhash:%v blockhash:%v blockheight:%v", tx.Hash, header.Hash, header.Height))
		}
		carrier := &model.ChainBlockTx{
			BlockHash:   header.Hash,
			BlockHeight: header.Height,
			BlockTime:   header.Timestamp,
			Version:     int(header.Version),
			TxFrom:      tx.From,
			TxTo:        tx.To,
			GasPrice:    gasPice,
			GasLimit:    gasLimit,
			GasUsed:     gasuesd,
			GasFee:      common.CalcGasFee(gasuesdFloat, tx.GasPrice).String(),
			Data:        tx.Data,
			Nonce:       tx.Nonce,
			Value:       tx.Value.String(),
			Signature:   tx.Signature,
			Hash:        tx.Hash,
			Status:      int(receipts.Status),
			Type:        0,
		}
		if carrier.TxTo == "" {
			carrier.Type = 1
		}
		if err := s.syncAccouts(header, tx); err != nil {
			global.GVA_LOG.Warn(err.Error())
			return
		}
		if err := s.recordHandle.writeTxs(carrier); err != nil {
			global.GVA_LOG.Warn(err.Error())
			return
		}
	}
}

func (s *syncService) syncAccouts(header *BlockHeader, tx *Transaction) error {
	addrs := []string{tx.From, tx.To}
	for _, v := range addrs {
		go s.updateAccount(header, tx, v)
	}
	return nil
}

func (s *syncService) updateAccount(header *BlockHeader, tx *Transaction, addr string) error {
	var (
		obj      *model.ChainAddress = new(model.ChainAddress)
		objChain *AccountState       = new(AccountState)
		carrier  *model.ChainAddress
	)

	obj = s.recordHandle.QueryAccount(addr)
	objChain = s.chainMgr.GetAccountInfo(addr)
	if objChain == nil {
		global.GVA_LOG.Warn(fmt.Sprintf("GetAccountInfo:%v error:%v", addr, objChain))
		return nil
	}

	carrier = new(model.ChainAddress)
	if obj == nil {
		carrier.Address = addr
		carrier.Balance = objChain.Balance
		carrier.Nonce = objChain.Nonce
		carrier.Extra = objChain.Extra
		carrier.Code = objChain.Code
		carrier.StateRoot = objChain.StateRoot
		carrier.Type = 1
		if carrier.Code == "" {
			carrier.Type = 0
		}
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
		err := s.recordHandle.writeAccount(carrier)
		if err != nil {
			global.GVA_LOG.Warn(err.Error())
		}
		return err
	}

	carrier.Address = tx.From
	carrier.Balance = objChain.Balance
	carrier.Nonce = objChain.Nonce
	carrier.Extra = objChain.Extra
	carrier.Code = objChain.Code
	carrier.StateRoot = objChain.StateRoot
	carrier.FromStateRoot = header.StateRoot
	carrier.FromBlockHeight = header.Height
	carrier.FromBlockHash = header.Hash
	carrier.TxCount = obj.TxCount + 1
	if header.Height > obj.CreateFromBlockHeight {
		carrier.CreateFromBlockHash = header.Hash
		carrier.CreateFromBlockHeight = header.Height
		carrier.CreateFromBlockHash = header.Hash
		carrier.CreateFromStateRoot = header.StateRoot
		carrier.CreateFromTxHash = tx.Hash
	}
	err := s.recordHandle.updateAccount(carrier)
	if err != nil {
		global.GVA_LOG.Warn(err.Error())
	}
	return err
}
