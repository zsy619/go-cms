package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
	"haedu.gov.cn/tools/xgeneric"
	"haedu.gov.cn/tools/xphp"
	"haedu.gov.cn/tools/xstring"

	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/global"
)

type BaseController struct {
	web.Controller

	ControllerName string
	ActionName     string
}

func (ctrl *BaseController) GetStringTrim(key string, def ...string) string {
	outStr := ctrl.GetString(key, def...)
	outStr = strings.TrimSpace(outStr)
	return outStr
}

// 重定向
func (ctrl *BaseController) redirect(url string) {
	ctrl.Redirect(url, 302)
	ctrl.StopRun()
}

// 是否POST提交
func (ctrl *BaseController) IsPost() bool {
	return ctrl.Ctx.Request.Method == "POST"
}

// GetPagingParameters 获取分页参数
func (ctrl *BaseController) GetPagingParameters() (page int, limit int) {
	page, _ = ctrl.GetInt("page", 1)
	limit, _ = ctrl.GetInt("limit", 10)
	return
}

func (ctrl *BaseController) Prepare() {
	fmt.Println("BaseController Prepare")
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

func (ctrl *BaseController) Finish() {
	fmt.Println("BaseController Finish")
}

// 公共返回方法
func (ctrl *BaseController) JSON(code lib.CodeResult, message string, data interface{}) {
	ctrl.Data["json"] = &lib.JSONResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
	_ = ctrl.ServeJSON()
	ctrl.StopRun()
}

func (ctrl *BaseController) JSONSuccess(message string, data interface{}) {
	ctrl.JSON(lib.CodeSuccess, message, data)
}

func (ctrl *BaseController) JSONError(message string) {
	ctrl.JSON(lib.CodeError, message, nil)
}

func (ctrl *BaseController) JSONErrorOfData(message string, data interface{}) {
	ctrl.JSON(lib.CodeError, message, data)
}

// JSONPage 返回分页信息
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

// JSONPageSuccess 返回分页信息
func (ctrl *BaseController) JSONPageSuccess(data interface{}, count int64) {
	ctrl.JSONPage(lib.CodeSuccess, "", data, count)
}

// JSONPageError 返回分页信息
func (ctrl *BaseController) JSONPageError(msg string, data interface{}, count int64) {
	ctrl.JSONPage(lib.CodeError, "", data, count)
}

// JSONData 公共返回方法
func (ctrl *BaseController) JSONData(data *lib.JSONResponse) {
	ctrl.Data["json"] = data
	_ = ctrl.ServeJSON()
	ctrl.StopRun()
}

// 公共返回方法
func (ctrl *BaseController) OutStatus(code lib.CodeResult, message string, data interface{}) {
	if ctrl.Ctx.Input.IsAjax() {
		ctrl.JSON(code, message, data)
	}
	ctrl.Abort("403")
}

func (ctrl *BaseController) History(msg string, url string) {
	if url == "" {
		ctrl.Ctx.WriteString("<script>alert('" + msg + "');window.history.go(-1);</script>")
		ctrl.StopRun()
	} else {
		ctrl.Redirect(url, 302)
	}
}

// 获取用户IP地址
func (ctrl *BaseController) GetClientIp() string {
	s := strings.Split(ctrl.Ctx.Request.RemoteAddr, ":")
	return s[0]
}

func (ctrl *BaseController) GetSessionString(sName string) string {
	fd := ctrl.GetSession(sName)
	if str, ok := fd.(string); ok && str != "" {
		return str
	}
	return ""
}

func (ctrl *BaseController) GetSessionBool(sName string) bool {
	fd := ctrl.GetSession(sName)
	if str, ok := fd.(string); ok && str != "" {
		rt, _ := strconv.ParseBool(str)
		return rt
	}
	return false
}

func (ctrl *BaseController) GetSessionInt(sName string) int {
	fd := ctrl.GetSession(sName)
	if str, ok := fd.(int); ok {
		return str
	}
	return 0
}

func (ctrl *BaseController) SetSessionBool(sName string, value bool) {
	if value {
		_ = ctrl.SetSession(sName, "true")
	} else {
		_ = ctrl.SetSession(sName, "false")
	}
}

func (ctrl *BaseController) SetDatas(datas map[string]interface{}) {
	for k, v := range datas {
		ctrl.Data[k] = v
	}
}

// GetSafeString 获取安全字符串
func (ctrl *BaseController) GetSafeString(key string, def ...string) string {
	data := ctrl.GetStringTrim(key, def...)
	if data == "" {
		return ""
	}
	return xstring.GetSafeString(data)
}

func (ctrl *BaseController) GetSafeStringTuple2(key1, key2 string) (cnt xgeneric.Tuple2[string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	return
}

func (ctrl *BaseController) GetSafeStringTuple3(key1, key2, key3 string) (cnt xgeneric.Tuple3[string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	return
}

func (ctrl *BaseController) GetSafeStringTuple4(key1, key2, key3, key4 string) (cnt xgeneric.Tuple4[string, string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	cnt.D = ctrl.GetSafeString(key4)
	return
}

func (ctrl *BaseController) GetSafeStringTuple5(key1, key2, key3, key4, key5 string) (cnt xgeneric.Tuple5[string, string, string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	cnt.D = ctrl.GetSafeString(key4)
	cnt.E = ctrl.GetSafeString(key5)
	return
}

func (ctrl *BaseController) GetSafeStringTuple6(key1, key2, key3, key4, key5, key6 string) (cnt xgeneric.Tuple6[string, string, string, string, string, string]) {
	cnt.A = ctrl.GetSafeString(key1)
	cnt.B = ctrl.GetSafeString(key2)
	cnt.C = ctrl.GetSafeString(key3)
	cnt.D = ctrl.GetSafeString(key4)
	cnt.E = ctrl.GetSafeString(key5)
	cnt.F = ctrl.GetSafeString(key6)
	return
}

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

func (ctrl *BaseController) GetInt64Tuple2(key1, key2 string) (cnt xgeneric.Tuple2[int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	return
}

func (ctrl *BaseController) GetInt64Tuple3(key1, key2, key3 string) (cnt xgeneric.Tuple3[int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	return
}

func (ctrl *BaseController) GetInt64Tuple4(key1, key2, key3, key4 string) (cnt xgeneric.Tuple4[int64, int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	cnt.D, _ = ctrl.GetInt64(key4)
	return
}

func (ctrl *BaseController) GetInt64Tuple5(key1, key2, key3, key4, key5 string) (cnt xgeneric.Tuple5[int64, int64, int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	cnt.D, _ = ctrl.GetInt64(key4)
	cnt.E, _ = ctrl.GetInt64(key5)
	return
}

func (ctrl *BaseController) GetInt64Tuple6(key1, key2, key3, key4, key5, key6 string) (cnt xgeneric.Tuple6[int64, int64, int64, int64, int64, int64]) {
	cnt.A, _ = ctrl.GetInt64(key1)
	cnt.B, _ = ctrl.GetInt64(key2)
	cnt.C, _ = ctrl.GetInt64(key3)
	cnt.D, _ = ctrl.GetInt64(key4)
	cnt.E, _ = ctrl.GetInt64(key5)
	cnt.F, _ = ctrl.GetInt64(key6)
	return
}

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

func (ctrl *BaseController) GetIntTuple2(key1, key2 string) (cnt xgeneric.Tuple2[int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	return
}

func (ctrl *BaseController) GetIntTuple3(key1, key2, key3 string) (cnt xgeneric.Tuple3[int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	return
}

func (ctrl *BaseController) GetIntTuple4(key1, key2, key3, key4 string) (cnt xgeneric.Tuple4[int, int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	cnt.D, _ = ctrl.GetInt(key4)
	return
}

func (ctrl *BaseController) GetIntTuple5(key1, key2, key3, key4, key5 string) (cnt xgeneric.Tuple5[int, int, int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	cnt.D, _ = ctrl.GetInt(key4)
	cnt.E, _ = ctrl.GetInt(key5)
	return
}

func (ctrl *BaseController) GetIntTuple6(key1, key2, key3, key4, key5, key6 string) (cnt xgeneric.Tuple6[int, int, int, int, int, int]) {
	cnt.A, _ = ctrl.GetInt(key1)
	cnt.B, _ = ctrl.GetInt(key2)
	cnt.C, _ = ctrl.GetInt(key3)
	cnt.D, _ = ctrl.GetInt(key4)
	cnt.E, _ = ctrl.GetInt(key5)
	cnt.F, _ = ctrl.GetInt(key6)
	return
}

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

func (ctrl *BaseController) GetInt32Tuple2(key1, key2 string) (cnt xgeneric.Tuple2[int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	return
}

func (ctrl *BaseController) GetInt32Tuple3(key1, key2, key3 string) (cnt xgeneric.Tuple3[int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	return
}

func (ctrl *BaseController) GetInt32Tuple4(key1, key2, key3, key4 string) (cnt xgeneric.Tuple4[int32, int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	cnt.D, _ = ctrl.GetInt32(key4)
	return
}

func (ctrl *BaseController) GetInt32Tuple5(key1, key2, key3, key4, key5 string) (cnt xgeneric.Tuple5[int32, int32, int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	cnt.D, _ = ctrl.GetInt32(key4)
	cnt.E, _ = ctrl.GetInt32(key5)
	return
}

func (ctrl *BaseController) GetInt32Tuple6(key1, key2, key3, key4, key5, key6 string) (cnt xgeneric.Tuple6[int32, int32, int32, int32, int32, int32]) {
	cnt.A, _ = ctrl.GetInt32(key1)
	cnt.B, _ = ctrl.GetInt32(key2)
	cnt.C, _ = ctrl.GetInt32(key3)
	cnt.D, _ = ctrl.GetInt32(key4)
	cnt.E, _ = ctrl.GetInt32(key5)
	cnt.F, _ = ctrl.GetInt32(key6)
	return
}

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
