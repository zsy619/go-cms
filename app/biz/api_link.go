package biz

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xgeneric"
)

var Cache_LinkFindByCategory map[string][]map[string]interface{}

func init() {
	Cache_LinkFindByCategory = make(map[string][]map[string]interface{})
}

type ApiLink struct{}

func NewApiLink() *ApiLink {
	return &ApiLink{}
}

func (this *ApiLink) InitCache() {
	Cache_LinkFindByCategory = make(map[string][]map[string]interface{})
}

// Find 获取链接列表
// 排序规则：is_top desc,order_id asc
func (this *ApiLink) Find(limit int, call_index string) ([]map[string]interface{}, int64, error) {
	outLink := []map[string]interface{}{}
	if call_index == "" {
		return outLink, 0, errors.New("调用链接分类标识不能为空")
	}
	if links, ok := Cache_LinkFindByCategory[call_index]; ok {
		logs.Debug("LinkFindByCategory[Cache]::", "callIndex", call_index, "links", links)
		// fmt.Println("LinkFindByCategory[Cache]::", "callIndex", callIndex, "links", links)
		return links, int64(len(links)), nil
	}
	_, linkDo := query.CmsLinkDo()
	sqlSelect := "a.link_id,a.site_id,a.channel_id,a.category_id,a.title,a.link_url,a.target,a.click,a.img_url,a.is_lock,a.is_red,a.is_hot,a.is_slide"
	sql := "SELECT " + sqlSelect + " FROM cms_link a LEFT JOIN cms_link_category b ON a.category_id = b.category_id WHERE b.call_index=? AND a.`status`=2 ORDER BY a.is_top desc,a.sort_id ASC"
	if limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", limit)
	}
	err := linkDo.UnderlyingDB().Raw(sql, call_index).Scan(&outLink).Error
	if err == nil {
		Cache_LinkFindByCategory[call_index] = outLink
	}
	return outLink, int64(len(outLink)), err
}

// Paginate 获取链接列表
// 排序规则：is_top desc,order_id asc
func (this *ApiLink) Paginate(page, limit int, call_index string) ([]map[string]interface{}, int64, error) {
	outLink := []map[string]interface{}{}
	_, linkDo := query.CmsLinkDo()
	sqlSelectRow := "a.link_id,a.site_id,a.channel_id,a.category_id,a.title,a.link_url,a.target,a.click,a.img_url,a.is_lock,a.is_red,a.is_hot,a.is_slide"
	sqlRow := "SELECT " + sqlSelectRow + " FROM cms_link a LEFT JOIN cms_link_category b ON a.category_id = b.category_id WHERE " +
		xgeneric.IFF(call_index == "", "", "b.call_index='"+call_index+"' And") +
		" a.`status`=2 ORDER BY a.is_top desc,a.sort_id ASC"

	sqlSelectCount := "count(1) as count"
	sqlCount := "SELECT " + sqlSelectCount + " FROM cms_link a LEFT JOIN cms_link_category b ON a.category_id = b.category_id WHERE " +
		xgeneric.IFF(call_index == "", "", "b.call_index='"+call_index+"' And") +
		" a.`status`=2"

	var count int64
	linkDo.UnderlyingDB().Raw(sqlCount).Pluck("count", &count)

	sqlRow += fmt.Sprintf(" LIMIT %d,%d", (page-1)*limit, limit)
	err := linkDo.UnderlyingDB().Raw(sqlRow).Scan(&outLink).Error
	if err == nil {
		Cache_LinkFindByCategory[call_index] = outLink
	}
	return outLink, count, err
}

func (this *ApiLink) Click(link_id int64) error {
	mdl, do := query.CmsLinkDo()
	do.Where(mdl.LinkID.Eq(link_id), mdl.Status.Eq(2)).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}
