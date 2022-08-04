package model

// 分页参数
type PageParam struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Params   string `json:"params"`
}

type ResultPage struct {
	Total    int16 `json:"total"`
	PageNum  int16 `json:"pageNum"`
	PageSize int16 `json:"pageSize"`
	Data     any   `json:"data"`
}
