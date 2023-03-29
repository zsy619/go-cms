package www

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/lib"
)

type ApiLinkController struct{ BaseController }

/**
 * @description: Find 获取链接列表
 * @param {int} limit 获取数量
 * @param {string} call_index 链接分类标识
 * @return {*}
 */
// @router /api/link/find [get]
func (this *ApiLinkController) Find() {
	call_index := this.GetString("call_index")
	limit, _ := this.GetInt("limit", 6)
	out, len, err := biz.NewApiLink().Find(limit, call_index)
	if err != nil {
		logs.Error("Find::", "callIndex", call_index, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}

/**
 * @description: Paginate 获取链接列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {string} call_index 链接分类标识
 * @return {*}
 */
// @router /api/link/paginate [get]
func (this *ApiLinkController) Paginate() {
	call_index := this.GetString("call_index")
	limit, _ := this.GetInt("limit", 12)
	page, _ := this.GetInt("page", 1)
	out, len, err := biz.NewApiLink().Paginate(page, limit, call_index)
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
	err := biz.NewApiLink().Click(link_id)
	if err != nil {
		logs.Error("Click", err.Error())
		this.JSONError(err.Error())
	}
	this.JSONSuccess("", nil)
}
