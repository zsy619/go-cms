package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
 * @description: 获取专题列表
 * @param {int} limit 限制数量
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func TopicNewExt(limit int, site_id int64, site_flag string, channel_id int64) []*bizmodel.ApiTopicModel {
	find, _, _ := biz.NewApiTopic().GetNew(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*bizmodel.ApiTopicModel{}
	}
	return find
}

func TopicNew(limit int) []*bizmodel.ApiTopicModel {
	return TopicNewExt(limit, 0, "", 0)
}

/**
 * @description: 获取最新专题列表
 * @param {int} limit 限制数量
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func TopicTopExt(limit int, site_id int64, site_flag string, channel_id int64) []*bizmodel.ApiTopicModel {
	find, _, _ := biz.NewApiTopic().Get(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*bizmodel.ApiTopicModel{}
	}
	return find
}

func TopicTop(limit int) []*bizmodel.ApiTopicModel {
	return TopicTopExt(limit, 0, "", 0)
}

func TopicArtilceTop(limit int, topic_name string) []*bizmodel.ApiArticleListModel {
	list, _, err := biz.NewApiTopic().ArtilceTop(limit, topic_name)
	if err != nil {
		return []*bizmodel.ApiArticleListModel{}
	}
	return list
}
