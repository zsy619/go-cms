package www

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
* @description: AdsGet 获取广告列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func (ctrl *BaseController) AdsGet(limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiAdsModel, int64, error) {
	return biz.NewApiAds().Get(limit, site_id, site_flag, category_id, call_index)
}

/**
* @description: AdsGetNew 获取最新广告列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func (ctrl *BaseController) AdsGetNew(limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiAdsModel, int64, error) {
	return biz.NewApiAds().GetNew(limit, site_id, site_flag, category_id, call_index)
}

/**
 * @description: AdsPaginate 获取广告列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} category_id 广告分类ID
 * @param {string} call_index 广告分类标识
 * @return {*}
 */
func (ctrl *BaseController) AdsPaginate(page, limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiAdsModel, int64, error) {
	return biz.NewApiAds().Paginate(page, limit, site_id, site_flag, category_id, call_index)
}

/**
 * @description: AdsClick 点击数+1
 * @param {int64} ads_id 广告ID
 * @return {*}
 */
func (ctrl *BaseController) AdsClick(ads_id int64) error {
	return biz.NewApiAds().Click(ads_id)
}
