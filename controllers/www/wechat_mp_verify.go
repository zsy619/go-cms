package www

import (
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"haedu.gov.cn/cms/app/cms/service"
)

type WechatMpVerifyController struct{ web.Controller }

/**
 * @description: 初始化微信公众号验证文件路由
 * @return {*}
 */
func InitWechatMpVerifyRouter() {
	if list, err := service.NewWeixinMpVerify().GetCache(); err != nil {
		logs.Debug("InitMpVerifyRouter--->", err.Error())
	} else {
		if len(list) > 0 {
			logs.Debug("开始 初始化微信公众号验证文件路由")
			for _, item := range list {
				if item.Status != 2 {
					continue
				}
				router := item.Path
				if !strings.HasPrefix(router, "/") {
					router = "/" + router
				}
				if strings.HasSuffix(router, "/") {
					router = router + "*"
				} else {
					router = router + "/*"
				}
				logs.Debug("     初始化微信公众号验证文件路由：", router)
				web.Router(router, &WechatMpVerifyController{}, "get:Verify")
			}
			logs.Debug("结束 初始化微信公众号验证文件路由")
		}
	}
}

/**
 * @description: 验证微信公众号文件
 * @return {*}
 */
func (ctrl *WechatMpVerifyController) Verify() {
	if list, err := service.NewWeixinMpVerify().GetCache(); err != nil {
		ctrl.Ctx.WriteString(err.Error())
	} else {
		if len(list) > 0 {
			orpath := ctrl.Ctx.Request.URL.Path
			for _, item := range list {
				if item.Status != 2 {
					continue
				}
				router := item.Path
				if !strings.HasPrefix(router, "/") {
					router = "/" + router
				}
				if strings.HasSuffix(router, "/") {
					router = router + item.FileName
				} else {
					router = router + "/" + item.FileName
				}
				if strings.HasSuffix(orpath, router) {
					path := item.FilePath
					logs.Debug(orpath, path)
					ctrl.Ctx.Request.Header.Set("Content-Type", "application/txt")
					http.ServeFile(ctrl.Ctx.ResponseWriter, ctrl.Ctx.Request, path[1:])
					ctrl.StopRun()
				}
			}
		}
	}
	ctrl.Ctx.WriteString("not found")
}
