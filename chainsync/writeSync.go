package chainsync

import (
	"xfschainbrowser/model"
)

type recordHandle struct {
	handleBlockHeader  *model.HandleChainBlockHeader
	handleBlockTxs     *model.HandleChainBlockTx
	handleChainAddress *model.HandleChainAddress
}

func newRecordHandle() *recordHandle {
	return &recordHandle{
		handleBlockHeader:  new(model.HandleChainBlockHeader),
		handleBlockTxs:     new(model.HandleChainBlockTx),
		handleChainAddress: new(model.HandleChainAddress),
	}
}

// blockchain operation
func (wr *recordHandle) writeChainHeader(data *model.ChainBlockHeader) error {
	return wr.handleBlockHeader.Insert(data)
}

// transfer operation
func (wr *recordHandle) writeTxs(data *model.ChainBlockTx) error {
	return wr.handleBlockTxs.Insert(data)
}

// Account operation
func (wr *recordHandle) writeAccount(data *model.ChainAddress) error {
	return wr.handleChainAddress.Insert(data)
}

func (wr *recordHandle) updateAccount(data *model.ChainAddress) error {
	return wr.handleChainAddress.Update(data)
}

func (wr *recordHandle) QueryAccount(addr string) *model.ChainAddress {
	return wr.handleChainAddress.Query(addr)
}
func (wr *recordHandle) QueryByHash(hash string) *model.ChainBlockHeader {
	return wr.handleBlockHeader.QueryByHash(hash)
}

func (wr *recordHandle) QueryUp() *model.ChainBlockHeader {
	return wr.handleBlockHeader.QueryUp()
}

func (wr *recordHandle) QueryDown() *model.ChainBlockHeader {
	return wr.handleBlockHeader.QueryDown()
}
