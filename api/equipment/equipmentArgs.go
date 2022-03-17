package api

type SetPriceArgs struct {
	Iccid  string `json:"Iccid"`
	APrice string `json:"aprice"`
	BPrice string `json:"bprice"`
}
