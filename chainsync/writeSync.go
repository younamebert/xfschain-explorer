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
	if v := wr.handleChainAddress.Query("address = ?", addr); len(v) > 0 {
		return v[0]
	} else {
		return nil
	}
}
func (wr *recordHandle) QueryByHash(hash string) (result *model.ChainBlockHeader) {

	v := wr.handleBlockHeader.Query("hash = ?", hash)
	if v != nil {
		result = v[0]
		return
	}
	return
}

func (wr *recordHandle) QueryUp() (result *model.ChainBlockHeader) {
	ra := wr.handleBlockHeader.QuerySort(1, "height desc")
	if len(ra) > 0 {
		result = ra[0]
		return
	}
	return
}

func (wr *recordHandle) QueryDown() (result *model.ChainBlockHeader) {
	ra := wr.handleBlockHeader.QuerySort(1, "height asc")

	if len(ra) > 0 {
		result = ra[0]
		return
	}
	return
}
