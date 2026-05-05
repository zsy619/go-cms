package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/zsy619/tools/xgeneric"
	"github.com/zsy619/tools/xphp"
	"github.com/zsy619/tools/xstring"

	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/global"
)

// BaseController 基础控制器,提供公共方法和属性
type BaseController struct {
	web.Controller

	ControllerName string // 控制器名称
	ActionName     string // 方法名称
}

// GetStringTrim 获取字符串并去除首尾空格
// @param key 参数键名
// @param def 默认值
// @return string 处理后的字符串
func (ctrl *BaseController) GetStringTrim(key string, def ...string) string {
	outStr := ctrl.GetString(key, def...)
	outStr = strings.TrimSpace(outStr)
	return outStr
}

// redirect 重定向到指定URL并停止请求处理
// @param url 重定向目标URL
func (ctrl *BaseController) redirect(url string) {
	ctrl.Redirect(url, 302)
	ctrl.StopRun()
}

// IsPost 判断是否为POST请求
// @return bool 是否为POST请求
func (ctrl *BaseController) IsPost() bool {
	return ctrl.Ctx.Request.Method == "POST"
}

// GetPagingParameters 获取分页参数(page和limit)
// @return page 页码,limit 每页数量
func (ctrl *BaseController) GetPagingParameters() (page int, limit int) {
	page, _ = ctrl.GetInt("page", 1)
	limit, _ = ctrl.GetInt("limit", 10)
	return
}

// Prepare 控制器前置处理,设置公共数据和上下文
func (ctrl *BaseController) Prepare() {
	controllerName, actionName := ctrl.GetControllerAndAction()
	ctrl.ControllerName = controllerName[0 : len(controllerName)-10]
	ctrl.ActionName = actionName

	ctrl.Data["version"], _ = web.AppConfig.String("version")
	ctrl.Data["siteName"], _ = web.AppConfig.String("site.name")
	ctrl.Data["curRoute"] = ctrl.ControllerName + "." + ctrl.ActionName
	ctrl.Data["curController"] = ctrl.ControllerName
	ctrl.Data["curAction"] = ctrl.ActionName
	ctrl.Data["time"] = xphp.Time()
	ctrl.Data["year"] = time.Now().Year()
	ctrl.Data["themePath"] = global.ThemePath
}

// Finish 控制器后置处理
func (ctrl *BaseController) Finish() {
	// 预留后置处理逻辑
}

// JSON 公共返回方法(统一响应格式)
// @param code 状态码
// @param message 消息
// @param data 数据
func (ctrl *BaseController) JSON(code lib.CodeResult, message string, data interface{}) {
	ctrl.Data["json"] = &lib.JSONResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
	_ = ctrl.ServeJSON()
	ctrl.StopRun()
}

// JSONSuccess 成功返回
// @param message 消息
// @param data 数据
func (ctrl *BaseController) JSONSuccess(message string, data interface{}) {
	ctrl.JSON(lib.CodeSuccess, message, data)
}

// JSONError 错误返回
// @param message 错误消息
func (ctrl *BaseController) JSONError(message string) {
	ctrl.JSON(lib.CodeError, message, nil)
}

// JSONErrorOfData 带数据的错误返回
// @param message 错误消息
// @param data 数据
func (ctrl *BaseController) JSONErrorOfData(message string, data interface{}) {
	ctrl.JSON(lib.CodeError, message, data)
}

