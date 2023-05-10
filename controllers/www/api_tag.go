package www

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/lib"
)

type ApiTagController struct{ BaseController }

// @router /api/tag/get [get]
func (this *ApiTagController) Get() {
	site_flag := this.GetString("site_flag")
	site_id, _ := this.GetInt64("site_id")
	channel_id, _ := this.GetInt64("channel_id")
	limit, _ := this.GetInt("limit", 6)
	out, len, err := this.BaseController.TagGet(limit, site_id, site_flag, channel_id)
	if err != nil {
		logs.Error("Get::", err)
		this.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	this.JSONPageSuccess(out, len)
}

// @router /api/tag/get/new [get]
func (this *ApiTagController) GetNew() {
	site_flag := this.GetString("site_flag")
	site_id, _ := this.GetInt64("site_id")
	channel_id, _ := this.GetInt64("channel_id")
	limit, _ := this.GetInt("limit", 6)
	out, len, err := this.BaseController.TagGetNew(limit, site_id, site_flag, channel_id)
	if err != nil {
		logs.Error("GetNew::", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}

/**
 * @description: Click 点击数+1
 * @param {int64} ads_id 广告ID
 * @return {*}
 */
// @router /api/tag/click [get]
func (this *ApiTagController) Click() {
	tag_id, _ := this.GetInt64("tag_id", 0)
	err := this.BaseController.TagClick(tag_id)
	if err != nil {
		logs.Error("Click", err.Error())
		this.JSONError(err.Error())
	}
	this.JSONSuccess("", nil)
}

/**
 * @description: ArticlePaginate 获取文章分页列表
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {string} tag_name 标签名称
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
// @router /api/tag/article/paginate [get]
func (this *ApiTagController) ArticlePaginate() {
	limit, _ := this.GetInt("limit", 6)
	page, _ := this.GetInt("page", 1)
	order_by := this.GetString("order_by", "sort_id")
	call_index := this.GetString("call_index")
	channel_name := this.GetString("channel_name")
	keyword := this.GetString("keyword")
	tag_name := this.GetString("tag_name")
	channel_id, _ := this.GetInt64("channel_id", -1)
	category_id, _ := this.GetInt64("category_id", -1)
	is_top, _ := this.GetInt("is_top", -1)
	is_red, _ := this.GetInt("is_red", -1)
	is_hot, _ := this.GetInt("is_hot", -1)
	is_slide, _ := this.GetInt("is_slide", -1)
	is_search, _ := this.GetInt("is_search", -1)
	outArticle, count, err := this.BaseController.TagArticlePaginate(page, limit, tag_name, channel_id, channel_name, category_id, call_index, keyword, is_top, is_red, is_hot, is_slide, is_search, order_by)
	if err != nil {
		this.JSONPage(lib.CodeError, err.Error(), outArticle, count)
	}
	this.JSONPageSuccess(outArticle, count)
}
