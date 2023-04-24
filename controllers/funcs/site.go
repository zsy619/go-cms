package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
 * @description: 获取默认站点
 * @return {*}
 */
func SiteDefault() *bizmodel.ApiSiteModel {
	find, err := biz.NewApiSite().Default()
	if err != nil {
		return &bizmodel.ApiSiteModel{}
	}
	return find
}

/**
 * @description: 站点菜单
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func SiteMenu(site_id int64) []*bizmodel.ApiNavFindModel {
	find, _, err := biz.NewApiSite().NavGet(site_id, 0)
	if err != nil {
		return []*bizmodel.ApiNavFindModel{}
	}
	return find
}
