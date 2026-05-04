package www

import (
	"github.com/zsy619/tools/xstring"

	"haedu.gov.cn/cms/app/cms/service"
)

// ChannelController 频道控制器
type ChannelController struct{ BaseController }

// Channel 频道首页
// @router /:flag/:name:string [get]
func (ctrl *ChannelController) Channel() {
	siteFlag := ctrl.Ctx.Input.Param(":flag")
	channelName := ctrl.Ctx.Input.Param(":name")
	if siteFlag == "" || channelName == "" {
		ctrl.Abort("404")
	}
	channelModel, err := service.NewApiChannel().Find(channelName, 0)
	if err != nil {
		ctrl.Ctx.WriteString(err.Error())
		ctrl.StopRun()
	}
	ctrl.Data["channel"] = channelModel
	ctrl.Data["channelName"] = channelName
	if channelModel.TmplChnl == "" {
		ctrl.TplName = ctrl.GetView(DefaultSite.Template, "channel.html")
	} else {
		if !xstring.HasSuffix(channelModel.TmplChnl, ".html", ".htm", ".tpl") {
			channelModel.TmplChnl += ".html"
		}
		ctrl.TplName = ctrl.GetView(DefaultSite.Template, channelModel.TmplChnl)
	}
}

// Category 频道分类
// @router /:flag/:name:string/:category:string [get]
func (ctrl *ChannelController) Category() {
	siteFlag := ctrl.Ctx.Input.Param(":flag")
	channelName := ctrl.Ctx.Input.Param(":name")
	categoryName := ctrl.Ctx.Input.Param(":category")
	if siteFlag == "" || channelName == "" || categoryName == "" {
		ctrl.Abort("404")
	}
	channelModel, channelErr := service.NewApiChannel().Find(channelName, 0)
	if channelErr != nil {
		ctrl.Ctx.WriteString(channelErr.Error())
		ctrl.StopRun()
	}
	categoryModel, categoryErr := service.NewApiArticle().CategoryFind(0, categoryName)
	if categoryErr != nil {
		ctrl.Ctx.WriteString(categoryErr.Error())
		ctrl.StopRun()
	}
	if categoryModel.ChannelID <= 0 {
		ctrl.Ctx.WriteString("分类不存在")
		ctrl.StopRun()
	}
	if categoryModel.ChannelID != channelModel.ChannelID {
		ctrl.Ctx.WriteString("分类与频道不匹配")
		ctrl.StopRun()
	}
	ctrl.Data["channel"] = channelModel
	ctrl.Data["category"] = categoryModel
	ctrl.Data["channelName"] = channelName
	ctrl.Data["categoryName"] = categoryName

	if categoryModel.TmplCat == "" {
		ctrl.TplName = ctrl.GetView(DefaultSite.Template, "category.html")
	} else {
		if !xstring.HasSuffix(categoryModel.TmplCat, ".html", ".htm", ".tpl") {
			categoryModel.TmplCat += ".html"
		}
		ctrl.TplName = ctrl.GetView(DefaultSite.Template, categoryModel.TmplCat)
	}
}
