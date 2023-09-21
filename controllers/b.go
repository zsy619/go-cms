package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/global"
	"haedu.gov.cn/tools/xphp"
	"haedu.gov.cn/tools/xstring"
)

type BaseController struct {
	web.Controller

	ControllerName string
	ActionName     string
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
	data := ctrl.GetString(key, def...)
	if data == "" {
		return ""
	}
	return xstring.GetSafeString(data)
}
