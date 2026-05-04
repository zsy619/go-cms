package funcs

import (
	"haedu.gov.cn/cms/app/cms/service"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
)

/**
* @description: AdsNewExtend 获取最新广告列表
* @param {int} limit 获取数量，小于等于0时按6条处理
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func AdsNewExtend(limit int, site_id int64, site_flag string, category_id int64, call_index string) []*service_model.ApiAdsModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, err := service.NewApiAds().GetNew(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		find = []*service_model.ApiAdsModel{}
	}
	return find
}

/**
 * @description:获取最新广告列表
 * @param {int} limit
 * @return {*}
 */
func AdsNew(limit int) []*service_model.ApiAdsModel {
	if limit <= 0 {
		limit = 6
	}
	return AdsNewExtend(limit, 0, "", 0, "")
}

/**
* @description: AdsTopExtend 获取广告列表
* @param {int} limit 获取数量，小于等于0时按6条处理
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func AdsTopExtend(limit int, site_id int64, site_flag string, category_id int64, call_index string) []*service_model.ApiAdsModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, err := service.NewApiAds().Get(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		find = []*service_model.ApiAdsModel{}
	}
	return find
}

/**
 * @description: 获取广告列表
 * @param {int} limit 获取数量,小于等于0时按6条处理
 * @return {*}
 */
func AdsTop(limit int) []*service_model.ApiAdsModel {
	if limit <= 0 {
		limit = 6
	}
	return AdsTopExtend(limit, 0, "", 0, "")
}
