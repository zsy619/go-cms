package www

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/lib"
)

type ApiAdController struct{ BaseController }

/**
* @description: Find 获取广告列表
* @param {int} limit 获取数量
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
// @router /api/ad/find [get]
func (this *ApiAdController) Find() {
	category_id, _ := this.GetInt64("category_id")
	call_index := this.GetString("call_index")
	limit, _ := this.GetInt("limit", 6)
	out, len, err := this.BaseController.AdFind(limit, category_id, call_index)
	if err != nil {
		logs.Error("Find::", "callIndex", call_index, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}

/**
 * @description: Paginate 获取广告列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {int64} category_id 广告分类ID
 * @param {string} call_index 广告分类标识
 * @return {*}
 */
// @router /api/ad/paginate [get]
func (this *ApiAdController) Paginate() {
	category_id, _ := this.GetInt64("category_id")
	call_index := this.GetString("call_index")
	limit, _ := this.GetInt("limit", 12)
	page, _ := this.GetInt("page", 1)
	out, len, err := this.BaseController.AdPaginate(page, limit, category_id, call_index)
	if err != nil {
		logs.Error("Paginate::", "callIndex", call_index, "err", err)
		this.JSONPage(lib.CodeError, err.Error(), out, len)
	}
	this.JSONPageSuccess(out, len)
}

/**
 * @description: Click 点击数+1
 * @param {int64} ad_id 广告ID
 * @return {*}
 */
// @router /api/link/click [get]
func (this *ApiAdController) Click() {
	ad_id, _ := this.GetInt64("ad_id", 0)
	err := this.BaseController.AdClick(ad_id)
	if err != nil {
		logs.Error("Click", err.Error())
		this.JSONError(err.Error())
	}
	this.JSONSuccess("", nil)
}
