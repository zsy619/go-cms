package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type HelpController struct{ BaseController }

// Help 文档手册统一入口
// @router /help/:chn [get]
// @router /help/:chn/:page [get]
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

	// 面包屑: 帮助 / [chn] / [page]
	breadCrumb := "帮助"
	if chn != "" {
		breadCrumb = "帮助 / " + chn
	}
	if page != "" {
		// 去掉 .html 后缀, 取文件名(不含扩展名)作为显示标题
		title := strings.TrimSuffix(page, ".html")
		breadCrumb = fmt.Sprintf("%s / %s", breadCrumb, title)
	}
	ctrl.Data["breadcrumb"] = breadCrumb
	ctrl.Data["year"] = time.Now().Year()
	// 是否首页 (/help/index.html)
	ctrl.Data["isHome"] = chn == ""

	ctrl.Layout = "help/layout.html"
	ctrl.TplName = tplName
}
