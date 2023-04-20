package www

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
)

type ApiTagController struct{ BaseController }

// @router /api/tag/find [get]
func (this *ApiTagController) Find() {
	site_id, _ := this.GetInt64("site_id")
	channel_id, _ := this.GetInt64("channel_id")
	limit, _ := this.GetInt("limit", 6)
	out, len, err := this.BaseController.TagFind(limit, site_id, channel_id)
	if err != nil {
		logs.Error("Find::", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}

// @router /api/tag/paginate [get]
func (this *ApiTagController) FindNew() {
	site_id, _ := this.GetInt64("site_id")
	channel_id, _ := this.GetInt64("channel_id")
	limit, _ := this.GetInt("limit", 6)
	out, len, err := this.BaseController.TagFindNew(limit, site_id, channel_id)
	if err != nil {
		logs.Error("Find::", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}

/**
 * @description: Click 点击数+1
 * @param {int64} ads_id 广告ID
 * @return {*}
 */
// @router /api/tag/click [get]
func (this *ApiTagController) Click() {
	tag_id, _ := this.GetInt64("tag_id", 0)
	err := this.BaseController.TagClick(tag_id)
	if err != nil {
		logs.Error("Click", err.Error())
		this.JSONError(err.Error())
	}
	this.JSONSuccess("", nil)
}
