package controllers

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

// FilterAdmin 后台管理访问过滤器
// 功能: 检查用户是否已登录,未登录则重定向到登录页面
var FilterAdmin = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("adminId").(int64)
	ok2 := strings.Contains(ctx.Request.RequestURI, "/cms/admin/login")
	if !ok && !ok2 {
		logs.Debug("后台访问拦截: uri=%s, 未登录", ctx.Request.RequestURI)
		ctx.Redirect(302, "/cms/admin/login")
	}
}