package httpxfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	hostUrl string
	timeOut string
}

type jsonRPCReq struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type jsonRPCResp struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error"`
	ID      int         `json:"id"`
}

var ErrtoSliceStr = errors.New("json: cannot unmarshal string into Go value of type []string")

func NewClient(url, timeOut string) *Client {
	return &Client{
		hostUrl: url,
		timeOut: timeOut,
	}
}

// CallMethod executes a JSON-RPC call with the given psrameters,which is important to the rpc server.
func (cli *Client) CallMethod(id int, methodname string, params interface{}, out interface{}) error {
	client := resty.New()

	timeDur, err := time.ParseDuration(cli.timeOut)
	if err != nil {
		return err
	}
	client = client.SetTimeout(timeDur)
	req := &jsonRPCReq{
		JsonRPC: "2.0",
		ID:      id,
		Method:  methodname,
		Params:  params,
	}
	// The result must be a pointer so that response json can unmarshal into it.
	var resp *jsonRPCResp = nil
	r, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&resp). // or SetResult(AuthSuccess{}).
		Post(cli.hostUrl)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("resp null")
	}
	e := resp.Error
	if e != nil {
		return e
	}

	js, err := json.Marshal(resp.Result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(js, out)
	if err != nil {
		return err
	}
	_ = r
	return nil
}

// // CallMethod executes a JSON-RPC call with the given psrameters,which is important to the rpc server.
// func (cli *Client) CallMethod(id int, methodname string, params interface{}, out interface{}) error {

// 	timeDur, err := time.ParseDuration(cli.timeOut)
// 	if err != nil {
// 		return err
// 	}

// 	client := &http.Client{Timeout: timeDur}

// 	req := &jsonRPCReq{
// 		JsonRPC: "2.0",
// 		ID:      id,
// 		Method:  methodname,
// 		Params:  params,
// 	}

// 	reqStr, err := json.Marshal(req)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := client.Post(cli.hostUrl, "application/json;charset=utf-8", bytes.NewBuffer(reqStr))
// 	if err != nil {
// 		return err
// 	}
// 	defer result.Body.Close()
// 	// The result must be a pointer so that response json can unmarshal into it.

// 	content, err := ioutil.ReadAll(result.Body)
// 	if err != nil {
// 		return err
// 	}

// 	info := make(map[string]interface{})
// 	if err := json.Unmarshal(content, &info); err != nil {
// 		return err
// 	}

// 	if _, ok := info["result"]; ok {

// 		js, err := json.Marshal(info["result"])
// 		if err != nil {
// 			return err
// 		}
// 		err = json.Unmarshal(js, out)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		js, err := json.Marshal(info["result"])
// 		if err != nil {
// 			return err
// 		}
// 		return fmt.Errorf("err:%v\n", string(js))
// 	}

// 	return nil
// }
