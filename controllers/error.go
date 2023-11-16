package controllers

import "github.com/beego/beego/v2/server/web"

type ErrorController struct{ web.Controller }

func (ctrl *ErrorController) Error401() {
	ctrl.Data["content"] = "未经授权，请求要求验证身份"
	ctrl.Data["title"] = "401"
	ctrl.Data["t1"] = 4
	ctrl.Data["t2"] = 0
	ctrl.Data["t3"] = 1
	ctrl.TplName = "error/401.html"
}

func (ctrl *ErrorController) Error403() {
	ctrl.Data["content"] = "服务器拒绝请求"
	ctrl.Data["title"] = "403"
	ctrl.Data["t1"] = 4
	ctrl.Data["t2"] = 0
	ctrl.Data["t3"] = 3
	ctrl.TplName = "error/403.html"
}

func (ctrl *ErrorController) Error404() {
	ctrl.Data["content"] = "很抱歉您访问的地址或者方法不存在"
	ctrl.Data["title"] = "404"
	ctrl.Data["t1"] = 4
	ctrl.Data["t2"] = 0
	ctrl.Data["t3"] = 4
	ctrl.TplName = "error/404.html"
}

func (ctrl *ErrorController) Error500() {
	ctrl.Data["content"] = "server error"
	ctrl.Data["title"] = "500"
	ctrl.Data["t1"] = 5
	ctrl.Data["t2"] = 0
	ctrl.Data["t3"] = 0
	ctrl.TplName = "error/500.html"
}

func (ctrl *ErrorController) Error503() {
	ctrl.Data["content"] = "服务器目前无法使用（由于超载或停机维护）"
	ctrl.Data["title"] = "503"
	ctrl.Data["t1"] = 5
	ctrl.Data["t2"] = 0
	ctrl.Data["t3"] = 3
	ctrl.TplName = "error/503.html"
}
