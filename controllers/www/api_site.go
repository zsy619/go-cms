package www

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
)

type ApiSiteController struct{ BaseController }

/**
 * @description: Default 获取站点信息
 * @return {*}
 */
// @router /api/site/default [get]
func (this *ApiSiteController) Default() {
	out, err := this.BaseController.SiteDefault()
	if err != nil {
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess("", out)
}

/*
 * @description: Get 获取站点信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
// @router /api/site/get [get]
func (this *ApiSiteController) Get() {
	site_id, _ := this.GetInt64("site_id")
	out, err := this.BaseController.SiteGet(site_id)
	if err != nil {
		logs.Error("Site Get::", "siteId", site_id, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess("", out)
}

/*
 * @description: ChannelFind 获取站点频道
 * @param {int64} site_id 站点ID
 * @return {*}
 */
// @router /api/channel/find [get]
func (this *ApiSiteController) ChannelFind() {
	site_id, _ := this.GetInt64("site_id")
	out, len, err := this.BaseController.ChannelFind(site_id)
	if err != nil {
		logs.Error("Channel Find::", "siteId", site_id, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}

/**
 * @description: Menu 获取站点菜单
 * @return {*}
 */
// @router /api/site/menu [get]
func (this *ApiSiteController) Menu() {
	site_id, _ := this.GetInt64("site_id")
	channel_id, _ := this.GetInt64("channel_id")
	out, _, err := this.BaseController.SiteMenu(site_id, channel_id)
	if err != nil {
		logs.Error("Site Menu::", "siteId", site_id, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess("", out)
}
