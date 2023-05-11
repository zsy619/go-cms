package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
* @description: LinkNewExtend 获取最新链接列表
* @param {int} limit 获取数量，小于等于0时按6条处理
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
func LinkNewExtend(limit int, site_id int64, site_flag string, category_id int64, call_index string) []*bizmodel.ApiLinkModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, err := biz.NewApiLink().GetNew(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		find = []*bizmodel.ApiLinkModel{}
	}
	return find
}

/**
 * @description:获取最新链接列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @return {*}
 */
func LinkNew(limit int) []*bizmodel.ApiLinkModel {
	if limit <= 0 {
		limit = 6
	}
	return LinkNewExtend(limit, 0, "", 0, "")
}

/**
* @description: LinkTopExtend 获取链接列表
* @param {int} limit 获取数量，小于等于0时按6条处理
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
func LinkTopExtend(limit int, site_id int64, site_flag string, category_id int64, call_index string) []*bizmodel.ApiLinkModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, err := biz.NewApiLink().Get(limit, site_id, site_flag, category_id, call_index)
	if err != nil {
		find = []*bizmodel.ApiLinkModel{}
	}
	return find
}

/**
 * @description: 获取链接列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @return {*}
 */
func LinkTop(limit int) []*bizmodel.ApiLinkModel {
	if limit <= 0 {
		limit = 6
	}
	return LinkTopExtend(limit, 0, "", 0, "")
}
