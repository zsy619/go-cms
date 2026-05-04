package vmodel

type LoginResult struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Url     string `json:"url"`
}

func NewLoginResult(code int, message string) *LoginResult {
	result := LoginResult{}
	result.Code = code
	result.Message = message
	return &result
}
