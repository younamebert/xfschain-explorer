package tcpserve

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"mi/common/crc16"
	"mi/events"
	"mi/global"
	"mi/model"
	"mi/tcpserve/common"
	"mi/tools"
	"net"

	"github.com/shopspring/decimal"
)

type Handle struct {
	model    *model.RecordHandle
	iccid    string
	tcpConn  net.Conn
	outTime  *time.Timer
	eventBus *events.EventBus
}

func NewHandle() *Handle {
	handle := &Handle{
		model:    model.NewRecordHandle(),
		eventBus: events.EventBusExample,
		outTime:  time.NewTimer(30 * time.Minute), //30分钟
	}
	// 处理广播
	go handle.MsgBroadcastLoop()
	return handle
}

// func (h *Handle)
func (h *Handle) Process(conn net.Conn) error {

	//打开注释的地方(devour)
	// defer conn.Close()
	// go /
	h.tcpConn = conn

	for {
		// h.MsgBroadcastLoop()
		reader := bufio.NewReader(conn)
		buf := make([]byte, 128)
		n, err := reader.Read(buf) // 读取数据
		if err != nil {
			global.GVA_LOG.Warn(err.Error())
			return err
		}
		result := buf[:n]
		write, err := h.chck(result)
		if err != nil {
			global.GVA_LOG.Warn(err.Error())
			return err
		}
		if _, err := conn.Write(write); err != nil {
			global.GVA_LOG.Warn(err.Error())
			return fmt.Errorf("write client msg, err: %v", err)
		}
	}
	return nil
}

func (h *Handle) resetOutTime(msgcode byte) {
	if msgcode != common.Pant {
		h.outTime.Reset(20 * time.Minute)
	}
}

func (h *Handle) MsgBroadcastLoop() {
	newEventSub := h.eventBus.Subscript(events.SendNoticeEvent{})
	defer newEventSub.Unsubscribe()
	for {
		// time.Sleep(time.Second * 1)
		select {
		case e := <-newEventSub.Chan():
			event := e.(events.SendNoticeEvent)
			if h.iccid == event.Iccid {
				msg, err := h.chck(event.Data)
				if err != nil {
					h.tcpConn.Close()
					return
				}
				if _, err := h.tcpConn.Write(msg); err != nil {
					global.GVA_LOG.Warn(err.Error())
					return
				}
			}
		case outTime := <-h.outTime.C:
			global.GVA_LOG.Error(fmt.Sprintf("Device timeout session iccid:%v outitme:%v", h.iccid, outTime.Unix()))
			h.tcpConn.Close()
		}
	}
}

