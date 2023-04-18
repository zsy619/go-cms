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

/**
* @description: Find 获取链接列表
* @param {int} limit 获取数量
* @param {int64} category_id 链接分类ID
* @param {string} call_index 链接分类标识
* @return {*}
 */
func (this *ApiLink) Find(limit int, category_id int64, call_index string) ([]*bizmodel.ApiLinkListModel, int64, error) {
	cacheKey := fmt.Sprintf("LinkFind::%d::%d::%s", limit, category_id, call_index)
	if found, item := ApiCache.Get(cacheKey); found {
		links := item.([]*bizmodel.ApiLinkListModel)
		logs.Debug("LinkFindByCategory[Cache]::", "cacheKey", cacheKey, "links", links)
		return links, int64(len(links)), nil
	}
	outLink := []*bizmodel.ApiLinkListModel{}

	_, linkDo := query.CmsLinkDo()
	sqlSelect := "a.link_id,a.site_id,a.channel_id,a.category_id,a.title,a.link_url,a.target,a.click,a.img_url1,a.img_url2,a.is_lock,a.is_red,a.is_hot,a.is_slide"
	sql := "SELECT " + sqlSelect + " FROM cms_link a LEFT JOIN cms_link_category b ON a.category_id = b.category_id WHERE a.`status`=2 " +
		xgeneric.IFF(call_index == "", "", " AND b.call_index = '"+call_index+"'") +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id = "+strconv.FormatInt(category_id, 10)) +
		" ORDER BY a.is_top desc,a.sort_id ASC"
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
