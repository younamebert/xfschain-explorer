package api

import (
	"errors"
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
		txAmount      int64    = 0
		totalWorkload *big.Int = big.NewInt(0)
		// tps           decimal.Decimal
		currentTime   int64 = time.Now().Unix()
		startTime     int64
		BlockTimeMulX int64
		txsCount      int64
		result        *StatusResp = new(StatusResp)
		blockHeader   *model.ChainBlockHeader
	)

	//最高区块
	headers := i.Handle.HandleBlockHeader.QuerySort(1, "height desc")
	if len(headers) < 1 {
		common.SendResponse(c, http.StatusOK, nil, result)
		return
	}

	blockHeader = headers[0]
	startTime = time.Now().Unix() - 86400

	//条件区间小时区块和交易
	afterBlock := i.Handle.HandleBlockHeader.Query("timestamp > ?", startTime)
	afterBlockLen := len(afterBlock)
	// if len(afterBlock) < 1 {
	// 	fmt.Println(12212121)
	// 	common.SendResponse(c, http.StatusOK, nil, result)
	// 	return
	// }

	for _, v := range afterBlock {
		txAmount += int64(v.TxCount)
		blockWorkload := common.BigByZip(new(big.Int).SetInt64(v.Bits))
		totalWorkload.Add(totalWorkload, new(big.Int).SetUint64(uint64(blockWorkload)))
	}

	totalWorkload.Div(totalWorkload, big.NewInt(24*60*60))
	tpsStatus, _ := common.Div(txAmount, 24*60*60).Float64()

	BlockTimeMulX = currentTime - startTime
	BlockTimeTotal := common.Div(BlockTimeMulX, int64(afterBlockLen))
	BlockTimeTotalSecond, _ := BlockTimeTotal.Float64()

	rewards, _ := common.BaseCoin2Atto("14")
	TxsInBlock := common.Div(txAmount, int64(afterBlockLen))

	//全部交易
	txsCount = i.Handle.HandleBlockHeader.QueryTxCountSumByTime(1)
	result = &StatusResp{
		LatestHeight: blockHeader.Height,
		Accounts:     i.Handle.HandleChainAddress.Count(nil, nil),
		BlockRewards: rewards.String(),
		BlockTime:    BlockTimeTotalSecond,
		Transactions: txsCount,
		Power:        totalWorkload.Int64(),
		Tps:          tpsStatus,
		TxsInBlock:   TxsInBlock.BigInt().Int64(),
		Difficulty:   int64(common.BigByZip(new(big.Int).SetInt64(blockHeader.Bits))),
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}

func (i *IndexLinkApi) LatestBlocksAndTxs(c *gin.Context) {

	var (
		lastTxsLimit   int = 10
		lastBlockLimit int = 10
	)

	txlimit, err := strconv.Atoi(c.Query("txlimit"))
	if err == nil {
		lastTxsLimit = txlimit
	}

	blockLimit, err := strconv.Atoi(c.Query("blocklimit"))
	if err == nil {
		lastBlockLimit = blockLimit
	}

	// txlimit = c.Query("txlimit")
	// if addr == "" {
	// 	common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
	// 	return
	// // }
	// lastTxsLimit := 10
	// lastBlockLimit := 10
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

	var (
		param      int
		today      time.Time = time.Now()
		oneday     int       = 86400
		result     []*TxCountByDayResp
		startTime  int64
		beforeTime int64
	)
	param, err := strconv.Atoi(c.Query("p"))
	if err != nil {
		param = 7
	}

	param = int(math.Abs(float64(param)))

	startTime = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).Unix()
	beforeTime = startTime - int64(oneday*param)
	blocks := i.Handle.HandleBlockHeader.Query("timestamp > ?", beforeTime)

	result = make([]*TxCountByDayResp, param)
	for i := 1; i <= param; i++ {
		nextTime := startTime - int64((oneday * i))
		currentTime := nextTime + int64(oneday)
		txCountDay := new(TxCountByDayResp)
		txCountDay.Timestamp = int64(currentTime)

		for _, v := range blocks {
			if (v.Timestamp < int64(nextTime)) && (v.Timestamp > int64(nextTime)) {
				txCountDay.TxCount = txCountDay.TxCount + int64(v.TxCount)
			}
		}
		result[i-1] = txCountDay
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}
