package biz

import (
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xgeneric"
)

type ApiLink struct{}

func NewApiLink() *ApiLink {
	return &ApiLink{}
}

func (this *ApiLink) get(prefixCache string, limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiLinkListModel, int64, error) {
	cacheKey := fmt.Sprintf("%s_%d_%d_%s_%d_%s", prefixCache, limit, site_id, site_flag, category_id, call_index)
	if found, item := ApiCache.Get(cacheKey); found {
		links := item.([]*bizmodel.ApiLinkListModel)
		logs.Debug("Get[Cache]::", "cacheKey", cacheKey, "links", links)
		return links, int64(len(links)), nil
	}
	outLink := []*bizmodel.ApiLinkListModel{}

	_, linkDo := query.CmsLinkDo()
	sqlSelect := "a.link_id,a.site_id,a.channel_id,a.category_id,a.title,a.link_url,a.target,a.click,a.img_url1,a.img_url2,a.is_lock,a.is_red,a.is_hot,a.is_slide"
	sql := `SELECT ` + sqlSelect + ` FROM cms_link a` +
		` LEFT JOIN cms_link_category b ON a.category_id = b.category_id` +
		` LEFT JOIN cms_site c ON a.site_id = c.site_id` +
		` WHERE a.status=2` +
		xgeneric.IFF(site_flag == "", "", " AND c.flag = '"+site_flag+"'") +
		xgeneric.IFF(site_id <= 0, "", " AND b.site_id = "+strconv.FormatInt(site_id, 10)) +
		xgeneric.IFF(call_index == "", "", " AND b.call_index = '"+call_index+"'") +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id = "+strconv.FormatInt(category_id, 10))
	if prefixCache == "ApiLink_Get" {
		sql += " ORDER BY a.is_top DESC,a.sort_id ASC"
	} else {
		sql += " ORDER BY a.create_time DESC,a.sort_id ASC"
	}
	if limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", limit)
	}
	err := linkDo.UnderlyingDB().Raw(sql).Scan(&outLink).Error
	if err == nil {
		ApiCache.Set(cacheKey, outLink, 1800)
	}
	return outLink, int64(len(outLink)), err
}

/**
* @description: Get 获取链接列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
func (this *ApiLink) Get(limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiLinkListModel, int64, error) {
	return this.get("ApiLink_Get", limit, site_id, site_flag, category_id, call_index)
}

/**
* @description: GetNew 获取最新链接列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
func (this *ApiLink) GetNew(limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiLinkListModel, int64, error) {
	return this.get("ApiLink_GetNew", limit, site_id, site_flag, category_id, call_index)
}

/**
 * @description: Paginate 获取链接列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {int64} category_id 链接分类ID
 * @param {string} call_index 链接分类标识
 * @return {*}
 */
func (this *ApiLink) Paginate(page, limit int, category_id int64, call_index string) ([]*bizmodel.ApiLinkListModel, int64, error) {
	outLink := []*bizmodel.ApiLinkListModel{}
	_, linkDo := query.CmsLinkDo()
	sqlSelectRow := "a.link_id,a.site_id,a.channel_id,a.category_id,a.title,a.link_url,a.target,a.click,a.img_url1,a.img_url2,a.is_lock,a.is_red,a.is_hot,a.is_slide"
	sqlRow := "SELECT " + sqlSelectRow + " FROM cms_link a LEFT JOIN cms_link_category b ON a.category_id = b.category_id WHERE a.`status`=2" +
		xgeneric.IFF(call_index == "", "", " AND b.call_index = '"+call_index+"'") +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id = "+strconv.FormatInt(category_id, 10)) +
		" ORDER BY a.is_top desc,a.sort_id ASC"

	sqlSelectCount := "count(1) as count"
	sqlCount := "SELECT " + sqlSelectCount + " FROM cms_link a LEFT JOIN cms_link_category b ON a.category_id = b.category_id WHERE a.`status`=2" +
		xgeneric.IFF(call_index == "", "", " AND b.call_index = '"+call_index+"'") +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id = "+strconv.FormatInt(category_id, 10))

	var count int64
	linkDo.UnderlyingDB().Raw(sqlCount).Pluck("count", &count)

	sqlRow += fmt.Sprintf(" LIMIT %d,%d", (page-1)*limit, limit)
	err := linkDo.UnderlyingDB().Raw(sqlRow).Scan(&outLink).Error
	return outLink, count, err
}

/**
 * @description: Click 点击数+1
 * @param {int64} link_id 链接ID
 * @return {*}
 */
func (this *ApiLink) Click(link_id int64) error {
	mdl, do := query.CmsLinkDo()
	do.Where(mdl.LinkID.Eq(link_id), mdl.Status.Eq(2)).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}
