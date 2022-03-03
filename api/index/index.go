package api

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
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

	var (
		txAmount      int64
		totalWorkload *big.Int = big.NewInt(0)
		// tps           decimal.Decimal
		currentTime   int64 = time.Now().Unix()
		startTime     int64
		BlockTimeMulX int64
		txs           int64
		blockHeader   *model.ChainBlockHeader
	)

	//最高区块
	headers := i.Handle.HandleBlockHeader.QuerySort(1, "height desc")

	if headers != nil {
		blockHeader = headers[0]
	}

	startTime = time.Now().AddDate(0, 0, -5).Unix()
	//全部交易
	txs = i.Handle.HandleBlockHeader.QueryTxCountSumByTime(1)

	//条件区间小时区块和交易
	afterBlock := i.Handle.HandleBlockHeader.Query("timestamp > ?", startTime)

	for _, v := range afterBlock {
		txAmount += int64(v.TxCount)
		blockWorkload := common.BigByZip(new(big.Int).SetInt64(v.Bits))
		totalWorkload.Add(totalWorkload, new(big.Int).SetUint64(uint64(blockWorkload)))
	}
	totalWorkload.Div(totalWorkload, big.NewInt(24*60*60))
	tpsStatus, _ := common.Div(txAmount, 24*60*60).Float64()

	BlockTimeMulX = currentTime - startTime
	BlockTimeTotal := common.Div(BlockTimeMulX, int64(len(afterBlock)))
	BlockTimeTotalSecond, _ := BlockTimeTotal.Float64()

	rewards, _ := common.BaseCoin2Atto("14")
	TxsInBlock := common.Div(txAmount, int64(len(afterBlock)))

	status := &StatusResp{
		LatestHeight: blockHeader.Height,
		Accounts:     i.Handle.HandleChainAddress.Count(nil, nil),
		BlockRewards: rewards.String(),
		BlockTime:    BlockTimeTotalSecond,
		Transactions: txs,
		Power:        totalWorkload.Int64(),
		Tps:          tpsStatus,
		TxsInBlock:   TxsInBlock.BigInt().Int64(),
		Difficulty:   int64(common.BigByZip(new(big.Int).SetInt64(blockHeader.Bits))),
	}
	common.SendResponse(c, http.StatusOK, nil, status)
}

func (i *IndexLinkApi) LatestBlocksAndTxs(c *gin.Context) {
	lastTxsLimit := 10
	lastBlockLimit := 10
	blocks := i.Handle.HandleBlockHeader.QuerySort(int64(lastBlockLimit), "height desc")
	txs := i.Handle.HandleBlockTxs.QueryLastBlockTxs(int64(lastTxsLimit))

	latest := &LatestResp{
		Blocks: blocks,
		Txs:    txs,
	}

	common.SendResponse(c, http.StatusOK, nil, latest)
}

func (i *IndexLinkApi) Search(c *gin.Context) {
	var (
		param  string
		result *SearchResp = new(SearchResp)
	)
	param = c.Query("q")
	if param == "" && len(param) > 100 {
		common.SendResponse(c, http.StatusBadRequest, errors.New("illegal parameter Error"), nil)
	}
	params := "%" + param + "%"
	fmt.Println(len(param))

	blocks := i.Handle.HandleBlockHeader.QueryLike("hash like ?", params)
	if (blocks != nil) && (len(blocks) > 0) {
		result.Type = 1
		result.PathValue = blocks[0].Hash
		common.SendResponse(c, http.StatusOK, nil, result)
		return
	}

	txs := i.Handle.HandleBlockTxs.QueryLikeTx("hash like ?", params)
	if (txs != nil) && (len(txs) > 0) {
		result.Type = 2
		result.PathValue = txs[0].Hash
		common.SendResponse(c, http.StatusOK, nil, result)
		return
	}

	accounts := i.Handle.HandleChainAddress.QueryLikeAccount("address = ?", param)
	if (accounts != nil) && (len(accounts) > 0) {
		result.Type = 3
		result.PathValue = accounts[0].Address
		common.SendResponse(c, http.StatusOK, nil, result)
		return
	}

	blocks = i.Handle.HandleBlockHeader.QueryLike("height = ?", param)
	if (blocks != nil) && (len(blocks) > 0) {
		result.Type = 4
		result.PathValue = blocks[0].Hash
		common.SendResponse(c, http.StatusOK, nil, result)
		return
	}

	common.SendResponse(c, http.StatusOK, nil, nil)
}

func (i *IndexLinkApi) TxCountByDay(c *gin.Context) {

	param, err := strconv.Atoi(c.Query("p"))
	if err != nil {
		param = 7
	}

	param = int(math.Abs(float64(param)))
	paramLenTime := ^param + 2
	beforeTime := common.GetBeforeTime(paramLenTime)
	result := make([]*TxCountByDayResp, param)

	blocks := i.Handle.HandleBlockHeader.Query("timestamp > ?", beforeTime)
	if blocks == nil {
		common.SendResponse(c, http.StatusOK, nil, result)
	}

	for i := 1; i <= param; i++ {
		upTime := int(beforeTime) + (86400 * i)
		stepTime := upTime - 86400
		txCountDay := new(TxCountByDayResp)
		txCountDay.Timestamp = int64(stepTime)
		for _, v := range blocks {
			if (v.Timestamp < int64(upTime)) && (v.Timestamp > int64(stepTime)) {
				txCountDay.TxCount = txCountDay.TxCount + int64(v.TxCount)
			}
		}
		result[i-1] = txCountDay
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}
