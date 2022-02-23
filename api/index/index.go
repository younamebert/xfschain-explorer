package api

import (
	"fmt"
	"math/big"
	"net/http"
	"time"
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

	//最高区块
	blockHeader := i.HandleBlockHeader.QueryUp()

	startTime := time.Now().AddDate(0, 0, -1).Unix()
	//全部交易
	txs := i.HandleBlockHeader.QueryTxCountSumByTime(1)

	//24小时区块和交易
	fmt.Println(startTime)
	afterBlock := i.HandleBlockHeader.QueryBlockHeadersByTime(startTime, 0)
	afterTxs := i.HandleBlockHeader.QueryTxCountSumByTime(startTime)

	fmt.Println(len(afterBlock), afterTxs)
	TxsInBlock := int(afterTxs) / len(afterBlock)
	// new(
	status := &StatusResp{
		LatestHeight: blockHeader.Height,
		Accounts:     i.HandleChainAddress.Count(),
		BlockRewards: "14.00",
		BlockTime:    blockHeader.Timestamp,
		Transactions: txs,
		Tps:          common.Div(txs, int64(3600)).String(),
		TxsInBlock:   int64(TxsInBlock),
		Difficulty:   int64(common.BigByZip(new(big.Int).SetInt64(blockHeader.Bits))),
	}
	common.SendResponse(c, http.StatusOK, nil, status)
}

func (i *IndexLinkApi) LatestBlocksAndTxs(c *gin.Context) {

	// LatestResp{}
}