// func (h *Handle)
func (h *Handle) chck(data []byte) ([]byte, error) {

	//验证数据是否AA开头
	if data[0] != common.Header {
		return nil, errors.New("header data error")
	}

	fmt.Println("crc校验:", data)

	//验证crc16
	crc16 := crc16.CRC(data[:len(data)-2])
	crc16msg := data[len(data)-2:]
	//bert end

	fmt.Println("检验:", crc16)

	if bytes.Compare(crc16, crc16msg) != int(0) {
		return nil, errors.New("header crc error")
	}

	// AA F3 04 01 60 01 60 02ED

	//验证数据长度
	dataLen := data[2:3]
	pending := data[3 : len(data)-2]

	if int(dataLen[0]) != len(pending) {
		return nil, errors.New("header data len error")
	}

	// if int(data[2])
	// 优化
	result, err := h.work(data[1], pending)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// arrary("nihao"=>["nihao":1])
func (h *Handle) work(funccode byte, data []byte) ([]byte, error) {
	if h.iccid == "" {
		if funccode != common.Registers {
			return nil, errors.New("no iccid, please register first")
		}
	}

	h.resetOutTime(funccode)
	switch funccode {
	case common.Registers:
		return h.registers(data)
	case common.Pant:
		return h.pant(data)
	case common.UploadOrder:
		return h.uploadOrder(data)
	case common.SetPrice:
		return h.setPrice(data)
	case common.Switchad:
		return h.switchad(data)
	case common.SwitchLamp:
		return h.switchLamp(data)
	case common.Switchs:
		return h.switchs(data)
	}

	return nil, nil
}

func (h *Handle) registers(data []byte) ([]byte, error) {
	iccdRune := make([]rune, 0)
	for _, v := range data {
		iccdRune = append(iccdRune, rune(v))
	}
	//Ascii码
	iccid := string(iccdRune)

	write := &model.MiEquipment{
		Iccid: iccid,
	}

	//判断数据存不存在 不存在走这不
	list := h.model.HandleMiEquipment.Query("iccid =?", iccid)

	if list.Iccid == "" {
		if err := h.model.HandleMiEquipment.Insert(write); err != nil {
			fmt.Println(err)
			return AddEquipmentRegistersErr, err
		}
	}
	h.iccid = iccid
	// 写入数据库
	return AddEquipmentRegisters, nil
}

func (h *Handle) pant(data []byte) ([]byte, error) {
	var (
		wareA int = 0 // A仓库
		wareB int = 0 // B仓库
	)
	fmt.Println(h.iccid)
	wareA = common.WareType(data[0])
	wareB = common.WareType(data[1])
	// 写入数据库(a,b仓库)
	a := wareA
	b := wareB

	//查询数据库存不存在

	EmptyA := h.model.HandleMiWarehouse.Query("iccid =? and status=?", h.iccid, 1)
	//查询设备的a,b仓价格
	list := h.model.HandleMiEquipment.Query("iccid =?", h.iccid)

	table1 := &model.MiWarehouse{
		Iccid:          h.iccid,
		WarehouseType:  a,
		WarehousePrice: list.AWarehousePrice,
		Status:         1,
	}
	if EmptyA.Iccid == "" {
		h.model.HandleMiWarehouse.Insert(table1)

	} else {
		//修改

		EmptyA.WarehousePrice = list.AWarehousePrice
		EmptyA.WarehouseType = a
		h.model.HandleMiWarehouse.SaveWare("iccid =?", h.iccid, "status=?", 1, EmptyA)
	}

	EmptyB := h.model.HandleMiWarehouse.Query("iccid =? and status=?", h.iccid, 2)
	table2 := &model.MiWarehouse{
		Iccid:          h.iccid,
		WarehouseType:  b,
		WarehousePrice: list.BWarehousePrice,
		Status:         2,
	}
	if EmptyB.Iccid == "" {
		h.model.HandleMiWarehouse.Insert(table2)
	} else {
		//修改

		EmptyB.WarehousePrice = list.BWarehousePrice
		EmptyB.WarehouseType = b
		h.model.HandleMiWarehouse.SaveWare("iccid =?", h.iccid, "status=?", 2, EmptyB)
	}

	mangeList := h.model.HandleWlMange.Query("iccid =?", h.iccid)

	//a仓位预警
	earlyA := h.model.HandleWlEarly.QueryOne("iccid =?", h.iccid, "warehouse=?", 1)

	if earlyA.Iccid == "" {

		earlyA_inert := &model.WlEarly{
			MangeId:   int(mangeList.MangeId),
			Number:    h.iccid,
			Type:      1,
			CStatus:   a,
			Warehouse: 1,
			Iccid:     h.iccid,
		}

		h.model.HandleWlEarly.Insert(earlyA_inert)
	} else {
		earlyA.CStatus = a
		h.model.HandleWlEarly.SaveEarly("iccid =?", h.iccid, "warehouse=?", 1, earlyA)
	}

	//b仓库预警
	earlyB := h.model.HandleWlEarly.QueryOne("iccid =?", h.iccid, "warehouse=?", 2)

	if earlyB.Iccid == "" {

		earlyB_inert := &model.WlEarly{
			MangeId:   int(mangeList.MangeId),
			Number:    h.iccid,
			Type:      1,
			CStatus:   b,
			Warehouse: 2,
			Iccid:     h.iccid,
		}

		h.model.HandleWlEarly.Insert(earlyB_inert)
	} else {
		earlyB.CStatus = b
		h.model.HandleWlEarly.SaveEarly("iccid =?", h.iccid, "warehouse=?", 2, earlyB)

	}

	return HeartbeatSucc, nil
}

//上传订单
func (h *Handle) uploadOrder(data []byte) ([]byte, error) {
	var (
		amounts decimal.Decimal
		payType int
		payCode string
	)

	moneyUint64 := common.Hex2int(&[]byte{data[0], data[1]}) // 1045
	amounts = common.Uint64toDecimal(int64(moneyUint64), 100)

	_, _ = payType, payCode

	cardNumber := h.model.HandleWlCardNumber.Query("number =?", "8939131")

	//余额不足
	money := decimal.NewFromFloat(cardNumber.Money).Round(2).BigFloat()

	//余额不足
	if money.Cmp(amounts.BigFloat()) < 0 {
		return UploadOrderError, nil
	}

	switch len(data[1:]) {

	//卡片支付
	case common.Card:
		fmt.Println(1314412)
		break

	//其它支付
	case common.Other:

		//微信支付

		break
	}

	write := &model.MiOrder{
		Iccid:         h.iccid,
		PayType:       len(data),
		PayCode:       strconv.FormatInt(time.Now().Unix()+int64(rand.Intn(100000)), 10),
		OrderNumber:   tools.GetOrder(),
		PaymentAmount: amounts,
		Number:        "8939131",
	}

	//查询设备id
	list := h.model.HandleWlMange.Query("number =?", h.iccid)

	amouts, _ := amounts.Float64()

	//设备出售米纪录表
	saleWari := &model.WlSale{
		MangeId: list.MangeId,
		Number:  2,
		Money:   amouts,
	}

	//购米记录表
	writes := &model.WlPurchaseRecord{
		Price:    amouts,
		OrderNo:  tools.GetOrder(),
		PayTime:  time.Now(),
		MemberId: cardNumber.MemberId,
		Iccid:    h.iccid,
		MangeId:  list.MangeId,
	}

	//余额足支付
	if money.Cmp(amounts.BigFloat()) > 0 {
		if err := h.model.HandleMiOrder.Insert(write); err != nil {
			global.GVA_LOG.Warn(err.Error())

			//返回错误码
			return UploadOrderSucc, err
		}

		//设备出售米纪录表
		if err := h.model.HandleWlSale.Insert(saleWari); err != nil {
			global.GVA_LOG.Warn(err.Error())
		}

		//加入购米记录表
		h.model.HandleWlPurchaseRecord.Insert(writes)

		if err := h.model.HandleWlSale.Insert(saleWari); err != nil {
			global.GVA_LOG.Warn(err.Error())
		}

		//扣余额
		money.Sub(money, amounts.BigFloat())

		moneyFloat, _ := money.Float64()

		cardNumber.Money = moneyFloat

		h.model.HandleWlCardNumber.Update("card_number_id=?", cardNumber.CardNumberId, cardNumber)

		return UploadOrderSucc, nil
	}

	if money.Cmp(amounts.BigFloat()) == 0 {
		if err := h.model.HandleMiOrder.Insert(write); err != nil {
			// 	//返回错误码
			global.GVA_LOG.Warn(err.Error())
			return AddEquipmentRegistersErr, err
		}

		//加入购米记录表
		if err := h.model.HandleWlPurchaseRecord.Insert(writes); err != nil {
			global.GVA_LOG.Warn(err.Error())
		}

		//设备出售米纪录表
		h.model.HandleWlSale.Insert(saleWari)
		moneyFloat, _ := money.Float64()
		cardNumber.Money = moneyFloat
		h.model.HandleWlCardNumber.Update("card_number_id=?", cardNumber.CardNumberId, cardNumber)

		return UploadOrderSucc, nil
	}

	return UploadOrderSucc, nil
}

func (h *Handle) setPrice(data []byte) ([]byte, error) {

	fmt.Println(3131)

	// var (
	// 	wareAPrice decimal.Decimal
	// 	wareBPrice decimal.Decimal
	// )

	// wareAPriceUint64 := common.Hex2int(&[]byte{data[0], data[1]}) //350

	// wareAPrice = common.Uint64toDecimal(int64(wareAPriceUint64), 100) //3.50

	// wareBPriceUint64 := common.Hex2int(&[]byte{data[2], data[3]}) //
	// wareBPrice = common.Uint64toDecimal(int64(wareBPriceUint64), 100)

	// // fmt.Printf("a:%v b:%v\n", wareAPrice.String(), wareBPrice.String())
	// whereA := make([]interface{}, 0)
	// whereA = append(whereA, h.Iccid, 1)
	// warehouseA := h.model.HandleMiWarehouse.Query("iccid = ? and warehouse_type = ?", whereA)
	// if len(warehouseA) > 0 {
	// 	warehouseAWrite := warehouseA[0]
	// 	rPrice, _ := wareAPrice.Round(2).Float64()
	// 	warehouseAWrite.WarehousePrice = rPrice
	// 	if err := h.model.HandleMiWarehouse.Update(warehouseAWrite); err != nil {
	// 		//错误的码
	// 		// common.ModifyError
	// 		return ModifyError, nil
	// 	}
	// } else {
	// 	warehouseAWrite := &model.MiWarehouse{
	// 		Iccid:          h.Iccid,
	// 		WarehouseType:  1,
	// 		WarehousePrice: float64(0),
	// 	}
	// 	if err := h.model.HandleMiWarehouse.Insert(warehouseAWrite); err != nil {
	// 		//错误的码
	// 		return ModifyError, nil
	// 	}
	// }

	// whereB := make([]interface{}, 0)
	// whereB = append(whereB, h.Iccid, 0)
	// warehouseB := h.model.HandleMiWarehouse.Query("iccid = ? and warehouse_type = ?", whereB)
	// if len(warehouseB) > 0 {
	// 	warehouseBWrite := warehouseB[0]
	// 	rPrice, _ := wareBPrice.Round(2).Float64()
	// 	warehouseBWrite.WarehousePrice = rPrice
	// 	if err := h.model.HandleMiWarehouse.Update(warehouseBWrite); err != nil {
	// 		//错误的码
	// 		return ModifyError, nil
	// 	}
	// } else {
	// warehouseBWrite := &model.MiWarehouse{
	// 		Iccid:          h.Iccid,
	// 		WarehouseType:  1,
	// 		WarehousePrice: float64(0),
	// 	}
	// 	if err := h.model.HandleMiWarehouse.Insert(warehouseBWrite); err != nil {
	// 		//错误的码
	// 		return ModifyError, nil
	// 	}
	// }
	return nil, nil
}

func (h *Handle) switchad(data []byte) ([]byte, error) {
	var switchadType int = 0
	switch data[0] {
	case byte(0x00):
		switchadType = 0
	case byte(0xFF):
		switchadType = 1
	}

	if err := h.model.HandleMiEquipment.SetSwitchad("iccid =?", switchadType); err != nil {
		return nil, nil
	}

	if switchadType == 0 {
		return SwitchAdvertisingError, nil
	}

	if switchadType == 1 {
		return SwitchAdvertisingSucc, nil
	}

	return nil, nil
}

func (h *Handle) switchLamp(data []byte) ([]byte, error) {
	var switchLampType int = 0
	switch data[0] {
	case byte(0x00):
		switchLampType = 0
	case byte(0xFF):
		switchLampType = 1
	}

	if err := h.model.HandleMiEquipment.SetSwitchadLed("iccid =?", switchLampType); err != nil {
		return nil, nil
	}

	if switchLampType == 0 {
		return SwitchBeltError, nil
	}

	if switchLampType == 1 {
		return SwitchBeltSucc, nil
	}

	return nil, nil

}

// []byte{0xAA, 0xF6, 0x03, 0x49, 0x51, 0x55, 0x02, 0x92}
// []byte{0xAA, 0xF6, 0x03, 0x49, 0x50, 0x51, 0x02, 0x94}
//开关机
func (h *Handle) switchs(data []byte) ([]byte, error) {

	var switchss int = 0

	switch data[1] {
	case byte(0x51):
		switchss = 1
	case byte(0x50):
		switchss = 0
	}

	fmt.Println("结果:", switchss)

	if err := h.model.HandleMiEquipment.SetSwitc("iccid =?", switchss); err != nil {
		return nil, nil
	}

	if switchss == 0 {
		h.tcpConn.Close()
		return SwitchCloseError, nil
	}

	if switchss == 1 {
		return SwitchOpenSucc, nil
	}
	return nil, nil
}
