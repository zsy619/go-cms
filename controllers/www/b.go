package www

import (
	"fmt"
	"path"
	"strings"
	"time"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers"
)

var (
	DefaultSite *bizmodel.ApiSiteModel // 默认站点
	SiteStatic  string                 // 站点静态文件
	SiteTheme   string                 // 站点模板
)

type BaseController struct {
	controllers.BaseController
}

/**
 * @description: 获取视图地址
 * @param {*} themeName 皮肤名称
 * @param {string} viewName view名称
 * @return {*}
 */
func (c *BaseController) GetView(themeName, viewName string) string {
	return "themes/" + themeName + "/views/" + viewName
}

/**
 * @description: Prepare
 * @return {*}
 */
func (c *BaseController) Prepare() {
	c.BaseController.Prepare()
	fmt.Println("www BaseController Prepare")

	// 根据域名获取站点信息
	DefaultSite, _ = c.SiteByHost(c.Ctx.Request.Host)
	if DefaultSite == nil {
		// 获取默认站点、模板
		DefaultSite, _ = c.SiteDefault()
		if DefaultSite.Template == "" {
			c.Ctx.WriteString("请设置默认模板")
			c.StopRun()
		}
	}
	SiteTheme = DefaultSite.Template
	SiteStatic = "/views/themes/" + SiteTheme + "/static/"
	c.Data["siteTheme"] = SiteTheme
	c.Data["siteStatic"] = SiteStatic
	c.Data["siteImages"] = path.Join(SiteStatic, "images")
	c.Data["siteJs"] = path.Join(SiteStatic, "js")
	c.Data["siteCss"] = path.Join(SiteStatic, "css")
	c.Data["siteViews"] = "themes/" + SiteTheme + "/views/"

	c.Data["site"] = DefaultSite
	channel, _, _ := c.ChannelGet(DefaultSite.SiteID)
	c.Data["channel"] = channel

	c.Data["webroot"] = "/static/www/"
	c.Data["year"] = time.Now().Year()
	c.Data["controllerName"] = strings.ToLower(c.ControllerName)
	c.Data["actionName"] = strings.ToLower(c.ActionName)
	debug := c.GetString("debug")
	c.Data["debug"] = debug
}

/**
 * @description: 加载完毕
 * @return {*}
 */
func (c *BaseController) Finish() {
	fmt.Println("www BaseController Finish")
}

/**
 * @description: SiteDefault 获取默认站点信息
 * @return {*}
 */
func (this *BaseController) SiteDefault() (*bizmodel.ApiSiteModel, error) {
	return biz.NewApiSite().Default()
}

/**
 * @description: SiteByHost 根据域名获取站点信息
 * @return {*}
 */
func (this *BaseController) SiteByHost(host string) (*bizmodel.ApiSiteModel, error) {
	return biz.NewApiSite().FindByHost(host)
}

/**
 * @description: SiteFind 获取站点信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *BaseController) SiteFind(site_id int64) (*model.CmsSite, error) {
	return biz.NewApiSite().Find(site_id)
}

/**
 * @description: ChannelGet 获取站点频道
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *BaseController) ChannelGet(site_id int64) ([]*bizmodel.ApiChannelModel, int64, error) {
	return biz.NewApiSite().ChannelGet(site_id)
}

/**
 * @description: SiteMenu 获取站点导航
 * @param {int64} site_id 站点ID
 * @param {int64} channel_id 频道ID
 * @return {*}
 */
func (this *BaseController) SiteMenu(site_id, channel_id int64) ([]*bizmodel.ApiNavModel, int64, error) {
	return biz.NewApiSite().NavGet(site_id, channel_id)
}

/**
 * @description: SiteMenuFlag 获取站点导航
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 频道ID
 * @return {*}
 */
func (this *BaseController) SiteMenuFlag(site_flag string, channel_id int64) ([]*bizmodel.ApiNavModel, int64, error) {
	return biz.NewApiSite().NavGetByFlag(site_flag, channel_id)
}

/**
 * @description: CategoryNav 获取栏目导航
 * @param {string} channel_name 频道名称
 * @param {in64} channel_id 频道ID
 * @param {string} call_index 栏目别名
 * @param {in64} category_id 栏目ID
 * @param {int64} article_id 文章ID
 * @return {*}
 */
func (this *BaseController) CategoryNav(channel_name string, channel_id int64, call_index string, category_id int64, article_id int64) ([]*bizmodel.ApiCategoryNav, error) {
	return biz.NewApiArticle().CategoryNav(channel_name, channel_id, call_index, category_id, article_id)
}

