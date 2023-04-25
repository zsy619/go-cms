package www

import "haedu.gov.cn/cms/app/biz"

// TopicController 专题控制器
type TopicController struct{ BaseController }

// Index 专题首页
// @router /:flag/topic/:name:string [get]
func (this *TopicController) Index() {
	flag := this.Ctx.Input.Param(":flag")
	if flag == "" {
		this.Ctx.WriteString("站点标识不能为空")
	}
	name := this.Ctx.Input.Param(":name")
	if name == "" {
		this.Ctx.WriteString("专题名称不能为空")
	}
	this.Data["flag"] = flag
	this.Data["name"] = name
	find, err := biz.NewApiTopic().Find(0, name)
	if err != nil {
		this.Ctx.WriteString("专题不存在")
		this.StopRun()
	}
	this.Data["topic"] = find
	if find.Template == "" {
		find.Template = "topic.html"
	}
	this.TplName = this.GetView(SiteTheme, find.Template)
}
