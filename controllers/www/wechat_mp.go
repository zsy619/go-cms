package www

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/app/wechat/models"
	"haedu.gov.cn/cms/app/wechat/mp"
)

type WechatMpController struct{ BaseController }

// 微信设置
func (this *WechatMpController) Index() {
	this.Data["website"] = lib.C_LOCAL_DOMAIN_Backslash()
}

// 微信认证
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
// @router /weixin/mp/index [get]
func (this *WechatMpController) Signature() {
	accountId, _ := this.GetInt64("accountId")
	if accountId <= 0 {
		this.Abort("500")
		return
	}
	// 获取微信公众号配置
	account, err := biz.NewWeixinAccount().AccountFindCache(accountId)
	if err != nil {
		logs.Error(err)
		this.Abort("500")
		return
	}
	// fmt.Println(account)

	// 开发者通过检验signature对请求进行校验（下面有校验方式）。
	//	若确认此次GET请求来自微信服务器，请原样返回echostr参数内容，则接入生效，成为开发者成功，否则接入失败。
	// 加密/校验流程如下：
	// 1）将token、timestamp、nonce三个参数进行字典序排序
	// 2）将三个参数字符串拼接成一个字符串进行sha1加密
	// 3）开发者获得加密后的字符串可与signature对比，标识该请求来源于微信
	method := this.Ctx.Request.Method
	switch method {
	case "GET":
		{
			timestamp := this.GetSafeString("timestamp")
			nonce := this.GetSafeString("nonce")
			signatureIn := this.GetSafeString("signature")
			echostr := this.GetSafeString("echostr")
			verify := mp.NewTokenVerify(models.TokenParam{
				Timestamp: timestamp,
				Nonce:     nonce,
				Signature: signatureIn,
				EchoStr:   echostr,
			}, account.Token)
			rt, err := verify.Verify()
			if err == nil {
				this.Ctx.WriteString(rt)
			} else {
				logs.Error(err.Error())
				this.Ctx.WriteString("")
			}
		}
		break
	case "POST":
		this.Message()
		break
	}
	this.StopRun()
}

// @router /weixin/mp/index [post]
func (this *WechatMpController) Message() {
	accountId, _ := this.GetInt64("accountId")
	if accountId <= 0 {
		this.Abort("500")
		return
	}
	// 获取微信公众号配置
	account, err := biz.NewWeixinAccount().AccountFindCache(accountId)
	if err != nil {
		logs.Error(err)
		this.Abort("500")
		return
	}
	fmt.Println(account)

	method := this.Ctx.Request.Method
	switch method {
	case "GET":
		this.Signature()
		break
	case "POST":
		{
		}
		break
	}
	this.StopRun()
}
