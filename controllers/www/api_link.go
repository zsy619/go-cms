package www

import (
	"strconv"

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
func (this *ApiLinkController) Get() {
	site_flag := this.GetString("site_flag")
	site_id, _ := this.GetInt64("site_id")
	category_id, _ := this.GetInt64("category_id")
	call_index := this.GetString("call_index")
	limit, _ := this.GetInt("limit", 6)
	out, len, err := this.BaseController.LinkGet(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		logs.Error("Get::", "callIndex", call_index, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
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
func (this *ApiLinkController) GetNew() {
	site_flag := this.GetString("site_flag")
	site_id, _ := this.GetInt64("site_id")
	category_id, _ := this.GetInt64("category_id")
	call_index := this.GetString("call_index")
	limit, _ := this.GetInt("limit", 6)
	out, len, err := this.BaseController.LinkGetNew(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		logs.Error("GetNew::", "callIndex", call_index, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}

/**
 * @description: Paginate 获取链接列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {int64} category_id 链接分类ID
 * @param {string} call_index 链接分类标识
 * @return {*}
 */
// @router /api/link/paginate [get]
func (this *ApiLinkController) Paginate() {
	category_id, _ := this.GetInt64("category_id")
	call_index := this.GetString("call_index")
	limit, _ := this.GetInt("limit", 12)
	page, _ := this.GetInt("page", 1)
	out, len, err := this.BaseController.LinkPaginate(page, limit, category_id, call_index)
	if err != nil {
		logs.Error("Paginate::", "callIndex", call_index, "err", err)
		this.JSONPage(lib.CodeError, err.Error(), out, len)
	}
	this.JSONPageSuccess(out, len)
}

/**
 * @description: Click 点击数+1
 * @param {int64} link_id 链接ID
 * @return {*}
 */
// @router /api/link/click [get]
func (this *ApiLinkController) Click() {
	link_id, _ := this.GetInt64("link_id", 0)
	err := this.BaseController.LinkClick(link_id)
	if err != nil {
		logs.Error("Click", err.Error())
		this.JSONError(err.Error())
	}
	this.JSONSuccess("", nil)
}
