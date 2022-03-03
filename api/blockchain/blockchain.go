package api

import (
	"net/http"
	"strconv"
	"xfschainbrowser/common"
	"xfschainbrowser/common/apis"
	"xfschainbrowser/conf"

	"github.com/gin-gonic/gin"
)

type BlocksLinkApi struct {
	Handle *apis.LinkApi
}

func (bc *BlocksLinkApi) GetBlocks(c *gin.Context) {

	var (
		page     int
		pageSize int
		result   *apis.Pages = new(apis.Pages)
	)
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = conf.Page
	}

	pageSize, err = strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = conf.PageSize
	}

	counts := bc.Handle.HandleBlockHeader.Count(nil, nil)
	blocks := bc.Handle.HandleBlockHeader.GetBlocks(nil, nil, page, pageSize)
	if len(blocks) == 0 || blocks == nil {
		common.SendResponse(c, http.StatusOK, nil, nil)
		return
	}

	result = &apis.Pages{
		Page:     int64(page),
		PageSize: int64(pageSize),
		Limits:   counts,
		Data:     blocks,
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}

func (bc *BlocksLinkApi) Detailed(c *gin.Context) {

	var (
		hash   string
		result *apis.Pages = new(apis.Pages)
	)
	hash = c.Query("hash")
	if hash == "" {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	blockObj := bc.Handle.HandleBlockHeader.Query("hash = ?", hash)
	if len(blockObj) == 0 || blockObj == nil {
		common.SendResponse(c, http.StatusOK, nil, nil)
		return
	}

	result = &apis.Pages{
		Data: blockObj[0],
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}

func (bc *BlocksLinkApi) DetailedTx(c *gin.Context) {
	var (
		blockhash string
		page      int
		pageSize  int
		result    *apis.Pages = new(apis.Pages)
	)

	blockhash = c.Query("blockhash")
	if blockhash == "" {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = conf.Page
	}

	pageSize, err = strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = conf.PageSize
	}

	txs := bc.Handle.HandleBlockTxs.GetTxs(`block_hash = ?`, blockhash, page, pageSize)
	if len(txs) == 0 || txs == nil {
		common.SendResponse(c, http.StatusOK, nil, nil)
		return
	}
	limits := bc.Handle.HandleBlockTxs.Count(`block_hash = ?`, blockhash)
	result = &apis.Pages{
		Page:     int64(page),
		PageSize: int64(pageSize),
		Limits:   limits,
		Data:     txs,
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}
