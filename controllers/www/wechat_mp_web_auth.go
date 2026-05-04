package www

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/app/wechat/models"
	"haedu.gov.cn/cms/app/wechat/mp"
)

// MpWebAuth 网页授权
// 1、引导用户进入授权页面同意授权，获取code
// 2、通过code换取网页授权access_token（与基础支持中的access_token不同）
// 3、如果需要，开发者可以刷新网页授权access_token，避免过期
// 4、通过网页授权access_token和openid获取用户基本信息（支持UnionID机制）
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html#0
type WechatMpWebAuthController struct{ web.Controller }

func (ctrl *WechatMpWebAuthController) GetStringTrim(key string, def ...string) string {
	outStr := ctrl.GetString(key, def...)
	outStr = strings.TrimSpace(outStr)
	return outStr
}

// GET /wechat/mp/tooauth2
func (ctrl *WechatMpWebAuthController) ToOauth2() {
	accountId, _ := ctrl.GetInt64("accountId")
	if accountId <= 0 {
		ctrl.Abort("500")
		return
	}
	// 获取微信公众号配置
	account, err := service.NewWeixinAccount().AccountFindCache(accountId)
	if err != nil {
		logs.Error(err)
		ctrl.Abort("500")
		return
	}
	fmt.Println(account)
	from := ctrl.GetStringTrim("from")
	fmt.Println("from ----> ", from)
	// redirect_uri := global.WebSite + `wechat/mp/redirect_uri`
	// switch from {
	// case "bind":
	// 	redirect_uri = global.WebSite + `mobile/x/wechat/bind`
	// 	break
	// case "login":
	// 	redirect_uri = global.WebSite + `mobile/wechat/login`
	// 	break
	// case "register":
	// 	redirect_uri = global.WebSite + `/mobile/register/one`
	// 	break
	// }

	oauth := mp.NewOAuth2(&models.MpConfig{
		AppId:          account.AppID,
		AppSecret:      account.AppSecret,
		EncodingAesKey: account.AppAesKey,
		Token:          account.Token,
	})
	fmt.Println(oauth)
	// oauth.Redirect(ctrl.Ctx.ResponseWriter, ctrl.Ctx.Request, redirect_uri, "snsapi_userinfo", "STATE")

	// url, err := GetRedirectURL(redirect_uri, "snsapi_userinfo", "STATE")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(url)
	// c.Redirect(url, 302)
}

// GET /wechat/mp/redirect_uri
func (c *WechatMpWebAuthController) RedirectUri() {
	// code说明 ：code作为换取access_token的票据，每次用户授权带上的code将不一样，code只能使用一次，5分钟未被使用自动过期。
	// 	返回码	说明
	// 10003	redirect_uri域名与后台配置不一致
	// 10004	此公众号被封禁
	// 10005	此公众号并没有这些scope的权限
	// 10006	必须关注此测试号
	// 10009	操作太频繁了，请稍后重试
	// 10010	scope不能为空
	// 10011	redirect_uri不能为空
	// 10012	appid不能为空
	// 10013	state不能为空
	// 10015	公众号未授权第三方平台，请检查授权状态
	// 10016	不支持微信开放平台的Appid，请使用公众号Appid

	code := c.GetStringTrim("code")
	state := c.GetStringTrim("state")
	logs.Debug("code --> ", code, " state --> ", state)

	oauth := mp.NewOAuth2(&models.MpConfig{
		// AppId:          global.WxMpConfig.AppId,
		// AppSecret:      global.WxMpConfig.AppSecret,
		// EncodingAesKey: global.WxMpConfig.EncodingAesKey,
		// Token:          global.WxMpConfig.Token,
	})

	token, openid, err := oauth.GetGlobalUserAccessToken(code)
	if err != nil {
		logs.Error(err.Error())
		c.Ctx.WriteString(err.Error())
		return
	}

	userInfo, err := oauth.GetUserInfo(token, openid)
	if err != nil {
		return
	}
	fmt.Println("获取用户信息：", userInfo)
	c.Data["userinfo"] = userInfo
	c.Ctx.WriteString(userInfo.Nickname + ` --> ` + userInfo.OpenId)
}
