package apis

type Pages struct {
	Page   int64 `json:"page"`
	Limits int64 `json:"limits"`
	// Total    int64       `json:"total"`
	PageSize int64       `json:"page_size"`
	Data     interface{} `json:"data"`
}
