package api

import (
	"net/http"
	"xfschainbrowser/common"
	"xfschainbrowser/model"

	"github.com/gin-gonic/gin"
)

type IndexLinkApi struct {
	HandleBlockHeader  *model.HandleChainBlockHeader
	HandleBlockTxs     *model.HandleChainBlockTx
	HandleChainAddress *model.HandleChainAddress
}

func (i *IndexLinkApi) Status(c *gin.Context) {
	blockHeader := i.HandleBlockHeader.QueryUp()

	// var startTime = time.Now().AddDate(0, 0, -1).Unix()

	_ = &StatusResp{
		LatestHeight: blockHeader.Height,
		Accounts:     i.HandleChainAddress.Count(),
		BlockRewards: "14.00",
		BlockTime:    blockHeader.Timestamp,
		Difficulty:   common.BitsUnzip(uint32(blockHeader.Bits)).Int64(),
	}
	common.SendResponse(c, http.StatusOK, nil, "nihao")
}

// func (i *IndexLinkApi)
