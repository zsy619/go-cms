package www

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
}
