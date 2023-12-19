package www

import (
	"github.com/beego/beego/v2/core/logs"

	lib "haedu.gov.cn/cms/app/tool"
)

type ApiTopicController struct{ BaseController }

// @router /api/topic/get [get]
func (ctrl *ApiTopicController) Get() {
	site_flag := ctrl.GetSafeString("site_flag")
	site_id, _ := ctrl.GetInt64("site_id")
	channel_id, _ := ctrl.GetInt64("channel_id")
	limit, _ := ctrl.GetInt("limit", 6)
	out, len, err := ctrl.BaseController.TopicGet(limit, site_id, site_flag, channel_id)
	if err != nil {
		logs.Error("Get::", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, len)
}

// @router /api/topic/get/new [get]
func (ctrl *ApiTopicController) GetNew() {
	site_flag := ctrl.GetSafeString("site_flag")
	site_id, _ := ctrl.GetInt64("site_id")
	channel_id, _ := ctrl.GetInt64("channel_id")
	limit, _ := ctrl.GetInt("limit", 6)
	out, len, err := ctrl.BaseController.TopicGetNew(limit, site_id, site_flag, channel_id)
	if err != nil {
		logs.Error("GetNew::", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, len)
}

/**
 * @description: Click 点击数+1
 * @param {int64} ads_id 广告ID
 * @return {*}
 */
// @router /api/topic/click [get]
func (ctrl *ApiTopicController) Click() {
	Topic_id, _ := ctrl.GetInt64("topic_id", 0)
	err := ctrl.BaseController.TopicClick(Topic_id)
	if err != nil {
		logs.Error("Click", err.Error())
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("", nil)
}

/**
 * @description: ArticlePaginate 获取文章分页列表
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {string} topic_name 标签名称
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
// @router /api/topic/article/paginate [get]
func (ctrl *ApiTopicController) ArticlePaginate() {
	limit, _ := ctrl.GetInt("limit", 6)
	page, _ := ctrl.GetInt("page", 1)
	site_id, _ := ctrl.GetInt64("site_id")
	site_flag := ctrl.GetSafeString("site_flag")
	order_by := ctrl.GetSafeString("order_by", "sort_id")
	call_index := ctrl.GetSafeString("call_index")
	channel_name := ctrl.GetSafeString("channel_name")
	keyword := ctrl.GetSafeString("keyword")
	topic_name := ctrl.GetSafeString("topic_name")
	channel_id, _ := ctrl.GetInt64("channel_id", -1)
	category_id, _ := ctrl.GetInt64("category_id", -1)
	is_top, _ := ctrl.GetInt("is_top", -1)
	is_red, _ := ctrl.GetInt("is_red", -1)
	is_hot, _ := ctrl.GetInt("is_hot", -1)
	is_slide, _ := ctrl.GetInt("is_slide", -1)
	is_search, _ := ctrl.GetInt("is_search", -1)
	outArticle, count, err := ctrl.BaseController.TopicArticlePaginate(page, limit, topic_name, site_id, site_flag, channel_id, channel_name, category_id, call_index, keyword, is_top, is_red, is_hot, is_slide, is_search, order_by)
	if err != nil {
		ctrl.JSONPage(lib.CodeError, err.Error(), outArticle, count)
	}
	ctrl.JSONPageSuccess(outArticle, count)
}
