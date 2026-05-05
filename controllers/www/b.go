package www

import (
	"path"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
	"haedu.gov.cn/cms/controllers"
)

// 全局站点信息
var (
	DefaultSite *service_model.ApiSiteModel // 默认站点
	SiteStatic  string                      // 站点静态文件路径
	SiteTheme   string                      // 站点模板名称
)

// BaseController 前台展示基础控制器
// 继承自controllers.BaseController,提供前台页面的公共功能
type BaseController struct{ controllers.BaseController }

// GetView 获取完整视图路径
// @param themeName 皮肤/模板名称
// @param viewName 视图名称
// @return string 完整视图路径
func (ctrl *BaseController) GetView(themeName, viewName string) string {
	return "themes/" + themeName + "/views/" + viewName
}

// Prepare 前台控制器前置处理
// 功能: 根据域名加载站点信息,初始化模板变量
func (ctrl *BaseController) Prepare() {
	ctrl.BaseController.Prepare()

	// 根据域名获取站点信息
	DefaultSite, _ = ctrl.SiteByHost(ctrl.Ctx.Request.Host)
	if DefaultSite == nil {
		DefaultSite, _ = ctrl.SiteDefault()
		if DefaultSite.Template == "" {
			ctrl.Ctx.WriteString("请设置默认模板")
			ctrl.StopRun()
		}
	}

	logs.Debug("前台站点初始化: host=%s, siteId=%d, template=%s",
		ctrl.Ctx.Request.Host, DefaultSite.SiteID, DefaultSite.Template)

	SiteTheme = DefaultSite.Template
	SiteStatic = "/views/themes/" + SiteTheme + "/static/"

	// 设置模板变量
	ctrl.Data["siteTheme"] = SiteTheme
	ctrl.Data["siteStatic"] = SiteStatic
	ctrl.Data["siteImages"] = path.Join(SiteStatic, "images")
	ctrl.Data["siteJs"] = path.Join(SiteStatic, "js")
	ctrl.Data["siteCss"] = path.Join(SiteStatic, "css")
	ctrl.Data["siteViews"] = "themes/" + SiteTheme + "/views/"
	ctrl.Data["site"] = DefaultSite

	channel, _, _ := ctrl.ChannelGet(DefaultSite.SiteID)
	ctrl.Data["channel"] = channel

	ctrl.Data["webroot"] = "/static/www/"
	ctrl.Data["year"] = time.Now().Year()
	ctrl.Data["controllerName"] = strings.ToLower(ctrl.ControllerName)
	ctrl.Data["actionName"] = strings.ToLower(ctrl.ActionName)
	ctrl.Data["debug"] = ctrl.GetSafeString("debug")
}

// Finish 前台控制器后置处理
func (ctrl *BaseController) Finish() {
	// 预留后置处理逻辑
}

// SiteDefault 获取默认站点信息
// @return *service_model.ApiSiteModel, error
func (ctrl *BaseController) SiteDefault() (*service_model.ApiSiteModel, error) {
	return service.NewApiSite().Default()
}

// SiteByHost 根据域名获取站点信息
// @param host 域名
// @return *service_model.ApiSiteModel, error
func (ctrl *BaseController) SiteByHost(host string) (*service_model.ApiSiteModel, error) {
	return service.NewApiSite().FindByHost(host)
}

// SiteFind 获取站点信息
// @param site_id 站点ID
// @return *domain.CmsSite, error
func (ctrl *BaseController) SiteFind(site_id int64) (*domain.CmsSite, error) {
	return service.NewApiSite().Find(site_id)
}

// ChannelGet 获取站点频道列表
// @param site_id 站点ID
// @return []*service_model.ApiChannelModel, int64, error
func (ctrl *BaseController) ChannelGet(site_id int64) ([]*service_model.ApiChannelModel, int64, error) {
	return service.NewApiSite().ChannelGet(site_id)
}

// SiteMenu 获取站点导航
// @param site_id 站点ID
// @param channel_id 频道ID
// @return []*service_model.ApiNavModel, int64, error
func (ctrl *BaseController) SiteMenu(site_id, channel_id int64) ([]*service_model.ApiNavModel, int64, error) {
	return service.NewApiSite().NavGet(site_id, channel_id)
}

// SiteMenuFlag 根据站点标识获取导航
// @param site_flag 站点标识
// @param channel_id 频道ID
// @return []*service_model.ApiNavModel, int64, error
func (ctrl *BaseController) SiteMenuFlag(site_flag string, channel_id int64) ([]*service_model.ApiNavModel, int64, error) {
	return service.NewApiSite().NavGetByFlag(site_flag, channel_id)
}

