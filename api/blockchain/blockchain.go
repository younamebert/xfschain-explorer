package api

import (
	"net/http"
	"strconv"
	"xfschainbrowser/common"
	"xfschainbrowser/common/apis"

	"github.com/gin-gonic/gin"
)

type BlockChainLinkApi struct {
	Handle *apis.LinkApi
}

func (bc *BlockChainLinkApi) GetBlocks(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		common.SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}
	pageSize := 20
	limit := 20
	if page > 0 && pageSize > 0 {
		txs := bc.Handle.HandleBlockHeader.GetBlocks(page, pageSize, limit)
		common.SendResponse(c, http.StatusOK, nil, txs)
		return
	}
}
