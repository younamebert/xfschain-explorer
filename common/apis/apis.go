package apis

import "xfschainbrowser/model"

type LinkApi struct {
	HandleBlockHeader  *model.HandleChainBlockHeader
	HandleBlockTxs     *model.HandleChainBlockTx
	HandleChainAddress *model.HandleChainAddress
}

func NewLinkApi() *LinkApi {
	return &LinkApi{
		HandleBlockHeader:  new(model.HandleChainBlockHeader),
		HandleBlockTxs:     new(model.HandleChainBlockTx),
		HandleChainAddress: new(model.HandleChainAddress),
	}
}

var ApiResource = NewLinkApi()
