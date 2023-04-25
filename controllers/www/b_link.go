package www

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
* @description: Get 获取链接列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
func (this *BaseController) LinkGet(limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiLinkModel, int64, error) {
	return biz.NewApiLink().Get(limit, site_id, site_flag, category_id, call_index)
}

/**
* @description: 获取最新链接列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
func (this *BaseController) LinkGetNew(limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiLinkModel, int64, error) {
	return biz.NewApiLink().GetNew(limit, site_id, site_flag, category_id, call_index)
}

/**
 * @description: LinkPaginate 获取链接列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} category_id 链接分类ID
 * @param {string} call_index 链接分类标识
 * @return {*}
 */
func (this *BaseController) LinkPaginate(page, limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiLinkModel, int64, error) {
	return biz.NewApiLink().Paginate(page, limit, site_id, site_flag, category_id, call_index)
}

/**
 * @description: Click 点击数+1
 * @param {int64} link_id 链接ID
 * @return {*}
 */
func (this *BaseController) LinkClick(link_id int64) error {
	return biz.NewApiLink().Click(link_id)
}
