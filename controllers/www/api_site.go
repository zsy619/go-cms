package www

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
)

type ApiSiteController struct{ BaseController }

// Get 获取站点信息
// @router /api/site/get [get]
func (this *ApiSiteController) Get() {
	site_id, _ := this.GetInt64("site_id")
	out, err := biz.NewApiSite().Get(site_id)
	if err != nil {
		logs.Error("Site Get::", "siteId", site_id, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess("", out)
}

// ChannelFind 获取站点栏目
// @router /api/channel/find [get]
func (this *ApiSiteController) ChannelFind() {
	site_id, _ := this.GetInt64("site_id")
	out, len, err := biz.NewApiSite().ChannelFind(site_id)
	if err != nil {
		logs.Error("Channel Find::", "siteId", site_id, "err", err)
		this.JSONErrorOfData(err.Error(), out)
	}
	this.JSONSuccess(strconv.FormatInt(len, 10), out)
}
