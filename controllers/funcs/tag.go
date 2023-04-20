package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
 * @description: 获取标签列表
 * @param {int} limit 限制数量
 * @param {*} siteId 站点ID
 * @param {int64} channelId 栏目ID
 * @return {*}
 */
func TagNewExt(limit int, siteId, channelId int64) []*bizmodel.ApiTagListModel {
	find, _, _ := biz.NewApiTag().FindNew(limit, siteId, channelId)
	if find == nil {
		return []*bizmodel.ApiTagListModel{}
	}
	return find
}

func TagNew(limit int) []*bizmodel.ApiTagListModel {
	return TagNewExt(limit, 0, 0)
}

/**
 * @description: 获取最新标签列表
 * @param {int} limit 限制数量
 * @param {*} siteId 站点ID
 * @param {int64} channelId 栏目ID
 * @return {*}
 */
func TagTopExt(limit int, siteId, channelId int64) []*bizmodel.ApiTagListModel {
	find, _, _ := biz.NewApiTag().Find(limit, siteId, channelId)
	if find == nil {
		return []*bizmodel.ApiTagListModel{}
	}
	return find
}

func TagTop(limit int) []*bizmodel.ApiTagListModel {
	return TagTopExt(limit, 0, 0)
}
