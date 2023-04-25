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
 * @description: Find 获取站点信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
// @router /api/site/get [get]
// @router /api/site/get/:site_id:int64 [get]
func (this *ApiSiteController) Find() {
	site_id, _ := this.GetInt64("site_id")
	if site_id == 0 {
		site_idx := this.Ctx.Input.Param(":site_id")
		site_id, _ = strconv.ParseInt(site_idx, 10, 64)
	}
	out, err := this.BaseController.SiteFind(site_id)
	if err != nil {
		logs.Error("Site Find::", "siteId", site_id, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess("", out)
}

/*
 * @description: ChannelGet 获取站点频道
 * @param {int64} site_id 站点ID
 * @return {*}
 */
// @router /api/channel/get [get]
func (this *ApiSiteController) ChannelGet() {
	site_id, _ := this.GetInt64("site_id")
	out, len, err := this.BaseController.ChannelGet(site_id)
	if err != nil {
		logs.Error("Channel Get::", "siteId", site_id, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}

/**
 * @description: Menu 获取站点菜单
 * @return {*}
 */
// @router /api/site/menu [get]
// @router /api/site/menu/:site_id:int64 [get]
func (this *ApiSiteController) Menu() {
	site_id, _ := this.GetInt64("site_id")
	if site_id == 0 {
		site_idx := this.Ctx.Input.Param(":site_id")
		site_id, _ = strconv.ParseInt(site_idx, 10, 64)
	}
	channel_id, _ := this.GetInt64("channel_id")
	out, _, err := this.BaseController.SiteMenu(site_id, channel_id)
	if err != nil {
		logs.Error("Site Menu::", "siteId", site_id, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess("", out)
}

/**
 * @description: MenuFlag 获取站点菜单
 * @return {*}
 */
// @router /api/site/menu/flag [get]
// @router /api/site/menu/flag/:site_flag:string [get]
func (this *ApiSiteController) MenuFlag() {
	site_flag := this.GetString("site_flag")
	if site_flag == "" {
		site_flag = this.Ctx.Input.Param(":site_flag")
	}
	channel_id, _ := this.GetInt64("channel_id")
	out, _, err := this.BaseController.SiteMenuFlag(site_flag, channel_id)
	if err != nil {
		logs.Error("Site MenuFlag::", "site_flag", site_flag, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess("", out)
}
