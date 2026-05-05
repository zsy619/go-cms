package controllers

import (
"fmt"

"github.com/beego/beego/v2/core/logs"
)

type HelpController struct{ BaseController }

func (ctrl *HelpController) Index() {
	ctrl.Data["pageTitle"] = "帮助中心"
	chn := ctrl.Ctx.Input.Param(":chn")
	page := ctrl.Ctx.Input.Param(":page")
	logs.Debug(chn, page)
	tplName := "help/index.html"
	if chn != "" && page != "" {
		tplName = fmt.Sprintf("help/%s/%s.html", chn, page)
	}
	if chn != "" && page == "" {
		tplName = fmt.Sprintf("help/%s/index.html", chn)
	}
	ctrl.Layout = "help/layout.html"
	ctrl.TplName = tplName
}
