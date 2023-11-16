package admin

import (
	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/biz"
)

// Response 消息记录
func (ctrl *WeixinController) Response() {
	{
		list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
		ctrl.Data["accountList"] = list
	}
	ctrl.display()
}

// ContentPaginate 消息记录分页
func (ctrl *WeixinController) ContentPaginate() {
	page, limit := ctrl.GetPagingParameters()
	accountId, _ := ctrl.GetInt64("account_id")
	list, total, err := biz.NewWeixinContent().ContentPaginate(page, limit, accountId)
	if err != nil {
		logs.Error("ContentPaginate", err.Error())
		ctrl.JSONPageError(err.Error(), list, total)
	}
	ctrl.JSONPageSuccess(list, total)
}
