package www

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
)

type ApiArticleController struct{ BaseController }

/**
 * @description: CategoryNav 获取栏目导航
 * @param {string} channel_name 频道名称
 * @param {in64} channel_id 频道ID
 * @param {string} call_index 栏目别名
 * @param {in64} category_id 栏目ID
 * @param {int64} article_id 文章ID
 * @return {*}
 */
// @router /api/category/nav [get]
func (this *ApiArticleController) CategoryNav() {
	channel_name := this.GetString("channel_name")
	channel_id, _ := this.GetInt64("channel_id")
	call_index := this.GetString("call_index")
	category_id, _ := this.GetInt64("category_id")
	article_id, _ := this.GetInt64("article_id")

	outNav, err := this.BaseController.CategoryNav(channel_name, channel_id, call_index, category_id, article_id)
	if err != nil {
		this.JSONErrorOfData(err.Error(), outNav)
	}
	this.JSONSuccess("", outNav)
}

/**
 * @description: CategoryFind 获取栏目列表
 * @param {string} channel_name 频道名称
 * @return {*}
 */
// @router /api/category/find [get]
func (this *ApiArticleController) CategoryFind() {
	channel_name := this.GetString("channel_name")
	if channel_name == "" {
		this.JSONErrorOfData("频道名称不能为空", nil)
	}
	outChannel, count, err := this.BaseController.CategoryFind(channel_name)
	if err != nil {
		this.JSONErrorOfData(err.Error(), outChannel)
	}
	this.JSONSuccess(strconv.FormatInt(count, 10), outChannel)
}

/**
 * @description: CategoryOne 获取栏目详情
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @return {*}
 */
