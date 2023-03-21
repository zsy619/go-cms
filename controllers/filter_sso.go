package controllers

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// FilterAdmin 过滤器登录判断
var FilterAdmin = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("a.id").(string)
	fmt.Println(ctx.Request.RequestURI)
	ok2 := strings.Contains(ctx.Request.RequestURI, "/cms/admin/login")
	if !ok && !ok2 {
		ctx.Redirect(302, "/cms/admin/login")
	}
}
