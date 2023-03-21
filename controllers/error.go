package controllers

import "github.com/beego/beego/v2/server/web"

type ErrorController struct {
	web.Controller
}

func (c *ErrorController) Error401() {
	c.Data["content"] = "未经授权，请求要求验证身份"
	c.Data["title"] = "401"
	c.Data["t1"] = 4
	c.Data["t2"] = 0
	c.Data["t3"] = 1
	c.TplName = "error/401.html"
}

func (c *ErrorController) Error403() {
	c.Data["content"] = "服务器拒绝请求"
	c.Data["title"] = "403"
	c.Data["t1"] = 4
	c.Data["t2"] = 0
	c.Data["t3"] = 3
	c.TplName = "error/403.html"
}

func (c *ErrorController) Error404() {
	c.Data["content"] = "很抱歉您访问的地址或者方法不存在"
	c.Data["title"] = "404"
	c.Data["t1"] = 4
	c.Data["t2"] = 0
	c.Data["t3"] = 4
	c.TplName = "error/404.html"
}

func (c *ErrorController) Error500() {
	c.Data["content"] = "server error"
	c.Data["title"] = "500"
	c.Data["t1"] = 5
	c.Data["t2"] = 0
	c.Data["t3"] = 0
	c.TplName = "error/500.html"
}

func (c *ErrorController) Error503() {
	c.Data["content"] = "服务器目前无法使用（由于超载或停机维护）"
	c.Data["title"] = "503"
	c.Data["t1"] = 5
	c.Data["t2"] = 0
	c.Data["t3"] = 3
	c.TplName = "error/503.html"
}
