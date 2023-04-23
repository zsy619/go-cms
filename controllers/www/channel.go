package www

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/tools/xstring"
)

// ChannelController 频道控制器
type ChannelController struct{ BaseController }

// Channel 频道首页
// @router /:flag/:name:string [get]
func (this *ChannelController) Channel() {
	siteFlag := this.Ctx.Input.Param(":flag")
	chanelName := this.Ctx.Input.Param(":name")
	if siteFlag == "" || chanelName == "" {
		this.Abort("404")
	}
	channelModel, err := biz.NewApiChannel().Find(chanelName, 0)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		this.StopRun()
	}
	this.Data["channel"] = channelModel
	this.Data["chanelName"] = chanelName
	if channelModel.Template == "" {
		this.TplName = this.getView(DefatulSite.Template, "channel.html")
	} else {
		if xstring.HasSuffix(channelModel.Template, ".html", ".htm", ".tpl") == false {
			channelModel.Template += ".html"
		}
		this.TplName = this.getView(DefatulSite.Template, channelModel.Template)
	}
}

// Category 频道分类
// @router /:flag/:name:string/:category:string [get]
func (this *ChannelController) Category() {
	siteFlag := this.Ctx.Input.Param(":flag")
	chanelName := this.Ctx.Input.Param(":name")
	categoryName := this.Ctx.Input.Param(":category")
	if siteFlag == "" || chanelName == "" || categoryName == "" {
		this.Abort("404")
	}
	channelModel, channelErr := biz.NewApiChannel().Find(chanelName, 0)
	if channelErr != nil {
		this.Ctx.WriteString(channelErr.Error())
		this.StopRun()
	}
	categoryModel, categoryErr := biz.NewApiArticle().CategoryOne(0, categoryName)
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
	this.Data["channel"] = channelModel
	this.Data["category"] = categoryModel
	this.Data["chanelName"] = chanelName
	this.Data["categoryName"] = categoryName

	if categoryModel.Template == "" {
		this.TplName = this.getView(DefatulSite.Template, "category.html")
	} else {
		if xstring.HasSuffix(categoryModel.Template, ".html", ".htm", ".tpl") == false {
			categoryModel.Template += ".html"
		}
		this.TplName = this.getView(DefatulSite.Template, categoryModel.Template)
	}
}
