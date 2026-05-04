package global

import (
	"github.com/beego/beego/v2/server/web"
)

const (
	SuperFlag = "supper"
)

var (
	// CMSConfig   *cms.CmsConfig
	// WxMpConfig  *cms.WechatMpConfig
	// WxMchConfig *mkt.WechatMchConfig
	AllowLogin = true
	WebSite    string
	PayKey     string
	WxOfPayUrl string

	ThemePath = "/static/themes/"
)

func init() {
	// cmsAction := cms.NewCmsConfigModel(models.DB_CMS)
	// find, err := cmsAction.FindOneByUserId(0)
	// if err != nil {
	// 	CMSConfig = new(cms.CmsConfig)
	// } else {
	// 	CMSConfig = find
	// }

	// cmsAction1 := cms.NewWechatMpConfigModel(models.DB_CMS)
	// find1, err := cmsAction1.FindOneByUserId(0)
	// if err != nil {
	// 	WxMpConfig = new(cms.WechatMpConfig)
	// } else {
	// 	WxMpConfig = find1
	// }

	// mktAction := mkt.NewWechatMchConfigModel(models.DB_Market)
	// findx, err := mktAction.FindOneByUserId(0)
	// if err != nil {
	// 	WxMchConfig = new(mkt.WechatMchConfig)
	// } else {
	// 	WxMchConfig = findx
	// }

	AllowLogin = web.AppConfig.DefaultBool("allow_login", true)
	PayKey = web.AppConfig.DefaultString("pay_key", "")
	WebSite = web.AppConfig.DefaultString("web_site", "")
}
