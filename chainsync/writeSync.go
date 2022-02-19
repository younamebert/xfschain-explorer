package chainsync

import (
	"xfschainbrowser/model"
)

type recordHandle struct {
	handleBlockHeader  *model.HandlerChainBlockHeader
	handleBlockTxs     *model.HandlerChainBlockTx
	handleChainAddress *model.HandleChainAddress
}

func newRecordHandle() *recordHandle {
	return &recordHandle{
		handleBlockHeader: new(model.HandlerChainBlockHeader),
		handleBlockTxs:    new(model.HandlerChainBlockTx),
	}
}

func (wr *recordHandle) writeChainHeader(data *model.ChainBlockHeader) error {
	return wr.handleBlockHeader.Insert(data)
}

func (wr *recordHandle) writeTxs(data *model.ChainBlockTx) error {
	return wr.handleBlockTxs.Insert(data)
}

func (wr *recordHandle) writeAccount(data *model.ChainAddress) error {
	return wr.handleChainAddress.Insert(data)
}
