package funcs

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
func TagNewExt(limit int, site_id int64, site_flag string, channel_id int64) []*bizmodel.ApiTagModel {
	find, _, _ := biz.NewApiTag().GetNew(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*bizmodel.ApiTagModel{}
	}
	return find
}

func TagNew(limit int) []*bizmodel.ApiTagModel {
	return TagNewExt(limit, 0, "", 0)
}

/**
 * @description: 获取最新标签列表
 * @param {int} limit 限制数量
 * @param {*} site_flag 站点标识
 * @param {*} site_id 站点ID
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func TagTopExt(limit int, site_id int64, site_flag string, channel_id int64) []*bizmodel.ApiTagModel {
	find, _, _ := biz.NewApiTag().Get(limit, site_id, site_flag, channel_id)
	if find == nil {
		return []*bizmodel.ApiTagModel{}
	}
	return find
}

func TagTop(limit int) []*bizmodel.ApiTagModel {
	return TagTopExt(limit, 0, "", 0)
}

func TagArtilceTop(limit int, tag_name string) []*bizmodel.ApiArticleListModel {
	list, _, err := biz.NewApiTag().ArtilceTop(limit, tag_name)
	if err != nil {
		return []*bizmodel.ApiArticleListModel{}
	}
	return list
}
