package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mi/common"
	"mi/common/apis"
	"mi/conf"
	"mi/events"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type EquipmentLinkApi struct {
	Handle *apis.LinkApi
}

func (ac *EquipmentLinkApi) EquipmentList(c *gin.Context) {
	var (
		page     int         = conf.Page
		pageSize int         = conf.PageSize
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
	count := ac.Handle.HandleMiEquipment.Count(nil, nil)
	list := ac.Handle.HandleMiEquipment.Querys(nil, nil, page, pageSize)
	if len(list) < 1 {
		common.SendResponse(c, http.StatusOK, nil, nil)
		return
	}

	result = &apis.Pages{
		Page:     int64(page),
		PageSize: int64(pageSize),
		Limits:   count,
		Data:     list,
	}
	common.SendResponse(c, http.StatusOK, nil, result)
}

func (ac *EquipmentLinkApi) EquipmentSwitch(c *gin.Context) {
	//   c.Qquer
	// ac.Handle.HandleMiEquipment.Insert()
}

//关闭/打开广告屏 0关 1开
func (ac *EquipmentLinkApi) SwitchAdvertising(c *gin.Context) {
	//设备id
	iccid := c.Query("iccid")
	status, err := strconv.Atoi(c.Query("status"))

	if iccid == "" || err != nil {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if (status != apis.OpenStatus) && (status != apis.CloseStatus) {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if err := ac.Handle.HandleMiEquipment.SetSwitchad(iccid, status); err != nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
	}

	// push events sub	   //创建事件
	ac.Handle.EventsBus.Publish(events.SendNoticeEvent{Iccid: iccid, Data: apis.Screen(status)})

	common.SendResponse(c, http.StatusOK, nil, nil)
}

//关闭/打开灯带0关 1开
func (ac *EquipmentLinkApi) SwitchLed(c *gin.Context) {
	//iccid
	iccid := c.Query("iccid")
	status, err := strconv.Atoi(c.Query("status"))

	if iccid == "" || err != nil {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if (status != apis.OpenStatus) && (status != apis.CloseStatus) {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if err := ac.Handle.HandleMiEquipment.SetSwitchadLed(iccid, status); err != nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
	}

	// push events sub	   //创建事件
	ac.Handle.EventsBus.Publish(events.SendNoticeEvent{Iccid: iccid, Data: apis.SwitchIed(status)})
	common.SendResponse(c, http.StatusOK, nil, nil)
}

type Data struct {
	Iccid  string  `json:"iccid"`
	APrice float64 `json:"a_price"`
	BPrice float64 `json:"b_price"`
}

//修改仓库价格
//handle 接收到修改单价命令的话，重新查询单价
func (ac *EquipmentLinkApi) UpdatePrice(c *gin.Context) {
	//截取参数
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
		return
	}
	//取
	parameter := new(SetPriceArgs)
	if err := json.Unmarshal(body, &parameter); err != nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
		return
	}
	if parameter.Iccid == "" {
		common.SendResponse(c, http.StatusBadRequest, common.NotIccidErr, nil)
		return
	}

	// code
	priceCode := make([]decimal.Decimal, 2)

	//iccid的所有长裤
	euqicondition := make(map[string]interface{}, 0)
	euqicondition["iccid"] = parameter.Iccid
	euqicondition["status"] = 1
	euqiA := ac.Handle.HandleMiWarehouse.BeartQuery(euqicondition)
	priceCode[0] = decimal.NewFromFloat(euqiA.WarehousePrice)
	if parameter.APrice != "" {
		APrice, err := decimal.NewFromString(parameter.APrice)
		if err != nil {
			common.SendResponse(c, http.StatusBadRequest, err, nil)
			return
		}
		price, ok := APrice.Round(2).Float64()
		if !ok {
			euqiA.WarehousePrice = price
			if err := ac.Handle.HandleMiWarehouse.Update(euqicondition, euqiA); err != nil {
				common.SendResponse(c, http.StatusBadRequest, err, nil)
				return
			}
		}
	}
	euqiB := ac.Handle.HandleMiWarehouse.BeartQuery(euqicondition)
	priceCode[1] = decimal.NewFromFloat(euqiB.WarehousePrice)
	if parameter.BPrice != "" {
		BPrice, err := decimal.NewFromString(parameter.BPrice)
		if err != nil {
			common.SendResponse(c, http.StatusBadRequest, err, nil)
			return
		}
		price, ok := BPrice.Round(2).Float64()
		if !ok {
			euqiB.WarehousePrice = price
			if err := ac.Handle.HandleMiWarehouse.Update(euqicondition, euqiB); err != nil {
				common.SendResponse(c, http.StatusBadRequest, err, nil)
				return
			}
		}
	}

	// 拼凑code)
	aCode := apis.SetPrice2byte(priceCode[0])
	bCode := apis.SetPrice2byte(priceCode[1])

	t := apis.SetPriceCode(aCode + bCode)

	fmt.Println("数据:", t)

	_ = t

	b, _ := hex.DecodeString(t)

	ac.Handle.EventsBus.Publish(events.SendNoticeEvent{Iccid: parameter.Iccid, Data: b})

}

//开关机
func (ac *EquipmentLinkApi) SwitchMachine(c *gin.Context) {
	iccid := c.Query("iccid")
	status, err := strconv.Atoi(c.Query("status"))

	if iccid == "" || err != nil {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if (status != apis.OpenStatus) && (status != apis.CloseStatus) {
		common.SendResponse(c, http.StatusBadRequest, common.NotParamErr, nil)
		return
	}

	if err := ac.Handle.HandleMiEquipment.SwitchMachine(iccid, status); err != nil {
		common.SendResponse(c, http.StatusBadRequest, err, nil)
	}

	// push events sub	   //创建事件
	ac.Handle.EventsBus.Publish(events.SendNoticeEvent{Iccid: iccid, Data: apis.SwitchMac(status)})

	common.SendResponse(c, http.StatusOK, nil, nil)
}
