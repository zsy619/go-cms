package models

type WxPayException struct {
	Code    int
	Message string
}

func NewWxPayException(message string) WxPayException {
	result := WxPayException{}
	result.Message = message
	return result
}

func (wpe WxPayException) Error() string {
	return wpe.Message
}
