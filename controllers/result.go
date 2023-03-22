package controllers

type ResCode int

const (
	CodeSuccess ResCode = iota
	CodeError
	CodeFail
)

// ResJson 返回统一的格式
type ResJson struct {
	Code    ResCode     `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// PageResJson 返回分页信息
type PageResJson struct {
	Count int64 `json:"count"`
	ResJson
}

type RouterMapping struct {
	Url    string
	Method string
}
