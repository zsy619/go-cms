package admin

import (
	"fmt"

	"haedu.gov.cn/cms/app/dal/model"

	"haedu.gov.cn/cms/controllers"
)

type BaseController struct {
	controllers.BaseController
}

func (c *BaseController) Prepare() {
	fmt.Println("Admin BaseController Prepare")
	c.BaseController.Prepare()

	if GlobalAdminId == 0 {
		user := c.GetSession("user").(*model.CmsAdmin)
		GlobalAdminId = user.UserID
		GlobalAuthFlag = int(user.UserType) // 1:管理员 2:学校
		GlobalAdminName = user.UserName
		GlobalRealName = user.RealName
	}
}

func (c *BaseController) Finish() {
	fmt.Println("Admin BaseController Finish")
}

// 渲染模版
func (this *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "admin/" + this.ControllerName + "/" + this.ActionName + ".html"
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
		tplname = "admin/" + this.ControllerName + "/" + this.ActionName + ".html"
	}
	this.TplName = tplname
}

// 登录人ID
func (this *BaseController) IsLogin() int64 {
	id := this.GetSession(`adminId`)
	if id == nil {
		return 0
	} else {
		switch id.(type) {
		case int64:
			rt := id.(int64)
			return rt
		default:
			return 0
		}
	}
}
