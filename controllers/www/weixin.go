package www

import "haedu.gov.cn/cms/app/lib"

type WeixinController struct{ BaseController }

// 微信设置
func (c *WeixinController) Index() {
	c.Data["website"] = lib.C_LOCAL_DOMAIN_Backslash()
}

// 微信认证
func (c *WeixinController) AccountAuth() {
	accountId, _ := c.GetInt64("accountId")
	if accountId <= 0 {
		return
	}
}
