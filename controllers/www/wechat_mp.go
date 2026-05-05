package www

import (
	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/app/wechat/models"
	"haedu.gov.cn/cms/app/wechat/mp"
)

// WechatMpController 微信公众平台控制器
type WechatMpController struct{ BaseController }

// Index 微信设置页面
func (ctrl *WechatMpController) Index() {
	ctrl.Data["website"] = lib.C_LOCAL_DOMAIN_Backslash()
}

// Signature 微信签名验证
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
// @router /weixin/mp/index [get]
func (ctrl *WechatMpController) Signature() {
	accountId, _ := ctrl.GetInt64("accountId")
	if accountId <= 0 {
		ctrl.Abort("500")
		return
	}
	account, err := service.NewWeixinAccount().AccountFindCache(accountId)
	if err != nil {
		logs.Error("获取微信公众号配置失败: accountId=%d, error=%v", accountId, err)
		ctrl.Abort("500")
		return
	}

	method := ctrl.Ctx.Request.Method
	switch method {
	case "GET":
		{
			timestamp := ctrl.GetSafeString("timestamp")
			nonce := ctrl.GetSafeString("nonce")
			signatureIn := ctrl.GetSafeString("signature")
			echostr := ctrl.GetSafeString("echostr")
			verify := mp.NewTokenVerify(models.TokenParam{
				Timestamp: timestamp,
				Nonce:     nonce,
				Signature: signatureIn,
				EchoStr:   echostr,
			}, account.Token)
			rt, err := verify.Verify()
			if err == nil {
				ctrl.Ctx.WriteString(rt)
			} else {
				logs.Error("微信签名验证失败: error=%v", err)
				ctrl.Ctx.WriteString("")
			}
		}
	case "POST":
		ctrl.Message()
	}
	ctrl.StopRun()
}

// Message 微信消息处理
// @router /weixin/mp/index [post]
func (ctrl *WechatMpController) Message() {
	accountId, _ := ctrl.GetInt64("accountId")
	if accountId <= 0 {
		ctrl.Abort("500")
		return
	}
	_, err := service.NewWeixinAccount().AccountFindCache(accountId)
	if err != nil {
		logs.Error("获取微信公众号配置失败: accountId=%d, error=%v", accountId, err)
		ctrl.Abort("500")
		return
	}
	logs.Debug("微信消息处理: accountId=%d", accountId)

	method := ctrl.Ctx.Request.Method
	switch method {
	case "GET":
		ctrl.Signature()
	case "POST":
		{
		}
	}
	ctrl.StopRun()
}
