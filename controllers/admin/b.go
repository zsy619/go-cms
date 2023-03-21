package admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers"
)

type BaseController struct {
	web.Controller
	controllerName string
	actionName     string
}

func (c *BaseController) Prepare() {
	fmt.Println("BaseController Prepare")
	controllerName, actionName := c.GetControllerAndAction()
	c.controllerName = controllerName[0 : len(controllerName)-10]
	c.actionName = actionName

	c.Data["version"], _ = web.AppConfig.String("version")
	c.Data["siteName"], _ = web.AppConfig.String("site.name")
	c.Data["curRoute"] = c.controllerName + "." + c.actionName
	c.Data["curController"] = c.controllerName
	c.Data["curAction"] = c.actionName
}

func (c *BaseController) Finish() {
	fmt.Println("BaseController Finish")
}

func (c *BaseController) SavaAdminState(user *model.CmsAdmin) {
	c.SetSession("adminId", user.UserID)
	c.SetSession("adminAccount", user.UserName)
	c.SetSession("adminName", user.UserName)
	c.SetSession("realName", user.RealName)
	c.SetSession("adminState", 1)
	c.SetSession("adminLevel", 1)
	c.SetSession("user", user)

	GlobalAdminId = user.UserID
	GlobalAuthFlag = 1 // 1:管理员 2:学校
	GlobalAdminName = user.UserName
	GlobalRealName = user.RealName
}

// 渲染模版
func (this *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "admin/" + this.controllerName + "/" + this.actionName + ".html"
	}
	this.Layout = "admin/layout/layout.html"
	this.TplName = tplname
}

// 渲染模版
func (this *BaseController) displayNoLayout(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "admin/" + this.controllerName + "/" + this.actionName + ".html"
	}
	this.TplName = tplname
}

// 公共返回方法
func (c *BaseController) OutJson(code int, message string, data interface{}) {
	c.Data["json"] = &controllers.ResJson{
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.ServeJSON()
	c.StopRun()
}

// OutPageJson 返回分页信息
func (c *BaseController) OutPageJson(code int, message string, count int64, data interface{}) {
	c.Data["json"] = &controllers.PageResJson{
		Count: count,
		ResJson: controllers.ResJson{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	c.ServeJSON()
	c.StopRun()
}

// 公共返回方法
func (c *BaseController) OutStatus(code int, message string, data interface{}) {
	if c.Ctx.Input.IsAjax() {
		c.OutJson(code, message, data)
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
