package www

import "haedu.gov.cn/cms/app/cms/service"

type TagController struct{ BaseController }

// Index 标签首页
// @router /:flag/tag/:name:string [get]
func (ctrl *TagController) Index() {
	flag := ctrl.Ctx.Input.Param(":flag")
	if flag == "" {
		ctrl.Ctx.WriteString("站点标识不能为空")
	}
	name := ctrl.Ctx.Input.Param(":name")
	if name == "" {
		ctrl.Ctx.WriteString("标签名称不能为空")
	}
	find, err := service.NewApiTag().Find(0, name)
	if err != nil || find.Name == "" {
		ctrl.Ctx.WriteString("标签不存在")
		ctrl.StopRun()
	}
	ctrl.Data["flag"] = flag
	ctrl.Data["name"] = name
	ctrl.Data["topic"] = find
	if find.Template == "" {
		find.Template = "tag.html"
	}
	ctrl.TplName = ctrl.GetView(SiteTheme, find.Template)
}
