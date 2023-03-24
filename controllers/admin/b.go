package admin

import (
	"fmt"
	"haedu.gov.cn/cms/app/lib"

	"haedu.gov.cn/cms/controllers"
)

type BaseController struct {
	controllers.BaseController
}

func (c *BaseController) Prepare() {
	fmt.Println("Admin BaseController Prepare")
	c.BaseController.Prepare()
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

// JSONPaging 返回分页信息
func (c *BaseController) JSONPaging(code lib.CodeResult, message string, data interface{}, count int64) {
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

// JSONData 公共返回方法
func (c *BaseController) JSONData(data *lib.JSONResponse) {
	c.Data["json"] = data
	c.ServeJSON()
	c.StopRun()
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
