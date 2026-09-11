package www

import (
	"log"
	"regexp"
	"strconv"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
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

/**
 * @description: SiteLanguageList 获取站点语言关联列表
 * 返回所有站点及其关联的语言信息（基于 SQL JOIN 查询）
 * @return {*}
 */
// @router /api/site/language/list [get]
func (ctrl *ApiSiteController) SiteLanguageList() {
	list, total, err := service.NewSiteLanguageService().SiteLanguageList()
	if err != nil {
		logs.Error("SiteLanguageList query failed:", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), nil, 0)
		return
	}
	ctrl.JSONPageSuccess(list, total)
}

/**
 * @description: SiteLanguageBySiteID 根据站点ID获取语言信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
// @router /api/site/language/find [get]
func (ctrl *ApiSiteController) SiteLanguageBySiteID() {
	siteID, _ := ctrl.GetInt64("site_id")
	if siteID == 0 {
		ctrl.JSONError("site_id is required")
		return
	}
	list, err := service.NewSiteLanguageService().SiteLanguageListBySiteID(siteID)
	if err != nil {
		logs.Error("SiteLanguageBySiteID query failed:", err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("获取成功", list)
}

/**
 * @description: SiteLanguageByCode 根据语言代码获取使用该语言的站点列表
 * @param {string} code 语言代码
 * @return {*}
 */
// @router /api/site/language/by-code [get]
func (ctrl *ApiSiteController) SiteLanguageByCode() {
	code := ctrl.GetStringTrim("code", "")
	if code == "" {
		ctrl.JSONError("code is required")
		return
	}
	list, err := service.NewSiteLanguageService().SiteLanguageListByCode(code)
	if err != nil {
		logs.Error("SiteLanguageByCode query failed:", err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("获取成功", list)
}
