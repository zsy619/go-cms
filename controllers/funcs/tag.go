package funcs

import (
	"haedu.gov.cn/cms/app/cms/service"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
)

/**
 * @description: 获取最新标签列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func TagNewExtend(limit int, site_id int64, site_flag string, channel_id int64) []*service_model.ApiTagModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, _ := service.NewApiTag().GetNew(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*service_model.ApiTagModel{}
	}
	return find
}

/**
 * @description: 获取最新标签列表
 * @param {int} limit
 * @return {*}
 */
func TagNew(limit int) []*service_model.ApiTagModel {
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
func TagTopExtend(limit int, site_id int64, site_flag string, channel_id int64) []*service_model.ApiTagModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, _ := service.NewApiTag().Get(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*service_model.ApiTagModel{}
	}
	return find
}

/**
 * @description: 获取标签列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @return {*}
 */
func TagTop(limit int) []*service_model.ApiTagModel {
	return TagTopExtend(limit, 0, "", 0)
}

/**
 * @description: 获取标签文章列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {string} tag_name
 * @return {*}
 */
func TagArtilceTop(limit int, tag_name string) []*service_model.ApiArticleListModel {
	if limit <= 0 {
		limit = 6
	}
	list, _, err := service.NewApiTag().ArtilceTop(limit, tag_name)
	if err != nil {
		return []*service_model.ApiArticleListModel{}
	}
	return list
}
