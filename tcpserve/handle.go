package tcpserve

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"mi/common/apis"
	"mi/common/crc16"
	"mi/events"
	"mi/global"
	"mi/model"
	"mi/pay/wechatpay"
	"mi/tcpserve/common"
	"net"

	"github.com/shopspring/decimal"
)

type Handle struct {
	model          *model.RecordHandle
	wechatpayServe *wechatpay.WeChatPay
	iccid          string
	tcpConn        net.Conn
	eventBus       *events.EventBus
}

func NewHandle() *Handle {
	handle := &Handle{
		model:          model.NewRecordHandle(),
		eventBus:       events.EventBusExample,
		wechatpayServe: wechatpay.DefaultWeChatPay(),
	}
	// 处理广播
	go handle.MsgBroadcastLoop()
	return handle
}

// func (h *Handle)
func (h *Handle) Process(conn net.Conn) error {
	//打开注释的地方(devour)
	// defer conn.Close()
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
		nowtime := time.Now()
		h.tcpConn.SetDeadline(nowtime.Add(10 * time.Minute)) //设置刷新超时
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
		}
	}
}

func (h *Handle) chck(data []byte) ([]byte, error) {

	//验证数据是否AA开头
	if data[0] != common.Header {
		return nil, errors.New("header data error")
	}

	//验证crc16
	crc16 := crc16.CRC(data[:len(data)-2])
	crc16msg := data[len(data)-2:]
	//bert end

	if bytes.Compare(crc16, crc16msg) != int(0) {
		return nil, errors.New("header crc error")
	}

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
	earlyA := h.model.HandleWlEarly.Query("iccid =?", h.iccid, "warehouse=?", 1)

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
	earlyB := h.model.HandleWlEarly.Query("iccid =?", h.iccid, "warehouse=?", 2)

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

func (h *Handle) WechatPay(amounts int, auto_code string) ([]byte, error) {

	//获取本机服务器IP
	ip, err := common.GetLocalIP()
	if err != nil {
		return UploadOrderError, err
	}
	//传入参数付款支付逻辑
	_, ordernumber, err := h.wechatpayServe.Micropay(amounts, auto_code, "mi", ip)
	if err != nil {
		return UploadOrderError, err
	}

	chanErr := make(chan error)
	//等待支付结果 15s
	timeout := time.After(time.Second * 15) //等待15秒在
	go func() {
		for {
			select {
			case <-timeout:
				//查询支付结果
				query, _, err := h.wechatpayServe.OrderQuery(ordernumber)
				if err != nil {
					chanErr <- err
					return
				}
				//判断支付结果
				// 交易成功判断条件： return_code、result_code和trade_state都为SUCCESS
				if query.ReturnCode != "SUCCESS" || query.ResultCode != "SUCCESS" || query.TradeState != "SUCCESS" {
					chanErr <- errors.New(query.ReturnMsg) //错误原因
					return
				}

				//订单记录写入数据库
				write := &model.MiOrder{
					Iccid:         h.iccid,
					PayType:       18,
					PayCode:       "testbert",  //唯一标识
					OrderNumber:   ordernumber, //订单号
					PaymentAmount: decimal.NewFromInt(int64(amounts)),
					Number:        auto_code, //付款码
				}

				if err := h.model.HandleMiOrder.Insert(write); err != nil {
					chanErr <- err
					return
				}
				//所有操作正确且成功
				chanErr <- nil

			}
		}
	}()

	//验证付款结果
	select {
	case err := <-chanErr:
		if err != nil {
			// global.GVA_LOG.Warn(err.Error())
			return UploadOrderError, err
		} else {
			return UploadOrderSucc, nil
		}
	}
}

func (h *Handle) BalancePay(amounts uint64, cardnumber string) ([]byte, error) {
	// 卡号查询余额
	cardingo := h.model.HandleWlCardNumber.Query("number = ?", cardnumber)
	if cardingo == nil {
		global.GVA_LOG.Warn("卡号不存在")
		return UploadOrderError, nil
	}
	// 对比余额
	amount := decimal.NewFromInt(int64(amounts)).Mul(decimal.NewFromInt(100))
	if amount.Cmp(cardingo.Balance) < 1 {
		global.GVA_LOG.Warn(fmt.Sprintf("用户id:%v err:余额不足", cardingo.MemberId))
		return UploadOrderError, nil
	}
	// 扣除余额
	cardingo.Balance = cardingo.Balance.Sub(decimal.NewFromInt(int64(amounts)))
	// 开启事物扣除卡号表余额度
	if err := h.model.HandleWlCardNumber.Update("member_id = ?", cardingo.MemberId, cardingo); err != nil {
		global.GVA_LOG.Warn(err.Error())
		return UploadOrderError, err
	}

	//记录支付信息写入数据库
	write := &model.MiOrder{
		Iccid:         h.iccid,
		PayType:       10,
		PayCode:       strconv.Itoa(cardingo.MemberId), //唯一标识
		OrderNumber:   common.GetOrder(),               //订单号
		PaymentAmount: amount.Div(decimal.NewFromInt(100)).Round(2),
		Number:        cardnumber, //付款码
	}
	if err := h.model.HandleMiOrder.Insert(write); err != nil {
		global.GVA_LOG.Warn(err.Error())
		//返回错误码
		return UploadOrderError, err
	}
	return UploadOrderSucc, nil
}

//上传订单
func (h *Handle) uploadOrder(data []byte) ([]byte, error) {
	//Ascii码
	amounts := common.Hex2int(&[]byte{data[0], data[1]}) // 1045

	deCodePay := make([]rune, 0)
	for _, v := range data {
		deCodePay = append(deCodePay, rune(v))
	}
	//付款码或者卡号
	payCode := string(deCodePay)

	// 微信支付付款码
	found, err := regexp.MatchString(`^1[0-5]\d{16}$`, payCode)
	if err != nil {
		global.GVA_LOG.Warn(err.Error())
	}
	// 微信支付
	if found {
		return h.WechatPay(int(amounts), payCode)
	}
	// 卡号支付
	return h.BalancePay(amounts, payCode)
}

func (h *Handle) setPrice(data []byte) ([]byte, error) {
	list := h.model.HandleMiEquipment.Query("iccid =?", h.iccid)

	priceCode := make([]decimal.Decimal, 2)
	priceCode[0] = decimal.NewFromFloat(list.AWarehousePrice)
	priceCode[1] = decimal.NewFromFloat(list.BWarehousePrice)

	aCode := apis.SetPrice2byte(priceCode[0])
	bCode := apis.SetPrice2byte(priceCode[1])
	code, err := apis.SetPriceCode(aCode + bCode)
	if err != nil {
		return nil, err
	}
	bs, err := hex.DecodeString(code)
	if err != nil {
		return nil, err
	}

	fmt.Println(bs)
	return bs, nil
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

	if err := h.model.HandleMiEquipment.SetSwitc("iccid =?", switchss); err != nil {
		return nil, nil
	}

	if switchss == 0 {
		defer h.tcpConn.Close()
		return SwitchCloseError, nil
	}

	if switchss == 1 {
		return SwitchOpenSucc, nil
	}
	return nil, nil
}
