package www

// TopicController 专题控制器
type TopicController struct{ BaseController }

// Index 专题首页
// @router /topic/:name:string [get]
func (this *TopicController) Index() {
	name := this.Ctx.Input.Param(":name")
	if name == "" {
		this.Ctx.WriteString("专题名称不能为空")
	}
}
