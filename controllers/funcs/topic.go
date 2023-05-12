package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
 * @description: 获取最新专题列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func TopicNewExtend(limit int, site_id int64, site_flag string, channel_id int64) []*bizmodel.ApiTopicModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, _ := biz.NewApiTopic().GetNew(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*bizmodel.ApiTopicModel{}
	}
	return find
}

/**
 * @description: 获取最新专题列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @return {*}
 */
func TopicNew(limit int) []*bizmodel.ApiTopicModel {
	return TopicNewExtend(limit, 0, "", 0)
}

/**
 * @description: 获取专题列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func TopicTopExtend(limit int, site_id int64, site_flag string, channel_id int64) []*bizmodel.ApiTopicModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, _ := biz.NewApiTopic().Get(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*bizmodel.ApiTopicModel{}
	}
	return find
}

/**
 * @description: 获取专题列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @return {*}
 */
func TopicTop(limit int) []*bizmodel.ApiTopicModel {
	return TopicTopExtend(limit, 0, "", 0)
}

/**
 * @description: 获取专题文章列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {string} topic_name
 * @return {*}
 */
func TopicArtilceTop(limit int, topic_name string) []*bizmodel.ApiArticleListModel {
	list, _, err := biz.NewApiTopic().ArtilceTop(limit, topic_name)
	if err != nil {
		return []*bizmodel.ApiArticleListModel{}
	}
	return list
}
