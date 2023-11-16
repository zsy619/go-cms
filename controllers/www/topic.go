package www

import "haedu.gov.cn/cms/app/biz"

// TopicController 专题控制器
type TopicController struct{ BaseController }

// Index 专题首页
// @router /:flag/topic/:name:string [get]
func (ctrl *TopicController) Index() {
	flag := ctrl.Ctx.Input.Param(":flag")
	if flag == "" {
		ctrl.Ctx.WriteString("站点标识不能为空")
	}
	name := ctrl.Ctx.Input.Param(":name")
	if name == "" {
		ctrl.Ctx.WriteString("专题名称不能为空")
	}
	find, err := biz.NewApiTopic().Find(0, name)
	if err != nil || find.Name == "" {
		ctrl.Ctx.WriteString("专题不存在")
		ctrl.StopRun()
	}
	ctrl.Data["flag"] = flag
	ctrl.Data["name"] = name
	ctrl.Data["topic"] = find
	if find.Template == "" {
		find.Template = "topic.html"
	}
	ctrl.TplName = ctrl.GetView(SiteTheme, find.Template)
}
