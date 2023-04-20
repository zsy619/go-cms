package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
* @description: AdsNewExt 获取最新广告列表
* @param {int} limit 获取数量
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func AdsNewExt(limit int, category_id int64, call_index string) []*bizmodel.ApiAdsListModel {
	find, _, err := biz.NewApiAds().FindNew(limit, category_id, call_index)
	if err != nil {
		find = []*bizmodel.ApiAdsListModel{}
	}
	return find
}

func AdsNew(limit int) []*bizmodel.ApiAdsListModel {
	return AdsNewExt(limit, 0, "")
}

/**
* @description: AdsTopExt 获取广告列表
* @param {int} limit 获取数量
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func AdsTopExt(limit int, category_id int64, call_index string) []*bizmodel.ApiAdsListModel {
	find, _, err := biz.NewApiAds().Find(limit, category_id, call_index)
	if err != nil {
		find = []*bizmodel.ApiAdsListModel{}
	}
	return find
}

func AdsTop(limit int) []*bizmodel.ApiAdsListModel {
	return AdsTopExt(limit, 0, "")
}
