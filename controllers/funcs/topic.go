package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
 * @description: 获取专题列表
 * @param {int} limit 限制数量
 * @param {*} siteId 站点ID
 * @param {int64} channelId 栏目ID
 * @return {*}
 */
func TopicNewExt(limit int, siteId, channelId int64) []*bizmodel.ApiTopicListModel {
	find, _, _ := biz.NewApiTopic().GetNew(limit, siteId, channelId)
	if find == nil {
		return []*bizmodel.ApiTopicListModel{}
	}
	return find
}

func TopicNew(limit int) []*bizmodel.ApiTopicListModel {
	return TopicNewExt(limit, 0, 0)
}

/**
 * @description: 获取最新专题列表
 * @param {int} limit 限制数量
 * @param {*} siteId 站点ID
 * @param {int64} channelId 栏目ID
 * @return {*}
 */
func TopicTopExt(limit int, siteId, channelId int64) []*bizmodel.ApiTopicListModel {
	find, _, _ := biz.NewApiTopic().Get(limit, siteId, channelId)
	if find == nil {
		return []*bizmodel.ApiTopicListModel{}
	}
	return find
}

func TopicTop(limit int) []*bizmodel.ApiTopicListModel {
	return TopicTopExt(limit, 0, 0)
}

func TopicArtilceTop(limit int, topic_name string) []*bizmodel.ApiArticleListModel {
	list, _, err := biz.NewApiTopic().ArtilceTop(limit, topic_name)
	if err != nil {
		return []*bizmodel.ApiArticleListModel{}
	}
	return list
}