// CategoryNav 获取栏目导航面包屑
// @param channel_name 频道名称
// @param channel_id 频道ID
// @param call_index 栏目别名
// @param category_id 栏目ID
// @param article_id 文章ID
// @return []*service_model.ApiCategoryNav, error
func (ctrl *BaseController) CategoryNav(channel_name string, channel_id int64, call_index string, category_id int64, article_id int64) ([]*service_model.ApiCategoryNav, error) {
	return service.NewApiArticle().CategoryNav(channel_name, channel_id, call_index, category_id, article_id)
}

// CategoryGet 获取栏目列表
// @param channel_name 频道名称
// @param call_index 栏目别名
// @return []*service_model.ApiCategoryGetModel, int64, error
func (ctrl *BaseController) CategoryGet(channel_name, call_index string) ([]*service_model.ApiCategoryGetModel, int64, error) {
	return service.NewApiArticle().CategoryGet(channel_name, call_index)
}

// CategoryFind 获取栏目详情
// @param category_id 栏目ID
// @param call_index 栏目别名
// @return *service_model.ApiCategoryFindModel, error
func (ctrl *BaseController) CategoryFind(category_id int64, call_index string) (*service_model.ApiCategoryFindModel, error) {
	return service.NewApiArticle().CategoryFind(category_id, call_index)
}

// ArticleGet 获取文章列表
// @param limit 获取数量
// @param channel_id 频道ID
// @param channel_name 频道名称
// @param category_id 栏目ID
// @param call_index 栏目别名
// @param is_top 是否置顶
// @param is_red 是否推荐
// @param is_hot 是否热门
// @param is_slide 是否幻灯片
// @param order_by 排序字段
// @return []*service_model.ApiArticleListModel, int64, error
func (ctrl *BaseController) ArticleGet(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) ([]*service_model.ApiArticleListModel, int64, error) {
	return service.NewApiArticle().ArticleGet(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
}

// ArticleGetNew 获取最新文章列表
func (ctrl *BaseController) ArticleGetNew(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) ([]*service_model.ApiArticleListModel, int64, error) {
	return service.NewApiArticle().ArticleGetNew(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
}

// ArticlePaginate 获取文章分页列表
func (ctrl *BaseController) ArticlePaginate(page, limit int, channel_id int64, channel_name string, category_id int64, call_index string, keyword string, is_top, is_red, is_hot, is_slide, is_search int, order_by string) ([]*service_model.ApiArticleListModel, int64, error) {
	return service.NewApiArticle().ArticlePaginate(page, limit, channel_id, channel_name, category_id, call_index, keyword, is_top, is_red, is_hot, is_slide, is_search, order_by)
}

// ArticleFind 根据article_id获取文章详情、相册、附件
func (ctrl *BaseController) ArticleFind(call_index string, article_id int64) (*service_model.ApiArticleOneModel, []*service_model.ApiAlbumModel, []*service_model.ApiAttachModel, []*service_model.ApiPropertyModel, error) {
	return service.NewApiArticle().ArticleFind(call_index, article_id)
}

// ArticlePrevNext 获取文章上一篇和下一篇
func (ctrl *BaseController) ArticlePrevNext(call_index string, category_id, article_id int64) (*service_model.ApiArticlePrevNextModel, *service_model.ApiArticlePrevNextModel) {
	return service.NewApiArticle().PrevNext(call_index, category_id, article_id)
}

// ArticleArticle 获取文章详情
func (ctrl *BaseController) ArticleArticle(call_index string, article_id int64) (*service_model.ApiArticleOneModel, error) {
	return service.NewApiArticle().Article(call_index, article_id)
}

// ArticleAlbum 获取文章相册列表
func (ctrl *BaseController) ArticleAlbum(call_index string, article_id int64, type_id int32) ([]*service_model.ApiAlbumModel, error) {
	return service.NewApiArticle().Album(call_index, article_id, type_id)
}

// ArticleAttach 获取文章附件列表
func (ctrl *BaseController) ArticleAttach(call_index string, article_id int64, type_id int32) ([]*service_model.ApiAttachModel, error) {
	return service.NewApiArticle().Attach(call_index, article_id, type_id)
}

// ArticleClick 文章点击数+1
func (ctrl *BaseController) ArticleClick(call_index string, article_id int64) error {
	return service.NewApiArticle().Click(call_index, article_id)
}

// ArticleLike 文章点赞数+1
func (ctrl *BaseController) ArticleLike(call_index string, article_id int64) error {
	return service.NewApiArticle().Like(call_index, article_id)
}

// AlbumClick 相册图片点击数+1
func (ctrl *BaseController) AlbumClick(article_id, ablum_id int64) error {
	return service.NewApiArticle().AlbumClick(article_id, ablum_id)
}

// Property 获取内容自定义属性
func (ctrl *BaseController) Property(page, limit int, parentId, articleId int64, callIndex, title string) ([]*domain.CmsArticleProperty, int64, error) {
	return service.NewApiArticle().Property(page, limit, parentId, articleId, callIndex, title)
}