package www

import (
	"fmt"
	"strings"
	"time"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers"
)

type BaseController struct {
	controllers.BaseController
}

func (c *BaseController) Prepare() {
	fmt.Println("www BaseController Prepare")
	c.BaseController.Prepare()
	defaultSite, _ := c.SiteDefault()
	c.Data["site"] = defaultSite
	c.Data["time"] = time.Now().Unix()
	c.Data["webroot"] = "/static/www/"
	c.Data["year"] = time.Now().Year()
	c.Data["controllerName"] = strings.ToLower(c.ControllerName)
	c.Data["actionName"] = strings.ToLower(c.ActionName)
	debug := c.GetString("debug")
	c.Data["debug"] = debug
}

func (c *BaseController) Finish() {
	fmt.Println("www BaseController Finish")
}

// 渲染模版
func (this *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "www/" + this.ControllerName + "/" + this.ActionName + ".html"
	}
	this.Layout = "www/layout/layout.html"
	this.TplName = tplname
}

func (this *BaseController) displayArticle(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "www/" + this.ControllerName + "/" + this.ActionName + ".html"
	}
	this.Layout = "www/layout/article.html"
	this.TplName = tplname
}

func (this *BaseController) displayCategory(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "www/" + this.ControllerName + "/" + this.ActionName + ".html"
	}
	this.Layout = "www/layout/category.html"
	this.TplName = tplname
}

// 渲染模版
func (this *BaseController) displayNoLayout(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "www/" + this.ControllerName + "/" + this.ActionName + ".html"
	}
	this.TplName = tplname
}

/**
 * @description: SiteDefault 获取站点信息
 * @return {*}
 */
func (this *BaseController) SiteDefault() (*model.CmsSite, error) {
	return biz.NewApiSite().Default()
}

/**
 * @description: SiteGet 获取站点信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *BaseController) SiteGet(site_id int64) (*model.CmsSite, error) {
	return biz.NewApiSite().Get(site_id)
}

/**
 * @description: ChannelFind 获取站点频道
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *BaseController) ChannelFind(site_id int64) ([]*bmodel.ApiChannelFindModel, int64, error) {
	return biz.NewApiSite().ChannelFind(site_id)
}

/**
* @description: LinkFind 获取链接列表
* @param {int} limit 获取数量
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
func (this *BaseController) LinkFind(limit int, category_id int64, call_index string) ([]*bmodel.ApiLinkListModel, int64, error) {
	return biz.NewApiLink().Find(limit, category_id, call_index)
}

/**
 * @description: LinkPaginate 获取链接列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {int64} category_id 链接分类ID
 * @param {string} call_index 链接分类标识
 * @return {*}
 */
func (this *BaseController) LinkPaginate(page, limit int, category_id int64, call_index string) ([]*bmodel.ApiLinkListModel, int64, error) {
	return biz.NewApiLink().Paginate(page, limit, category_id, call_index)
}

/**
 * @description: Click 点击数+1
 * @param {int64} link_id 链接ID
 * @return {*}
 */
func (this *BaseController) LinkClick(link_id int64) error {
	return biz.NewApiLink().Click(link_id)
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
func (this *BaseController) CategoryNav(channel_name string, channel_id int64, call_index string, category_id int64, article_id int64) ([]*bmodel.ApiCategoryNav, error) {
	return biz.NewApiArticle().CategoryNav(channel_name, channel_id, call_index, category_id, article_id)
}

/**
 * @description: CategoryFind 获取栏目列表
 * @param {string} channel_name 频道名称
 * @return {*}
 */
func (this *BaseController) CategoryFind(channel_name string) ([]*bmodel.ApiCategoryFindModel, int64, error) {
	return biz.NewApiArticle().CategoryFind(channel_name)
}

/**
 * @description: CategoryOne 获取栏目详情
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @return {*}
 */
func (this *BaseController) CategoryOne(category_id int64, call_index string) (*bmodel.ApiCategoryOneModel, error) {
	return biz.NewApiArticle().CategoryOne(category_id, call_index)
}

/**
 * @description: Find 获取文章列表
 * @param {int} limit 获取数量
 * @param {int64} channel_id 频道ID
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @param {bool} is_cache 是否使用缓存
 * @return {*}
 */
func (this *BaseController) ArticleFind(limit int, channel_id, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) ([]*bmodel.ApiArticleListModel, int64, error) {
	return biz.NewApiArticle().Find(limit, channel_id, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
}

/**
 * @description: ArticlePaginate 获取文章分页列表
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {int64} channel_id 频道ID
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
func (this *BaseController) ArticlePaginate(page, limit int, channel_id, category_id int64, call_index string, keyword string, is_top, is_red, is_hot, is_slide, is_search int, order_by string) ([]*bmodel.ApiArticleListModel, int64, error) {
	return biz.NewApiArticle().Paginate(page, limit, channel_id, category_id, call_index, keyword, is_top, is_red, is_hot, is_slide, is_search, order_by)
}

/**
 * @description: ArticleGet 根据article_id获取文章详情、相册、附件
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *BaseController) ArticleGet(call_index string, article_id int64) (*bmodel.ApiArticleOneModel, []*model.CmsArticleAlbum, []*model.CmsArticleAttach, error) {
	return biz.NewApiArticle().Get(call_index, article_id)
}

/**
 * @description: Article 获取文章详情
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *BaseController) ArticleArticle(call_index string, article_id int64) (*bmodel.ApiArticleOneModel, error) {
	return biz.NewApiArticle().Article(call_index, article_id)
}

/**
 * @description: ArticleAlbum 获取文章相册列表
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *BaseController) ArticleAlbum(call_index string, article_id int64) ([]*model.CmsArticleAlbum, error) {
	return biz.NewApiArticle().Album(call_index, article_id)
}

/**
 * @description: ArticleAttach 获取文章附件列表
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *BaseController) ArticleAttach(call_index string, article_id int64) ([]*model.CmsArticleAttach, error) {
	return biz.NewApiArticle().Attach(call_index, article_id)
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

type CategoryBaseController struct{ BaseController }

func (this *CategoryBaseController) Prepare() {
	this.BaseController.Prepare()
	category, _ := this.CategoryOne(0, this.ActionName)
	this.Data["category"] = category
}
