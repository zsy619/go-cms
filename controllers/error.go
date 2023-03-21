package controllers

import "github.com/beego/beego/v2/server/web"

type ErrorController struct {
	web.Controller
}

func (c *ErrorController) Error401() {
	c.Data["isLogin"] = 0
	c.Data["isCompanyLogin"] = 0
	c.Data["content"] = "未经授权，请求要求验证身份"
	c.TplName = "error/401.html"
}

func (c *ErrorController) Error403() {
	c.Data["isLogin"] = 0
	c.Data["isCompanyLogin"] = 0
	c.Data["content"] = "服务器拒绝请求"
	c.TplName = "error/403.html"
}

func (c *ErrorController) Error404() {
	c.Data["isLogin"] = 0
	c.Data["isCompanyLogin"] = 0
	c.Data["content"] = "很抱歉您访问的地址或者方法不存在"
	c.TplName = "erorr/404.html"
}

func (c *ErrorController) Error500() {
	c.Data["isLogin"] = 0
	c.Data["isCompanyLogin"] = 0
	c.Data["content"] = "server error"
	c.TplName = "error/500.html"
}

func (c *ErrorController) Error503() {
	c.Data["isLogin"] = 0
	c.Data["isCompanyLogin"] = 0
	c.Data["content"] = "服务器目前无法使用（由于超载或停机维护）"
	c.TplName = "error/503.html"
}