/**
 * @description: CategoryGet 获取栏目列表
 * @param {string} channel_name 频道名称
 * @param {string} call_index 栏目别名
 * @return {*}
 */
func (this *BaseController) CategoryGet(channel_name, call_index string) ([]*bizmodel.ApiCategoryGetModel, int64, error) {
	return biz.NewApiArticle().CategoryGet(channel_name, call_index)
}

/**
 * @description: CategoryFind 获取栏目详情
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @return {*}
 */
func (this *BaseController) CategoryFind(category_id int64, call_index string) (*bizmodel.ApiCategoryFindModel, error) {
	return biz.NewApiArticle().CategoryFind(category_id, call_index)
}

/**
 * @description: ArticleGet 获取文章列表
 * @param {int} limit 获取数量
 * @param {int64} channel_id 频道ID
 * @param {string} channel_name 频道名称
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @return {*}
 */
func (this *BaseController) ArticleGet(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	return biz.NewApiArticle().ArticleGet(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
}

/**
 * @description: ArticleGetNew 获取最新文章列表
 * @param {int} limit 获取数量
 * @param {int64} channel_id 频道ID
 * @param {string} channel_name 频道名称
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @return {*}
 */
func (this *BaseController) ArticleGetNew(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	return biz.NewApiArticle().ArticleGetNew(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
}

/**
 * @description: ArticlePaginate 获取文章分页列表
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {int64} channel_id 频道ID
 * @param {string} channel_name 频道名称
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {string} keyword 关键词：按标题、摘要进行搜索
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {int} is_search 是否搜索
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @return {*}
 */
func (this *BaseController) ArticlePaginate(page, limit int, channel_id int64, channel_name string, category_id int64, call_index string, keyword string, is_top, is_red, is_hot, is_slide, is_search int, order_by string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	return biz.NewApiArticle().ArticlePaginate(page, limit, channel_id, channel_name, category_id, call_index, keyword, is_top, is_red, is_hot, is_slide, is_search, order_by)
}

/**
 * @description: ArticleFind 根据article_id获取文章详情、相册、附件
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *BaseController) ArticleFind(call_index string, article_id int64) (*bizmodel.ApiArticleOneModel, []*bizmodel.ApiAlbumModel, []*bizmodel.ApiAttachModel, error) {
	return biz.NewApiArticle().ArticleFind(call_index, article_id)
}

/**
 * @description: Get 根据article_id获取文章上一个、下一个
 * @param {string} call_index 栏目调用别名
 * @param {int64} category_id 栏目id
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *BaseController) ArticlePrevNext(call_index string, category_id, article_id int64) (*bizmodel.ApiArticlePrevNextModel, *bizmodel.ApiArticlePrevNextModel) {
	return biz.NewApiArticle().PrevNext(call_index, category_id, article_id)
}

/**
 * @description: Article 获取文章详情
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *BaseController) ArticleArticle(call_index string, article_id int64) (*bizmodel.ApiArticleOneModel, error) {
	return biz.NewApiArticle().Article(call_index, article_id)
}

/**
 * @description: Album 获取文章相册列表
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @param {int32} type_id 分类
 * @return {*}
 */
func (this *BaseController) ArticleAlbum(call_index string, article_id int64, type_id int32) ([]*bizmodel.ApiAlbumModel, error) {
	return biz.NewApiArticle().Album(call_index, article_id, type_id)
}

/**
 * @description: Attach 获取文章附件列表
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @param {int32} type_id 分类
 * @return {*}
 */
func (this *BaseController) ArticleAttach(call_index string, article_id int64, type_id int32) ([]*bizmodel.ApiAttachModel, error) {
	return biz.NewApiArticle().Attach(call_index, article_id, type_id)
}

/**
 * @description: ArticleClick 点击数+1
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *BaseController) ArticleClick(call_index string, article_id int64) error {
	return biz.NewApiArticle().Click(call_index, article_id)
}

/**
 * @description: Like 点赞数+1
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *BaseController) ArticleLike(call_index string, article_id int64) error {
	return biz.NewApiArticle().Like(call_index, article_id)
}

/**
 * @description: AlbumClick 点击数+1
 * @param {int64} article_id 文章id
 * @param {int64} ablum_id 图片id
 * @return {*}
 */
func (this *BaseController) AlbumClick(article_id, ablum_id int64) error {
	return biz.NewApiArticle().AlbumClick(article_id, ablum_id)
}
