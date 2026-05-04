package www

import (
	"github.com/beego/beego/v2/core/logs"

	service_model "haedu.gov.cn/cms/app/cms/service/model"
	lib "haedu.gov.cn/cms/app/tool"
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
func (ctrl *ApiArticleController) CategoryNav() {
	channel_name := ctrl.GetSafeString("channel_name")
	channel_id, _ := ctrl.GetInt64("channel_id")
	call_index := ctrl.GetSafeString("call_index")
	category_id, _ := ctrl.GetInt64("category_id")
	article_id, _ := ctrl.GetInt64("article_id")

	outNav, err := ctrl.BaseController.CategoryNav(channel_name, channel_id, call_index, category_id, article_id)
	if err != nil {
		ctrl.JSONPage(lib.CodeError, err.Error(), outNav, 0)
	}
	ctrl.JSONPageSuccess(outNav, int64(len(outNav)))
}

/**
 * @description: CategoryGet 获取栏目列表
 * @param {string} channel_name 频道名称
 * @param {string} call_index 栏目别名
 * @return {*}
 */
// @router /api/category/get [get]
func (ctrl *ApiArticleController) CategoryGet() {
	channel_name := ctrl.GetSafeString("channel_name")
	call_index := ctrl.GetSafeString("call_index")

	if channel_name == "" && call_index == "" {
		ctrl.JSONErrorOfData("频道编码或栏目编码不能为空", nil)
	}

	outChannel, count, err := ctrl.BaseController.CategoryGet(channel_name, call_index)
	if err != nil {
		ctrl.JSONPage(lib.CodeError, err.Error(), outChannel, count)
	}
	ctrl.JSONPageSuccess(outChannel, count)
}

/**
 * @description: CategoryFind 获取栏目详情
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @return {*}
 */
// @router /api/category/find [get]
func (ctrl *ApiArticleController) CategoryFind() {
	category_id, _ := ctrl.GetInt64("category_id")
	call_index := ctrl.GetSafeString("call_index")
	outChannel, err := ctrl.BaseController.CategoryFind(category_id, call_index)
	if err != nil {
		ctrl.JSONPage(lib.CodeError, err.Error(), outChannel, 0)
	}
	ctrl.JSONPageSuccess(outChannel, 1)
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
// @router /api/article/get [get]
func (ctrl *ApiArticleController) Get() {
	limit, _ := ctrl.GetInt("limit", 6)
	call_index := ctrl.GetSafeString("call_index")
	channel_name := ctrl.GetSafeString("channel_name")
	order_by := ctrl.GetSafeString("order_by", "")
	channel_id, _ := ctrl.GetInt64("channel_id", 0)
	category_id, _ := ctrl.GetInt64("category_id", 0)
	is_top, _ := ctrl.GetInt("is_top", 0)
	is_red, _ := ctrl.GetInt("is_red", 0)
	is_hot, _ := ctrl.GetInt("is_hot", 0)
	is_slide, _ := ctrl.GetInt("is_slide", 0)
	outArticle, count, err := ctrl.BaseController.ArticleGet(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
	if err != nil {
		ctrl.JSONPage(lib.CodeError, err.Error(), outArticle, count)
	}
	ctrl.JSONPageSuccess(outArticle, count)
}

/**
 * @description: GetNew 获取最新文章列表
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
// @router /api/article/get/new [get]
func (ctrl *ApiArticleController) GetNew() {
	limit, _ := ctrl.GetInt("limit", 6)
	call_index := ctrl.GetSafeString("call_index")
	channel_name := ctrl.GetSafeString("channel_name")
	order_by := ctrl.GetSafeString("order_by", "")
	channel_id, _ := ctrl.GetInt64("channel_id", 0)
	category_id, _ := ctrl.GetInt64("category_id", 0)
	is_top, _ := ctrl.GetInt("is_top", 0)
	is_red, _ := ctrl.GetInt("is_red", 0)
	is_hot, _ := ctrl.GetInt("is_hot", 0)
	is_slide, _ := ctrl.GetInt("is_slide", 0)
	outArticle, count, err := ctrl.BaseController.ArticleGetNew(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
	if err != nil {
		ctrl.JSONPage(lib.CodeError, err.Error(), outArticle, count)
	}
	ctrl.JSONPageSuccess(outArticle, count)
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
// @router /api/article/paginate [get]
func (ctrl *ApiArticleController) Paginate() {
	limit, _ := ctrl.GetInt("limit", 6)
	page, _ := ctrl.GetInt("page", 1)
	order_by := ctrl.GetSafeString("order_by", "")
	call_index := ctrl.GetSafeString("call_index")
	channel_name := ctrl.GetSafeString("channel_name")
	keyword := ctrl.GetSafeString("keyword")
	channel_id, _ := ctrl.GetInt64("channel_id", -1)
	category_id, _ := ctrl.GetInt64("category_id", -1)
	is_top, _ := ctrl.GetInt("is_top", -1)
	is_red, _ := ctrl.GetInt("is_red", -1)
	is_hot, _ := ctrl.GetInt("is_hot", -1)
	is_slide, _ := ctrl.GetInt("is_slide", -1)
	is_search, _ := ctrl.GetInt("is_search", -1)
	outArticle, count, err := ctrl.BaseController.ArticlePaginate(page, limit, channel_id, channel_name, category_id, call_index, keyword, is_top, is_red, is_hot, is_slide, is_search, order_by)
	if err != nil {
		ctrl.JSONPage(lib.CodeError, err.Error(), outArticle, count)
	}
	ctrl.JSONPageSuccess(outArticle, count)
}

/**
 * @description: 根据article_id获取文章详情、相册、附件
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/find [get]
func (ctrl *ApiArticleController) Find() {
	article_id, _ := ctrl.GetInt64("article_id", 0)
	call_index := ctrl.GetSafeString("call_index")
	aritcle, album, attatch, property, err := ctrl.BaseController.ArticleFind(call_index, article_id)
	result := service_model.ApiArticleModel{
		Article:  aritcle,
		Album:    album,
		Attach:   attatch,
		Property: property,
	}
	if err != nil {
		logs.Error("", err.Error())
		ctrl.JSONPage(lib.CodeError, err.Error(), result, 0)
	}
	ctrl.JSONPageSuccess(result, 1)
}

/**
 * @description: Get 根据article_id获取文章上一个、下一个
 * @param {string} call_index 栏目调用别名
 * @param {int64} category_id 栏目id
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/prev_next [get]
func (ctrl *ApiArticleController) PrevNext() {
	article_id, _ := ctrl.GetInt64("article_id", 0)
	category_id, _ := ctrl.GetInt64("category_id", 0)
	call_index := ctrl.GetSafeString("call_index")
	prev, next := ctrl.BaseController.ArticlePrevNext(call_index, category_id, article_id)
	result := struct {
		Prev *service_model.ApiArticlePrevNextModel `json:"prev"`
		Next *service_model.ApiArticlePrevNextModel `json:"next"`
	}{
		Prev: prev,
		Next: next,
	}
	var count int64
	if result.Prev.CategoryID > 0 {
		count++
	}
	if result.Next.CategoryID > 0 {
		count++
	}
	ctrl.JSONPageSuccess(result, count)
}

/**
 * @description: Article 获取文章详情
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/article [get]
func (ctrl *ApiArticleController) Article() {
	article_id, _ := ctrl.GetInt64("article_id", 0)
	call_index := ctrl.GetSafeString("call_index")
	article, err := ctrl.BaseController.ArticleArticle(call_index, article_id)
	if err != nil {
		logs.Error("", err.Error())
		ctrl.JSONPage(lib.CodeError, err.Error(), article, 0)
	}
	ctrl.JSONPageSuccess(article, 1)
}

/**
 * @description: Album 获取文章相册列表
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/album [get]
func (ctrl *ApiArticleController) Album() {
	article_id, _ := ctrl.GetInt64("article_id", 0)
	call_index := ctrl.GetSafeString("call_index")
	type_id, _ := ctrl.GetInt32("type_id", 0)
	album, err := ctrl.BaseController.ArticleAlbum(call_index, article_id, type_id)
	if err != nil {
		logs.Error("", err.Error())
		ctrl.JSONPage(lib.CodeError, err.Error(), album, 0)
	}
	ctrl.JSONPageSuccess(album, int64(len(album)))
}

/**
 * @description: Attach 获取文章附件列表
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/attach [get]
func (ctrl *ApiArticleController) Attach() {
	article_id, _ := ctrl.GetInt64("article_id", 0)
	call_index := ctrl.GetSafeString("call_index")
	type_id, _ := ctrl.GetInt32("type_id", 0)
	attach, err := ctrl.BaseController.ArticleAttach(call_index, article_id, type_id)
	if err != nil {
		logs.Error("", err.Error())
		ctrl.JSONPage(lib.CodeError, err.Error(), attach, 0)
	}
	ctrl.JSONPageSuccess(attach, int64(len(attach)))
}

/**
 * @description: Click 点击数+1
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/click [get]
func (ctrl *ApiArticleController) Click() {
	article_id, _ := ctrl.GetInt64("article_id", 0)
	call_index := ctrl.GetSafeString("call_index")
	err := ctrl.BaseController.ArticleClick(call_index, article_id)
	if err != nil {
		logs.Error("Click", err.Error())
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("", nil)
}

/**
 * @description: Like 点赞数+1
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
// @router /api/article/like [get]
func (ctrl *ApiArticleController) Like() {
	article_id, _ := ctrl.GetInt64("article_id", 0)
	call_index := ctrl.GetSafeString("call_index")
	err := ctrl.BaseController.ArticleLike(call_index, article_id)
	if err != nil {
		logs.Error("Like", err.Error())
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("", nil)
}

/**
 * @description: AlbumClick 点击数+1
 * @param {int64} article_id 文章id
 * @param {int64} ablum_id 图片id
 * @return {*}
 */
// @router /api/article/album/click [get]
func (ctrl *ApiArticleController) AlbumClick() {
	article_id, _ := ctrl.GetInt64("article_id", 0)
	ablum_id, _ := ctrl.GetInt64("ablum_id", 0)
	err := ctrl.BaseController.AlbumClick(article_id, ablum_id)
	if err != nil {
		logs.Error("AlbumClick", err.Error())
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("", nil)
}

/**
 * @description: Property 内容自定义属性
 * @param {int64} parentId 父id
 * @param {int64} articleId 文章id
 * @param {string} callIndex 属性别名
 * @param {string} title 属性名称
 * @return {*}
 */
// @router /api/article/property [get]
func (ctrl *ApiArticleController) Property() {
	page := 1
	limit := 99999
	parentId, _ := ctrl.GetInt64("parentId")
	articleId, _ := ctrl.GetInt64("articleId")
	callIndex := ctrl.GetSafeString("callIndex")
	title := ctrl.GetSafeString("title")
	property, count, err := ctrl.BaseController.Property(page, limit, parentId, articleId, callIndex, title)
	if err != nil {
		logs.Error("Property", err.Error())
		ctrl.JSONPage(lib.CodeError, err.Error(), property, 0)
	}
	ctrl.JSONPageSuccess(property, count)
}
