package tool

type CodeResult int

const (
	CodeSuccess    CodeResult = iota // 成功
	CodeError                        // 失败
	CodeInvalid                      // 无效
	CodeNoAuth                       // 无权限
	CodeParamError                   // 参数错误
	CodeNoLogin                      // 未登录
	CodeNoData                       // 无数据
	CodeFatal                        // 严重错误
)

const (
	Code400 CodeResult = 400 + iota
	Code401
	Code402
	Code403
	Code404
)

const (
	Code500 CodeResult = 500 + iota
	Code501
	Code502
	Code503
	Code504
	Code505
)

// JSONResponse 返回统一的格式
type JSONResponse struct {
	Code    CodeResult `json:"code"`
	Message string    `json:"msg"`
	Data    any       `json:"data"`
}

// SetResult 返回参数赋值
func (rep *JSONResponse) SetResult(code CodeResult, message string) {
	rep.Code = code
	rep.Message = message
}

func NewJSONResponse(code CodeResult, message string) *JSONResponse {
	return &JSONResponse{
		Code:    code,
		Message: message,
	}
}

func NewJSONDataResponse(code CodeResult, message string, data any) *JSONResponse {
	return &JSONResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

type JSONResponseApi struct {
	Code    CodeResult `json:"code"`
	Message string    `json:"msg"`
	Data    any       `json:"data"`
}

// JSONUploadFile 文件上传
type JSONUploadFile struct {
	FileSize int64  `json:"fileSize"` // 文件大小
	FileExt  string `json:"fileExt"`  // 文件后缀
	FileUrl1 string `json:"fileUrl1"` // 文件相对路径
	FileUrl2 string `json:"fileUrl2"` // 文件绝对路径
	FileName string `json:"fileName"` // 保存文件名称
}

// JSONResponsePage 返回分页信息
type JSONResponsePage struct {
	Count int64 `json:"count"`
	JSONResponse
}
