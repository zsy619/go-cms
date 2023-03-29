package biz

import (
	"errors"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/dal/query"
)

var Cache_LinkFindByCategory map[string][]map[string]interface{}

func init() {
	Cache_LinkFindByCategory = make(map[string][]map[string]interface{})
}

type ApiLink struct{}

func NewWebRoot() *ApiLink {
	return &ApiLink{}
}

// FindByCategory 获取链接列表
// 排序规则：is_top desc,order_id asc
func (this *ApiLink) FindByCategory(callIndex string) ([]map[string]interface{}, int64, error) {
	outLink := []map[string]interface{}{}
	if callIndex == "" {
		return outLink, 0, errors.New("调用链接分类标识不能为空")
	}
	if links, ok := Cache_LinkFindByCategory[callIndex]; ok {
		logs.Debug("LinkFindByCategory[Cache]::", "callIndex", callIndex, "links", links)
		// fmt.Println("LinkFindByCategory[Cache]::", "callIndex", callIndex, "links", links)
		return links, int64(len(links)), nil
	}
	_, linkDo := query.CmsLinkDo()
	sqlSelect := "a.link_id,a.site_id,a.channel_id,a.category_id,a.title,a.link_url,a.target,a.click,a.img_url,a.is_lock,a.is_red,a.is_hot,a.is_slide"
	sql := "SELECT " + sqlSelect + " FROM cms_link a LEFT JOIN cms_link_category b ON a.category_id = b.category_id WHERE b.call_index=? AND a.`status`=2 ORDER BY a.is_top desc,a.sort_id ASC"
	err := linkDo.UnderlyingDB().Raw(sql, callIndex).Scan(&outLink).Error
	if err == nil {
		Cache_LinkFindByCategory[callIndex] = outLink
	}
	return outLink, int64(len(outLink)), err
}
