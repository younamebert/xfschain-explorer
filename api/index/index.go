package api

import (
	"fmt"
	"math"
	"math/big"
	"net/http"
	"time"
	"xfschainbrowser/common"
	"xfschainbrowser/common/apis"

	"xfschainbrowser/model"

	"github.com/gin-gonic/gin"
)

type IndexLinkApi struct {
	Handle *apis.LinkApi
}

func (i *IndexLinkApi) Status(c *gin.Context) {

	//最高区块
	var blockHeader *model.ChainBlockHeader
	headers := i.Handle.HandleBlockHeader.QueryUp(1)

	if headers != nil {
		blockHeader = headers[0]
	}

	startTime := time.Now().AddDate(0, 0, -1).Unix()
	//全部交易
	txs := i.Handle.HandleBlockHeader.QueryTxCountSumByTime(1)

	//24小时区块和交易
	fmt.Println(startTime)
	afterBlock := i.Handle.HandleBlockHeader.QueryBlockHeadersByTime(startTime)
	var txAmount int64
	for _, v := range afterBlock {
		txAmount += int64(v.TxCount)
	}

	TxsInBlock := common.Div(txAmount, int64(len(afterBlock)))
	status := &StatusResp{
		LatestHeight: blockHeader.Height,
		Accounts:     i.Handle.HandleChainAddress.Count(),
		BlockRewards: "14.00",
		BlockTime:    blockHeader.Timestamp,
		Transactions: txs,
		Tps:          common.Div(txs, int64(3600)).String(),
		TxsInBlock:   TxsInBlock.BigInt().Int64(),
		Difficulty:   int64(common.BigByZip(new(big.Int).SetInt64(blockHeader.Bits))),
	}
	common.SendResponse(c, http.StatusOK, nil, status)
}

func (i *IndexLinkApi) LatestBlocksAndTxs(c *gin.Context) {
	lastTxsLimit := 10
	lastBlockLimit := 10
	blocks := i.Handle.HandleBlockHeader.QueryUp(int64(lastBlockLimit))
	txs := i.Handle.HandleBlockTxs.QueryLastBlockTxs(int64(lastTxsLimit))

	latest := &LatestResp{
		Blocks: blocks,
		Txs:    txs,
	}
	common.SendResponse(c, http.StatusOK, nil, latest)
}

func (i *IndexLinkApi) Search(c *gin.Context) {
	// tx_count_by_day
	// txcountByday :=

}

func (i *IndexLinkApi) TxCountByDay(c *gin.Context) {
	day := -7
	startTime := time.Now().AddDate(0, 0, day).Unix()
	endDay := int(math.Abs(float64(day)))
	blocks := i.Handle.HandleBlockHeader.QueryBlockHeadersByTime(startTime)

	times := time.Now().AddDate(0, 0, -1)
	timeDay := time.Date(times.Year(), times.Month(), times.Day(), 0, 0, 0, 0, times.Location()).Unix()

	TxCountByDays := make([]*TxCountByDayResp, int(day))
	for i := 1; i <= endDay; i++ {
		for _, v := range blocks {
			upTime := timeDay + (3600 * 1)
			stepTime := upTime - 3600
			if (v.Timestamp < upTime) && (v.Timestamp > stepTime) {
				TxCountByDayResp := &TxCountByDayResp{
					Timestamp: stepTime,
					TxCount:   int64(v.TxCount),
				}
				TxCountByDays = append(TxCountByDays, TxCountByDayResp)
			}
		}
	}
	common.SendResponse(c, http.StatusOK, nil, TxCountByDays)
}
