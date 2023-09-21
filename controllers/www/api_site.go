package www

import (
	"haedu.gov.cn/cms/app/biz"
	"log"
	"regexp"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/lib"
)

type ApiSiteController struct{ BaseController }

/**
 * @description: Default 默认站点信息
 * @return {*}
 */
// @router /api/site/default [get]
func (this *ApiSiteController) Default() {
	out, err := this.BaseController.SiteDefault()
	if err != nil {
		logs.Error(err.Error())
		this.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	this.JSONPageSuccess(out, 1)
}

/**
 * @description: 通过域名获取站点信息（使用缓存）
 * @param {string} domain 域名
 * @return {*}
 */
// @router /api/site/find/domain [get]
func (this *ApiSiteController) FindDomain() {
	url := "http://" + this.Ctx.Request.Host
	patt := `^((http://)|(https://))?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}(/)`
	reg := regexp.MustCompile(patt)
	preUrl := reg.FindString(url)
	// 查询站点表获取对应站点ID
	domainMdl := biz.NewCmsSiteDomainModel().One(preUrl)
	if domainMdl == nil || domainMdl.SiteID > 0 {
		this.JSONError("获取数据失败")
	}
	// 获取站点信息
	siteMdl, err := biz.NewCmsSite().SiteOne(domainMdl.SiteID)
	if err != nil {
		this.JSONError("获取数据失败")
	}
	log.Println(siteMdl)
	this.JSONSuccess("获取成功", siteMdl)
}

/*
 * @description: Find 获取站点信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
// @router /api/site/find [get]
// @router /api/site/find/:site_id:int64 [get]
func (this *ApiSiteController) Find() {
	site_id, _ := this.GetInt64("site_id")
	if site_id == 0 {
		site_idx := this.Ctx.Input.Param(":site_id")
		site_id, _ = strconv.ParseInt(site_idx, 10, 64)
	}
	out, err := this.BaseController.SiteFind(site_id)
	if err != nil {
		logs.Error("Site Find::", "siteId", site_id, "err", err)
		this.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	this.JSONPageSuccess(out, 1)
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
		this.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	this.JSONPageSuccess(out, len)
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
	out, len, err := this.BaseController.SiteMenu(site_id, channel_id)
	if err != nil {
		logs.Error("Site Menu::", "siteId", site_id, "err", err)
		this.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	this.JSONPageSuccess(out, len)
}

/**
 * @description: MenuFlag 获取站点菜单
 * @return {*}
 */
// @router /api/site/menu/flag [get]
// @router /api/site/menu/flag/:site_flag:string [get]
func (this *ApiSiteController) MenuFlag() {
	site_flag := this.GetSafeString("site_flag")
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
