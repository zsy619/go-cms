package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
 * @description: 获取最新标签列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func TagNewExtend(limit int, site_id int64, site_flag string, channel_id int64) []*bizmodel.ApiTagModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, _ := biz.NewApiTag().GetNew(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*bizmodel.ApiTagModel{}
	}
	return find
}

/**
 * @description: 获取最新标签列表
 * @param {int} limit
 * @return {*}
 */
func TagNew(limit int) []*bizmodel.ApiTagModel {
	return TagNewExtend(limit, 0, "", 0)
}

/**
 * @description: 获取标签列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func TagTopExtend(limit int, site_id int64, site_flag string, channel_id int64) []*bizmodel.ApiTagModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, _ := biz.NewApiTag().Get(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*bizmodel.ApiTagModel{}
	}
	return find
}

/**
 * @description: 获取标签列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @return {*}
 */
func TagTop(limit int) []*bizmodel.ApiTagModel {
	return TagTopExtend(limit, 0, "", 0)
}

/**
 * @description: 获取标签文章列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {string} tag_name
 * @return {*}
 */
func TagArtilceTop(limit int, tag_name string) []*bizmodel.ApiArticleListModel {
	if limit <= 0 {
		limit = 6
	}
	list, _, err := biz.NewApiTag().ArtilceTop(limit, tag_name)
	if err != nil {
		return []*bizmodel.ApiArticleListModel{}
	}
	return list
}
