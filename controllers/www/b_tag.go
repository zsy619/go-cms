package www

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
 * @description: 获取标签列表
 * @param {int} limit 限制数量
 * @param {*} site_flag 站点标识
 * @param {*} site_id 站点ID
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func (this *BaseController) TagGet(limit int, site_id int64, site_flag string, channel_id int64) ([]*bizmodel.ApiTagModel, int64, error) {
	return biz.NewApiTag().Get(limit, site_id, site_flag, channel_id)
}

/**
 * @description: 获取最新标签列表
 * @param {int} limit 限制数量
 * @param {*} site_flag 站点标识
 * @param {*} site_id 站点ID
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func (this *BaseController) TagGetNew(limit int, site_id int64, site_flag string, channel_id int64) ([]*bizmodel.ApiTagModel, int64, error) {
	return biz.NewApiTag().GetNew(limit, site_id, site_flag, channel_id)
}

/**
 * @description: TagClick 点击数+1
 * @param {int64} tag_id 标签ID
 * @return {*}
 */
func (this *BaseController) TagClick(tag_id int64) error {
	return biz.NewApiTag().Click(tag_id)
}

/**
 * @description: TagArticlePaginate 获取文章分页列表
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
func (this *BaseController) TagArticlePaginate(page, limit int, tag_name string, channel_id int64, channel_name string, category_id int64, call_index string, keyword string, is_top, is_red, is_hot, is_slide, is_search int, order_by string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	return biz.NewApiTag().ArticlePaginate(page, limit, tag_name, channel_id, channel_name, category_id, call_index, keyword, is_top, is_red, is_hot, is_slide, is_search, order_by)
}
