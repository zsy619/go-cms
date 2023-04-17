package www

import "haedu.gov.cn/cms/app/biz"

// ChannelController 频道控制器
type ChannelController struct{ BaseController }

// Index 频道首页
// @router /channel/:name:string [get]
// @router /:name:string [get]
func (this *ChannelController) Index() {
	name := this.Ctx.Input.Param(":name")
	if name == "" {
		this.Abort("404")
	}
	channelModel, err := biz.NewApiChannel().ChannelFind(name, 0)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		this.StopRun()
	}
	this.Ctx.WriteString(channelModel.Name)
	this.Ctx.WriteString("ChannelController.Index")
}

// Category 频道分类
// @router /channel/:name:string/:category:string [get]
// @router /:name:string/:category:string [get]
func (this *ChannelController) Category() {
	name := this.Ctx.Input.Param(":name")
	category := this.Ctx.Input.Param(":category")
	if name == "" || category == "" {
		this.Abort("404")
	}
	channelModel, channelErr := biz.NewApiChannel().ChannelFind(name, 0)
	if channelErr != nil {
		this.Ctx.WriteString(channelErr.Error())
		this.StopRun()
	}
	categoryModel, categoryErr := biz.NewApiArticle().CategoryOne(0, category)
	if categoryErr != nil {
		this.Ctx.WriteString(categoryErr.Error())
		this.StopRun()
	}
	if categoryModel.ChannelID <= 0 {
		this.Ctx.WriteString("分类不存在")
		this.StopRun()
	}
	if categoryModel.ChannelID != channelModel.ChannelID {
		this.Ctx.WriteString("分类与频道不匹配")
		this.StopRun()
	}
	this.Ctx.WriteString(categoryModel.Title)
	// this.Data["channel"] = channel
	// this.Data["category"] = catModel
}
