package www

import (
	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/lib"
)

type ApiLinkController struct{ BaseController }

/**
* @description: Get 获取链接列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
// @router /api/link/get [get]
func (ctrl *ApiLinkController) Get() {
	site_flag := ctrl.GetSafeString("site_flag")
	site_id, _ := ctrl.GetInt64("site_id")
	category_id, _ := ctrl.GetInt64("category_id")
	call_index := ctrl.GetSafeString("call_index")
	limit, _ := ctrl.GetInt("limit", 6)
	out, len, err := ctrl.BaseController.LinkGet(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		logs.Error("Get::", "callIndex", call_index, "err", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, len)
}

/**
* @description: GetNew 获取最新链接列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
// @router /api/link/get/new [get]
func (ctrl *ApiLinkController) GetNew() {
	site_flag := ctrl.GetSafeString("site_flag")
	site_id, _ := ctrl.GetInt64("site_id")
	category_id, _ := ctrl.GetInt64("category_id")
	call_index := ctrl.GetSafeString("call_index")
	limit, _ := ctrl.GetInt("limit", 6)
	out, len, err := ctrl.BaseController.LinkGetNew(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		logs.Error("GetNew::", "callIndex", call_index, "err", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, len)
}

/**
 * @description: Paginate 获取链接列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} category_id 链接分类ID
 * @param {string} call_index 链接分类标识
 * @return {*}
 */
// @router /api/link/paginate [get]
func (ctrl *ApiLinkController) Paginate() {
	site_id, _ := ctrl.GetInt64("site_id")
	site_flag := ctrl.GetSafeString("site_flag")
	category_id, _ := ctrl.GetInt64("category_id")
	call_index := ctrl.GetSafeString("call_index")
	limit, _ := ctrl.GetInt("limit", 12)
	page, _ := ctrl.GetInt("page", 1)
	out, len, err := ctrl.BaseController.LinkPaginate(page, limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		logs.Error("Paginate::", "callIndex", call_index, "err", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, len)
	}
	ctrl.JSONPageSuccess(out, len)
}

/**
 * @description: Click 点击数+1
 * @param {int64} link_id 链接ID
 * @return {*}
 */
// @router /api/link/click [get]
func (ctrl *ApiLinkController) Click() {
	link_id, _ := ctrl.GetInt64("link_id", 0)
	err := ctrl.BaseController.LinkClick(link_id)
	if err != nil {
		logs.Error("Click", err.Error())
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("", nil)
}
