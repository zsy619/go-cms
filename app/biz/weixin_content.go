package biz

import (
	"fmt"

	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/query"
)

type WeixinContent struct{}

func NewWeixinContent() *WeixinContent {
	return &WeixinContent{}
}

// ContentPaginate 分页
func (svc *WeixinContent) ContentPaginate(page, limit int, accountId int64) ([]*bizmodel.Weixin_ContentModel, int64, error) {
	_, do := query.WeixinResponseContentDo()
	field := `a.*,b.name as account_name`
	sqlCount := `SELECT Count(1) as count FROM weixin_response_content a LEFT JOIN weixin_account b ON a.account_id=b.account_id`
	sql := `SELECT ` + field + ` FROM weixin_response_content a LEFT JOIN weixin_account b ON a.account_id=b.account_id`
	if accountId > 0 {
		sqlCount += fmt.Sprintf(` WHERE a.account_id = %d`, accountId)
		sql += fmt.Sprintf(` WHERE a.account_id = %d`, accountId)
	}
	sql += ` ORDER BY a.add_time DESC`
	sql += fmt.Sprintf(` LIMIT %d,%d`, (page-1)*limit, limit)
	list := []*bizmodel.Weixin_ContentModel{}
	count := int64(0)
	do.UnderlyingDB().Raw(sql).Find(&list)
	do.UnderlyingDB().Raw(sqlCount).Pluck("count", &count)
	return list, count, nil
}
