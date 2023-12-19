package www

import (
	"log"
	"regexp"
	"strconv"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/app/lib"
)

type ApiSiteController struct{ BaseController }

/**
 * @description: Default 默认站点信息
 * @return {*}
 */
// @router /api/site/default [get]
func (ctrl *ApiSiteController) Default() {
	out, err := ctrl.BaseController.SiteDefault()
	if err != nil {
		logs.Error(err.Error())
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, 1)
}

/**
 * @description: 通过域名获取站点信息（使用缓存）
 * @param {string} domain 域名
 * @return {*}
 */
// @router /api/site/find/domain [get]
func (ctrl *ApiSiteController) FindDomain() {
	url := "http://" + ctrl.Ctx.Request.Host
	patt := `^((http://)|(https://))?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}(/)`
	reg := regexp.MustCompile(patt)
	preUrl := reg.FindString(url)
	// 查询站点表获取对应站点ID
	domainMdl := service.NewCmsSiteDomainModel().One(preUrl)
	if domainMdl == nil || domainMdl.SiteID > 0 {
		ctrl.JSONError("获取数据失败")
	}
	// 获取站点信息
	siteMdl, err := service.NewCmsSite().SiteOne(domainMdl.SiteID)
	if err != nil {
		ctrl.JSONError("获取数据失败")
	}
	log.Println(siteMdl)
	ctrl.JSONSuccess("获取成功", siteMdl)
}

/*
 * @description: Find 获取站点信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
// @router /api/site/find [get]
// @router /api/site/find/:site_id:int64 [get]
func (ctrl *ApiSiteController) Find() {
	site_id, _ := ctrl.GetInt64("site_id")
	if site_id == 0 {
		site_idx := ctrl.Ctx.Input.Param(":site_id")
		site_id, _ = strconv.ParseInt(site_idx, 10, 64)
	}
	out, err := ctrl.BaseController.SiteFind(site_id)
	if err != nil {
		logs.Error("Site Find::", "siteId", site_id, "err", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, 1)
}

/*
 * @description: ChannelGet 获取站点频道
 * @param {int64} site_id 站点ID
 * @return {*}
 */
// @router /api/channel/get [get]
func (ctrl *ApiSiteController) ChannelGet() {
	site_id, _ := ctrl.GetInt64("site_id")
	out, len, err := ctrl.BaseController.ChannelGet(site_id)
	if err != nil {
		logs.Error("Channel Get::", "siteId", site_id, "err", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, len)
}

/**
 * @description: Menu 获取站点菜单
 * @return {*}
 */
// @router /api/site/menu [get]
// @router /api/site/menu/:site_id:int64 [get]
func (ctrl *ApiSiteController) Menu() {
	site_id, _ := ctrl.GetInt64("site_id")
	if site_id == 0 {
		site_idx := ctrl.Ctx.Input.Param(":site_id")
		site_id, _ = strconv.ParseInt(site_idx, 10, 64)
	}
	channel_id, _ := ctrl.GetInt64("channel_id")
	out, len, err := ctrl.BaseController.SiteMenu(site_id, channel_id)
	if err != nil {
		logs.Error("Site Menu::", "siteId", site_id, "err", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), out, 0)
	}
	ctrl.JSONPageSuccess(out, len)
}

/**
 * @description: MenuFlag 获取站点菜单
 * @return {*}
 */
// @router /api/site/menu/flag [get]
// @router /api/site/menu/flag/:site_flag:string [get]
func (ctrl *ApiSiteController) MenuFlag() {
	site_flag := ctrl.GetSafeString("site_flag")
	if site_flag == "" {
		site_flag = ctrl.Ctx.Input.Param(":site_flag")
	}
	channel_id, _ := ctrl.GetInt64("channel_id")
	out, _, err := ctrl.BaseController.SiteMenuFlag(site_flag, channel_id)
	if err != nil {
		logs.Error("Site MenuFlag::", "site_flag", site_flag, "err", err)
		ctrl.JSONErrorOfData(err.Error(), out)
	}
	ctrl.JSONSuccess("", out)
}
