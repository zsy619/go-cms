package funcs

import (
	"haedu.gov.cn/cms/app/cms/service"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
)

/**
 * @description: 获取默认站点
 * @return {*}
 */
func SiteDefault() *service_model.ApiSiteModel {
	find, err := service.NewApiSite().Default()
	if err != nil {
		return &service_model.ApiSiteModel{}
	}
	return find
}

/**
 * @description: 站点菜单
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func SiteMenu(site_id int64) []*service_model.ApiNavModel {
	find, _, err := service.NewApiSite().NavGet(site_id, 0)
	if err != nil {
		return []*service_model.ApiNavModel{}
	}
	return find
}

/**
 * @description: 站点菜单
 * @param {string} site_flag 站点标识
 * @return {*}
 */
func SiteMenuFlag(site_flag string) []*service_model.ApiNavModel {
	find, _, err := service.NewApiSite().NavGetByFlag(site_flag, 0)
	if err != nil {
		return []*service_model.ApiNavModel{}
	}
	return find
}

/**
 * @description: 站点频道
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func SiteChannel(site_id int64) []*service_model.ApiChannelModel {
	find, _, err := service.NewApiSite().ChannelGet(site_id)
	if err != nil {
		return []*service_model.ApiChannelModel{}
	}
	return find
}