// @router /api/category/one [get]
func (this *ApiArticleController) CategoryOne() {
	category_id, _ := this.GetInt64("category_id")
	call_index := this.GetString("call_index")
	outChannel, err := this.BaseController.CategoryOne(category_id, call_index)
	if err != nil {
		this.JSONErrorOfData(err.Error(), outChannel)
	}
	this.JSONSuccess("1", outChannel)
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
// @router /api/article/find [get]
func (this *ApiArticleController) Find() {
	limit, _ := this.GetInt("limit", 6)
	call_index := this.GetString("call_index")
	order_by := this.GetString("order_by", "sort_id")
	channel_id, _ := this.GetInt64("channel_id", 0)
	category_id, _ := this.GetInt64("category_id", 0)
	is_top, _ := this.GetInt("is_top", 0)
	is_red, _ := this.GetInt("is_red", 0)
	is_hot, _ := this.GetInt("is_hot", 0)
	is_slide, _ := this.GetInt("is_slide", 0)
	outArticle, count, err := this.BaseController.ArticleFind(limit, channel_id, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
	if err != nil {
		this.JSONErrorOfData(err.Error(), outArticle)
	}
	this.JSONSuccess(strconv.FormatInt(count, 10), outArticle)
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
// @router /api/article/paginate [get]
func (this *ApiArticleController) Paginate() {
	limit, _ := this.GetInt("limit", 6)
	page, _ := this.GetInt("page", 1)
	order_by := this.GetString("order_by", "sort_id")
	call_index := this.GetString("call_index")
	keyword := this.GetString("keyword")
	channel_id, _ := this.GetInt64("channel_id", -1)
	category_id, _ := this.GetInt64("category_id", -1)
	is_top, _ := this.GetInt("is_top", -1)
	is_red, _ := this.GetInt("is_red", -1)
	is_hot, _ := this.GetInt("is_hot", -1)
	is_slide, _ := this.GetInt("is_slide", -1)
	is_search, _ := this.GetInt("is_search", -1)
	outArticle, count, err := this.BaseController.ArticlePaginate(page, limit, channel_id, category_id, call_index, keyword, is_top, is_red, is_hot, is_slide, is_search, order_by)
	if err != nil {
		this.JSONPage(lib.CodeError, err.Error(), outArticle, count)
	}
	this.JSONPageSuccess(outArticle, count)
}

/**
 * @description: One 根据article_id获取文章详情、相册、附件
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/one [get]
func (this *ApiArticleController) One() {
	article_id, _ := this.GetInt64("article_id", 0)
	call_index := this.GetString("call_index")
	aritcle, album, attatch, err := this.BaseController.ArticleGet(call_index, article_id)
	result := struct {
		Article *bizmodel.ApiArticleOneModel `json:"article"`
		Album   []*model.CmsAlbum            `json:"album"`
		Attach  []*model.CmsAttach           `json:"attatch"`
	}{
		Article: aritcle,
		Album:   album,
		Attach:  attatch,
	}
	if err != nil {
		logs.Error("", err.Error())
		this.JSONErrorOfData(err.Error(), result)
	}
	this.JSONSuccess("", result)
}

/**
 * @description: Get 根据article_id获取文章上一个、下一个
 * @param {string} call_index 栏目调用别名
 * @param {int64} category_id 栏目id
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/prev_next [get]
func (this *ApiArticleController) PrevNext() {
	article_id, _ := this.GetInt64("article_id", 0)
	category_id, _ := this.GetInt64("category_id", 0)
	call_index := this.GetString("call_index")
	prev, next := this.BaseController.ArticlePrevNext(call_index, category_id, article_id)
	result := struct {
		Prev *bizmodel.ApiArticleOneModel `json:"prev"`
		Next *bizmodel.ApiArticleOneModel `json:"next"`
	}{
		Prev: prev,
		Next: next,
	}
	this.JSONSuccess("", result)
}

/**
 * @description: Article 获取文章详情
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/article [get]
func (this *ApiArticleController) Article() {
	article_id, _ := this.GetInt64("article_id", 0)
	call_index := this.GetString("call_index")
	article, err := this.BaseController.ArticleArticle(call_index, article_id)
	if err != nil {
		logs.Error("", err.Error())
		this.JSONErrorOfData(err.Error(), article)
	}
	this.JSONSuccess("", article)
}

/**
 * @description: Album 获取文章相册列表
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/album [get]
func (this *ApiArticleController) Album() {
	article_id, _ := this.GetInt64("article_id", 0)
	call_index := this.GetString("call_index")
	type_id, _ := this.GetInt32("type_id", 0)
	album, err := this.BaseController.ArticleAlbum(call_index, article_id, type_id)
	if err != nil {
		logs.Error("", err.Error())
		this.JSONErrorOfData(err.Error(), album)
	}
	this.JSONSuccess("", album)
}

/**
 * @description: Attach 获取文章附件列表
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/attach [get]
func (this *ApiArticleController) Attach() {
	article_id, _ := this.GetInt64("article_id", 0)
	call_index := this.GetString("call_index")
	type_id, _ := this.GetInt32("type_id", 0)
	attach, err := this.BaseController.ArticleAttach(call_index, article_id, type_id)
	if err != nil {
		logs.Error("", err.Error())
		this.JSONErrorOfData(err.Error(), attach)
	}
	this.JSONSuccess("", attach)
}

/**
 * @description: Click 点击数+1
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/click [get]
func (this *ApiArticleController) Click() {
	article_id, _ := this.GetInt64("article_id", 0)
	call_index := this.GetString("call_index")
	err := this.BaseController.ArticleClick(call_index, article_id)
	if err != nil {
		logs.Error("Click", err.Error())
		this.JSONError(err.Error())
	}
	this.JSONSuccess("", nil)
}

/**
 * @description: Like 点赞数+1
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/like [get]
func (this *ApiArticleController) Like() {
	article_id, _ := this.GetInt64("article_id", 0)
	call_index := this.GetString("call_index")
	err := this.BaseController.ArticleLike(call_index, article_id)
	if err != nil {
		logs.Error("Like", err.Error())
		this.JSONError(err.Error())
	}
	this.JSONSuccess("", nil)
}

/**
 * @description: AlbumClick 点击数+1
 * @param {int64} article_id 文章id
 * @param {int64} ablum_id 图片id
 * @return {*}
 */
// @router /api/article/album/click [get]
func (this *ApiArticleController) AlbumClick() {
	article_id, _ := this.GetInt64("article_id", 0)
	ablum_id, _ := this.GetInt64("ablum_id", 0)
	err := this.BaseController.AlbumClick(article_id, ablum_id)
	if err != nil {
		logs.Error("AlbumClick", err.Error())
		this.JSONError(err.Error())
	}
	this.JSONSuccess("", nil)
}
