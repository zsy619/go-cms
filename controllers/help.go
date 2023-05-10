package controllers

import "fmt"

type HelpController struct{ BaseController }

func (this *HelpController) Index() {
	this.Data["pageTitle"] = "帮助中心"
	chn := this.Ctx.Input.Param(":chn")
	page := this.Ctx.Input.Param(":page")
	fmt.Println(chn, page)
	tplName := "help/index.html"
	if chn != "" && page != "" {
		tplName = fmt.Sprintf("help/%s/%s.html", chn, page)
	}
	if chn != "" && page == "" {
		tplName = fmt.Sprintf("help/%s/index.html", chn)
	}
	this.Layout = "help/layout.html"
	this.TplName = tplName
}
