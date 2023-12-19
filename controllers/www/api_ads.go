package www

import (
	"github.com/beego/beego/v2/core/logs"

	lib "haedu.gov.cn/cms/app/tool"
)

type ApiAdsController struct{ BaseController }

/**
* @description: Get 获取广告列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
// @router /api/ads/get [get]
func (ctrl *ApiAdsController) Get() {
	site_flag := ctrl.GetSafeString("site_flag")
	site_id, _ := ctrl.GetInt64("site_id")
	category_id, _ := ctrl.GetInt64("category_id")
	call_index := ctrl.GetSafeString("call_index")
	limit, _ := ctrl.GetInt("limit", 6)
	out, len, err := ctrl.BaseController.AdsGet(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		logs.Error("Get::", "callIndex", call_index, "err", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, len)
}

/**
* @description: GetNew 获取最新广告列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
// @router /api/ads/get/new [get]
func (ctrl *ApiAdsController) GetNew() {
	site_flag := ctrl.GetSafeString("site_flag")
	site_id, _ := ctrl.GetInt64("site_id")
	category_id, _ := ctrl.GetInt64("category_id")
	call_index := ctrl.GetSafeString("call_index")
	limit, _ := ctrl.GetInt("limit", 6)
	out, len, err := ctrl.BaseController.AdsGetNew(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		logs.Error("GetNew::", "callIndex", call_index, "err", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, len)
}

/**
 * @description: Paginate 获取广告列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} category_id 广告分类ID
 * @param {string} call_index 广告分类标识
 * @return {*}
 */
// @router /api/ads/paginate [get]
func (ctrl *ApiAdsController) Paginate() {
	site_id, _ := ctrl.GetInt64("site_id")
	site_flag := ctrl.GetSafeString("site_flag")
	category_id, _ := ctrl.GetInt64("category_id")
	call_index := ctrl.GetSafeString("call_index")
	limit, _ := ctrl.GetInt("limit", 12)
	page, _ := ctrl.GetInt("page", 1)
	out, len, err := ctrl.BaseController.AdsPaginate(page, limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		logs.Error("Paginate::", "callIndex", call_index, "err", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, len)
	}
	ctrl.JSONPageSuccess(out, len)
}

/**
 * @description: Click 点击数+1
 * @param {int64} ads_id 广告ID
 * @return {*}
 */
// @router /api/ads/click [get]
func (ctrl *ApiAdsController) Click() {
	ads_id, _ := ctrl.GetInt64("ads_id", 0)
	err := ctrl.BaseController.AdsClick(ads_id)
	if err != nil {
		logs.Error("Click", err.Error())
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("", nil)
}