// JSONPage 返回分页信息
// @param code 状态码
// @param message 消息
// @param data 数据列表
// @param count 总数
func (ctrl *BaseController) JSONPage(code lib.CodeResult, message string, data interface{}, count int64) {
	ctrl.Data["json"] = &lib.JSONResponsePage{
		Count: count,
		JSONResponse: lib.JSONResponse{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	_ = ctrl.ServeJSON()
	ctrl.StopRun()
}

// JSONPageSuccess 分页成功返回
// @param data 数据列表
// @param count 总数
func (ctrl *BaseController) JSONPageSuccess(data interface{}, count int64) {
	ctrl.JSONPage(lib.CodeSuccess, "", data, count)
}

// JSONPageError 分页错误返回
// @param msg 错误消息
// @param data 数据
// @param count 总数
func (ctrl *BaseController) JSONPageError(msg string, data interface{}, count int64) {
	ctrl.JSONPage(lib.CodeError, "", data, count)
}

// JSONData 直接返回JSONResponse对象
// @param data JSONResponse对象
func (ctrl *BaseController) JSONData(data *lib.JSONResponse) {
	ctrl.Data["json"] = data
	_ = ctrl.ServeJSON()
	ctrl.StopRun()
}

// OutStatus 输出状态(ajax请求返回JSON,其他请求返回403)
// @param code 状态码
// @param message 消息
// @param data 数据
func (ctrl *BaseController) OutStatus(code lib.CodeResult, message string, data interface{}) {
	if ctrl.Ctx.Input.IsAjax() {
		ctrl.JSON(code, message, data)
	}
	ctrl.Abort("403")
}

// History 历史记录返回(带提示信息)
// @param msg 提示消息
// @param url 跳转URL(为空则返回上一页)
func (ctrl *BaseController) History(msg string, url string) {
	if url == "" {
		ctrl.Ctx.WriteString("<script>alert('" + msg + "');window.history.go(-1);</script>")
		ctrl.StopRun()
	} else {
		ctrl.Redirect(url, 302)
	}
}

// GetClientIp 获取客户端IP地址
// @return string 客户端IP
func (ctrl *BaseController) GetClientIp() string {
	s := strings.Split(ctrl.Ctx.Request.RemoteAddr, ":")
	return s[0]
}

// GetSessionString 获取Session字符串值
// @param sName Session键名
// @return string Session值
func (ctrl *BaseController) GetSessionString(sName string) string {
	fd := ctrl.GetSession(sName)
	if str, ok := fd.(string); ok && str != "" {
		return str
	}
	return ""
}

// GetSessionBool 获取Session布尔值
// @param sName Session键名
// @return bool Session值
func (ctrl *BaseController) GetSessionBool(sName string) bool {
	fd := ctrl.GetSession(sName)
	if str, ok := fd.(string); ok && str != "" {
		rt, _ := strconv.ParseBool(str)
		return rt
	}
	return false
}

// GetSessionInt 获取Session整数值
// @param sName Session键名
// @return int Session值
func (ctrl *BaseController) GetSessionInt(sName string) int {
	fd := ctrl.GetSession(sName)
	if str, ok := fd.(int); ok {
		return str
	}
	return 0
}

// SetSessionBool 设置Session布尔值
// @param sName Session键名
// @param value 值
func (ctrl *BaseController) SetSessionBool(sName string, value bool) {
	if value {
		_ = ctrl.SetSession(sName, "true")
	} else {
		_ = ctrl.SetSession(sName, "false")
	}
}

// SetDatas 批量设置模板数据
// @param datas 数据字典
func (ctrl *BaseController) SetDatas(datas map[string]interface{}) {
	for k, v := range datas {
		ctrl.Data[k] = v
	}
}

// GetSafeString 获取安全的字符串参数(防XSS)
// @param key 参数键名
// @param def 默认值
// @return string 安全的字符串
func (ctrl *BaseController) GetSafeString(key string, def ...string) string {
	data := ctrl.GetStringTrim(key, def...)
	if data == "" {
		return ""
	}
	return xstring.GetSafeString(data)
}

// GetSafeStringTuple2 获取2个安全字符串参数
func (ctrl *BaseController) GetSafeStringTuple2(key1, key2 string) (cnt xgeneric.Tuple2[string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	return
}

// GetSafeStringTuple3 获取3个安全字符串参数
func (ctrl *BaseController) GetSafeStringTuple3(key1, key2, key3 string) (cnt xgeneric.Tuple3[string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	return
}

// GetSafeStringTuple4 获取4个安全字符串参数
func (ctrl *BaseController) GetSafeStringTuple4(key1, key2, key3, key4 string) (cnt xgeneric.Tuple4[string, string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	cnt.D = ctrl.GetSafeString(key4)
	return
}

// GetSafeStringTuple5 获取5个安全字符串参数
func (ctrl *BaseController) GetSafeStringTuple5(key1, key2, key3, key4, key5 string) (cnt xgeneric.Tuple5[string, string, string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	cnt.D = ctrl.GetSafeString(key4)
	cnt.E = ctrl.GetSafeString(key5)
	return
}

// GetSafeStringTuple6 获取6个安全字符串参数
func (ctrl *BaseController) GetSafeStringTuple6(key1, key2, key3, key4, key5, key6 string) (cnt xgeneric.Tuple6[string, string, string, string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	cnt.D = ctrl.GetSafeString(key4)
	cnt.E = ctrl.GetSafeString(key5)
	cnt.F = ctrl.GetSafeString(key6)
	return
}

// GetSafeStringTuple7 获取7个安全字符串参数
func (ctrl *BaseController) GetSafeStringTuple7(key1, key2, key3, key4, key5, key6, key7 string) (cnt xgeneric.Tuple7[string, string, string, string, string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	cnt.D = ctrl.GetSafeString(key4)
	cnt.E = ctrl.GetSafeString(key5)
	cnt.F = ctrl.GetSafeString(key6)
	cnt.G = ctrl.GetSafeString(key7)
	return
}

// GetSafeStringTuple8 获取8个安全字符串参数
func (ctrl *BaseController) GetSafeStringTuple8(key1, key2, key3, key4, key5, key6, key7, key8 string) (cnt xgeneric.Tuple8[string, string, string, string, string, string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	cnt.D = ctrl.GetSafeString(key4)
	cnt.E = ctrl.GetSafeString(key5)
	cnt.F = ctrl.GetSafeString(key6)
	cnt.G = ctrl.GetSafeString(key7)
	cnt.H = ctrl.GetSafeString(key8)
	return
}

// GetSafeStringTuple9 获取9个安全字符串参数
func (ctrl *BaseController) GetSafeStringTuple9(key1, key2, key3, key4, key5, key6, key7, key8, key9 string) (cnt xgeneric.Tuple9[string, string, string, string, string, string, string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	cnt.D = ctrl.GetSafeString(key4)
	cnt.E = ctrl.GetSafeString(key5)
	cnt.F = ctrl.GetSafeString(key6)
	cnt.G = ctrl.GetSafeString(key7)
	cnt.H = ctrl.GetSafeString(key8)
	cnt.I = ctrl.GetSafeString(key9)
	return
}

// GetInt64Tuple2 获取2个int64参数
func (ctrl *BaseController) GetInt64Tuple2(key1, key2 string) (cnt xgeneric.Tuple2[int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	return
}

// GetInt64Tuple3 获取3个int64参数
func (ctrl *BaseController) GetInt64Tuple3(key1, key2, key3 string) (cnt xgeneric.Tuple3[int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	return
}

// GetInt64Tuple4 获取4个int64参数
func (ctrl *BaseController) GetInt64Tuple4(key1, key2, key3, key4 string) (cnt xgeneric.Tuple4[int64, int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	cnt.D, _ = ctrl.GetInt64(key4)
	return
}

// GetInt64Tuple5 获取5个int64参数
func (ctrl *BaseController) GetInt64Tuple5(key1, key2, key3, key4, key5 string) (cnt xgeneric.Tuple5[int64, int64, int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	cnt.D, _ = ctrl.GetInt64(key4)
	cnt.E, _ = ctrl.GetInt64(key5)
	return
}

// GetInt64Tuple6 获取6个int64参数
func (ctrl *BaseController) GetInt64Tuple6(key1, key2, key3, key4, key5, key6 string) (cnt xgeneric.Tuple6[int64, int64, int64, int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	cnt.D, _ = ctrl.GetInt64(key4)
	cnt.E, _ = ctrl.GetInt64(key5)
	cnt.F, _ = ctrl.GetInt64(key6)
	return
}

// GetInt64Tuple7 获取7个int64参数
func (ctrl *BaseController) GetInt64Tuple7(key1, key2, key3, key4, key5, key6, key7 string) (cnt xgeneric.Tuple7[int64, int64, int64, int64, int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	cnt.D, _ = ctrl.GetInt64(key4)
	cnt.E, _ = ctrl.GetInt64(key5)
	cnt.F, _ = ctrl.GetInt64(key6)
	cnt.G, _ = ctrl.GetInt64(key7)
	return
}

// GetInt64Tuple8 获取8个int64参数
func (ctrl *BaseController) GetInt64Tuple8(key1, key2, key3, key4, key5, key6, key7, key8 string) (cnt xgeneric.Tuple8[int64, int64, int64, int64, int64, int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	cnt.D, _ = ctrl.GetInt64(key4)
	cnt.E, _ = ctrl.GetInt64(key5)
	cnt.F, _ = ctrl.GetInt64(key6)
	cnt.G, _ = ctrl.GetInt64(key7)
	cnt.H, _ = ctrl.GetInt64(key8)
	return
}

// GetInt64Tuple9 获取9个int64参数
func (ctrl *BaseController) GetInt64Tuple9(key1, key2, key3, key4, key5, key6, key7, key8, key9 string) (cnt xgeneric.Tuple9[int64, int64, int64, int64, int64, int64, int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	cnt.D, _ = ctrl.GetInt64(key4)
	cnt.E, _ = ctrl.GetInt64(key5)
	cnt.F, _ = ctrl.GetInt64(key6)
	cnt.G, _ = ctrl.GetInt64(key7)
	cnt.H, _ = ctrl.GetInt64(key8)
	cnt.I, _ = ctrl.GetInt64(key9)
	return
}

// GetIntTuple2 获取2个int参数
func (ctrl *BaseController) GetIntTuple2(key1, key2 string) (cnt xgeneric.Tuple2[int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	return
}

// GetIntTuple3 获取3个int参数
func (ctrl *BaseController) GetIntTuple3(key1, key2, key3 string) (cnt xgeneric.Tuple3[int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	return
}

// GetIntTuple4 获取4个int参数
func (ctrl *BaseController) GetIntTuple4(key1, key2, key3, key4 string) (cnt xgeneric.Tuple4[int, int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	cnt.D, _ = ctrl.GetInt(key4)
	return
}

// GetIntTuple5 获取5个int参数
func (ctrl *BaseController) GetIntTuple5(key1, key2, key3, key4, key5 string) (cnt xgeneric.Tuple5[int, int, int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	cnt.D, _ = ctrl.GetInt(key4)
	cnt.E, _ = ctrl.GetInt(key5)
	return
}

// GetIntTuple6 获取6个int参数
func (ctrl *BaseController) GetIntTuple6(key1, key2, key3, key4, key5, key6 string) (cnt xgeneric.Tuple6[int, int, int, int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	cnt.D, _ = ctrl.GetInt(key4)
	cnt.E, _ = ctrl.GetInt(key5)
	cnt.F, _ = ctrl.GetInt(key6)
	return
}

// GetIntTuple7 获取7个int参数
func (ctrl *BaseController) GetIntTuple7(key1, key2, key3, key4, key5, key6, key7 string) (cnt xgeneric.Tuple7[int, int, int, int, int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	cnt.D, _ = ctrl.GetInt(key4)
	cnt.E, _ = ctrl.GetInt(key5)
	cnt.F, _ = ctrl.GetInt(key6)
	cnt.G, _ = ctrl.GetInt(key7)
	return
}

// GetIntTuple8 获取8个int参数
func (ctrl *BaseController) GetIntTuple8(key1, key2, key3, key4, key5, key6, key7, key8 string) (cnt xgeneric.Tuple8[int, int, int, int, int, int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	cnt.D, _ = ctrl.GetInt(key4)
	cnt.E, _ = ctrl.GetInt(key5)
	cnt.F, _ = ctrl.GetInt(key6)
	cnt.G, _ = ctrl.GetInt(key7)
	cnt.H, _ = ctrl.GetInt(key8)
	return
}

// GetIntTuple9 获取9个int参数
func (ctrl *BaseController) GetIntTuple9(key1, key2, key3, key4, key5, key6, key7, key8, key9 string) (cnt xgeneric.Tuple9[int, int, int, int, int, int, int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	cnt.D, _ = ctrl.GetInt(key4)
	cnt.E, _ = ctrl.GetInt(key5)
	cnt.F, _ = ctrl.GetInt(key6)
	cnt.G, _ = ctrl.GetInt(key7)
	cnt.H, _ = ctrl.GetInt(key8)
	cnt.I, _ = ctrl.GetInt(key9)
	return
}

// GetInt32Tuple2 获取2个int32参数
func (ctrl *BaseController) GetInt32Tuple2(key1, key2 string) (cnt xgeneric.Tuple2[int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	return
}

// GetInt32Tuple3 获取3个int32参数
func (ctrl *BaseController) GetInt32Tuple3(key1, key2, key3 string) (cnt xgeneric.Tuple3[int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	return
}

// GetInt32Tuple4 获取4个int32参数
func (ctrl *BaseController) GetInt32Tuple4(key1, key2, key3, key4 string) (cnt xgeneric.Tuple4[int32, int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	cnt.D, _ = ctrl.GetInt32(key4)
	return
}

// GetInt32Tuple5 获取5个int32参数
func (ctrl *BaseController) GetInt32Tuple5(key1, key2, key3, key4, key5 string) (cnt xgeneric.Tuple5[int32, int32, int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	cnt.D, _ = ctrl.GetInt32(key4)
	cnt.E, _ = ctrl.GetInt32(key5)
	return
}

// GetInt32Tuple6 获取6个int32参数
func (ctrl *BaseController) GetInt32Tuple6(key1, key2, key3, key4, key5, key6 string) (cnt xgeneric.Tuple6[int32, int32, int32, int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	cnt.D, _ = ctrl.GetInt32(key4)
	cnt.E, _ = ctrl.GetInt32(key5)
	cnt.F, _ = ctrl.GetInt32(key6)
	return
}

// GetInt32Tuple7 获取7个int32参数
func (ctrl *BaseController) GetInt32Tuple7(key1, key2, key3, key4, key5, key6, key7 string) (cnt xgeneric.Tuple7[int32, int32, int32, int32, int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	cnt.D, _ = ctrl.GetInt32(key4)
	cnt.E, _ = ctrl.GetInt32(key5)
	cnt.F, _ = ctrl.GetInt32(key6)
	cnt.G, _ = ctrl.GetInt32(key7)
	return
}

// GetInt32Tuple8 获取8个int32参数
func (ctrl *BaseController) GetInt32Tuple8(key1, key2, key3, key4, key5, key6, key7, key8 string) (cnt xgeneric.Tuple8[int32, int32, int32, int32, int32, int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	cnt.D, _ = ctrl.GetInt32(key4)
	cnt.E, _ = ctrl.GetInt32(key5)
	cnt.F, _ = ctrl.GetInt32(key6)
	cnt.G, _ = ctrl.GetInt32(key7)
	cnt.H, _ = ctrl.GetInt32(key8)
	return
}

// GetInt32Tuple9 获取9个int32参数
func (ctrl *BaseController) GetInt32Tuple9(key1, key2, key3, key4, key5, key6, key7, key8, key9 string) (cnt xgeneric.Tuple9[int32, int32, int32, int32, int32, int32, int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	cnt.D, _ = ctrl.GetInt32(key4)
	cnt.E, _ = ctrl.GetInt32(key5)
	cnt.F, _ = ctrl.GetInt32(key6)
	cnt.G, _ = ctrl.GetInt32(key7)
	cnt.H, _ = ctrl.GetInt32(key8)
	cnt.I, _ = ctrl.GetInt32(key9)
	return
}