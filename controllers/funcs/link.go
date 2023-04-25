package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
* @description: LinkNewExt 获取最新广告列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func LinkNewExt(limit int, site_id int64, site_flag string, category_id int64, call_index string) []*bizmodel.ApiLinkModel {
	find, _, err := biz.NewApiLink().GetNew(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		find = []*bizmodel.ApiLinkModel{}
	}
	return find
}

func LinkNew(limit int) []*bizmodel.ApiLinkModel {
	return LinkNewExt(limit, 0, "", 0, "")
}

/**
* @description: LinkTopExt 获取广告列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func LinkTopExt(limit int, site_id int64, site_flag string, category_id int64, call_index string) []*bizmodel.ApiLinkModel {
	find, _, err := biz.NewApiLink().Get(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		find = []*bizmodel.ApiLinkModel{}
	}
	return find
}

func LinkTop(limit int) []*bizmodel.ApiLinkModel {
	return LinkTopExt(limit, 0, "", 0, "")
}
