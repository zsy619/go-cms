package www

import "haedu.gov.cn/cms/app/biz"

type TagController struct{ BaseController }

// Index 标签首页
// @router /:flag/tag/:name:string [get]
func (this *TagController) Index() {
	flag := this.Ctx.Input.Param(":flag")
	if flag == "" {
		this.Ctx.WriteString("站点标识不能为空")
	}
	name := this.Ctx.Input.Param(":name")
	if name == "" {
		this.Ctx.WriteString("标签名称不能为空")
	}
	find, err := biz.NewApiTag().Find(0, name)
	if err != nil || find.Name == "" {
		this.Ctx.WriteString("标签不存在")
		this.StopRun()
	}
	this.Data["flag"] = flag
	this.Data["name"] = name
	this.Data["topic"] = find
	if find.Template == "" {
		find.Template = "tag.html"
	}
	this.TplName = this.GetView(SiteTheme, find.Template)
}
