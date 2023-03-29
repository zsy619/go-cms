package www

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
)

type ApiLinkController struct{ BaseController }

// LinkFindByCategory 获取链接列表
// @router /api/link/find-by-category [get]
func (this *ApiLinkController) LinkFindByCategory() {
	callIndex := this.GetString("callIndex")
	out, len, err := biz.NewWebRoot().FindByCategory(callIndex)
	if err != nil {
		logs.Error("LinkFindByCategory::", "callIndex", callIndex, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}
