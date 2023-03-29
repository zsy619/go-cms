package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/tools/xphp"
)

type BaseController struct {
	web.Controller

	ControllerName string
	ActionName     string
}

// GetPagingParameters 获取分页参数
func (c *BaseController) GetPagingParameters() (page int, limit int) {
	page, _ = c.GetInt("page", 1)
	limit, _ = c.GetInt("limit", 10)
	return
}

func (c *BaseController) Prepare() {
	fmt.Println("BaseController Prepare")
	controllerName, actionName := c.GetControllerAndAction()
	c.ControllerName = controllerName[0 : len(controllerName)-10]
	c.ActionName = actionName

	c.Data["version"], _ = web.AppConfig.String("version")
	c.Data["siteName"], _ = web.AppConfig.String("site.name")
	c.Data["curRoute"] = c.ControllerName + "." + c.ActionName
	c.Data["curController"] = c.ControllerName
	c.Data["curAction"] = c.ActionName
	c.Data["time"] = xphp.Time()
}

func (c *BaseController) Finish() {
	fmt.Println("BaseController Finish")
}

// 公共返回方法
func (c *BaseController) JSON(code lib.CodeResult, message string, data interface{}) {
	c.Data["json"] = &lib.JSONResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.ServeJSON()
	c.StopRun()
}

func (c *BaseController) JSONSuccess(message string, data interface{}) {
	c.JSON(lib.CodeSuccess, message, data)
}

func (c *BaseController) JSONError(message string) {
	c.JSON(lib.CodeError, message, nil)
}

func (c *BaseController) JSONErrorOfData(message string, data interface{}) {
	c.JSON(lib.CodeError, message, data)
}

// JSONPage 返回分页信息
func (c *BaseController) JSONPage(code lib.CodeResult, message string, data interface{}, count int64) {
	c.Data["json"] = &lib.JSONResponsePage{
		Count: count,
		JSONResponse: lib.JSONResponse{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	c.ServeJSON()
	c.StopRun()
}

// JSONPageSuccess 返回分页信息
func (c *BaseController) JSONPageSuccess(data interface{}, count int64) {
	c.JSONPage(lib.CodeSuccess, "", data, count)
}

// JSONPageError 返回分页信息
func (c *BaseController) JSONPageError(msg string, data interface{}, count int64) {
	c.JSONPage(lib.CodeError, "", data, count)
}

// JSONData 公共返回方法
func (c *BaseController) JSONData(data *lib.JSONResponse) {
	c.Data["json"] = data
	c.ServeJSON()
	c.StopRun()
}

// 公共返回方法
func (c *BaseController) OutStatus(code lib.CodeResult, message string, data interface{}) {
	if c.Ctx.Input.IsAjax() {
		c.JSON(code, message, data)
	}
	c.Abort("403")
}

func (c *BaseController) History(msg string, url string) {
	if url == "" {
		c.Ctx.WriteString("<script>alert('" + msg + "');window.history.go(-1);</script>")
		c.StopRun()
	} else {
		c.Redirect(url, 302)
	}
}

// 获取用户IP地址
func (c *BaseController) GetClientIp() string {
	s := strings.Split(c.Ctx.Request.RemoteAddr, ":")
	return s[0]
}

func (c *BaseController) GetSessionString(sName string) string {
	fd := c.GetSession(sName)
	if str, ok := fd.(string); ok && str != "" {
		return str
	}
	return ""
}

func (c *BaseController) GetSessionBool(sName string) bool {
	fd := c.GetSession(sName)
	if str, ok := fd.(string); ok && str != "" {
		rt, _ := strconv.ParseBool(str)
		return rt
	}
	return false
}

func (c *BaseController) GetSessionInt(sName string) int {
	fd := c.GetSession(sName)
	if str, ok := fd.(int); ok {
		return str
	}
	return 0
}

func (c *BaseController) SetSessionBool(sName string, value bool) {
	if value {
		c.SetSession(sName, "true")
	} else {
		c.SetSession(sName, "false")
	}
}
