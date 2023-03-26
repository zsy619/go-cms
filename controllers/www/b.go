package www

import (
	"fmt"

	"haedu.gov.cn/cms/controllers"
)

type BaseController struct {
	controllers.BaseController
}

func (c *BaseController) Prepare() {
	fmt.Println("www BaseController Prepare")
	c.BaseController.Prepare()
}

func (c *BaseController) Finish() {
	fmt.Println("www BaseController Finish")
}

// 渲染模版
func (this *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "www/" + this.ControllerName + "/" + this.ActionName + ".html"
	}
	this.Layout = "www/layout/layout.html"
	this.TplName = tplname
}

// 渲染模版
func (this *BaseController) displayNoLayout(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "www/" + this.ControllerName + "/" + this.ActionName + ".html"
	}
	this.TplName = tplname
}
