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
	if v := wr.handleChainAddress.Query("address = ?", addr); v != nil {
		return v[0]
	} else {
		return nil
	}
}
func (wr *recordHandle) QueryByHash(hash string) *model.ChainBlockHeader {
	if v := wr.handleBlockHeader.Query("hash = ?", hash); v != nil {
		return v[0]
	} else {
		return nil
	}
}

func (wr *recordHandle) QueryUp() *model.ChainBlockHeader {
	ra := wr.handleBlockHeader.QueryUp(1)
	if len(ra) > 1 {
		return ra[0]
	} else {
		return nil
	}
}

func (wr *recordHandle) QueryDown() *model.ChainBlockHeader {
	ra := wr.handleBlockHeader.QueryDown(1)
	if len(ra) > 1 {
		return ra[0]
	} else {
		return nil
	}
}
