package www

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"haedu.gov.cn/cms/app/biz"
)

type WechatMpVerifyController struct{ web.Controller }

/**
 * @description: 初始化微信公众号验证文件路由
 * @return {*}
 */
func InitWechatMpVerifyRouter() {
	if list, err := biz.NewWeixinMpVerify().GetCache(); err != nil {
		fmt.Println("InitMpVerifyRouter--->", err.Error())
	} else {
		if len(list) > 0 {
			fmt.Println("开始 初始化微信公众号验证文件路由")
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
				fmt.Println("     初始化微信公众号验证文件路由：", router)
				web.Router(router, &WechatMpVerifyController{}, "get:Verify")
			}
			fmt.Println("结束 初始化微信公众号验证文件路由")
		}
	}
}

/**
 * @description: 验证微信公众号文件
 * @return {*}
 */
func (c *WechatMpVerifyController) Verify() {
	if list, err := biz.NewWeixinMpVerify().GetCache(); err != nil {
		c.Ctx.WriteString(err.Error())
	} else {
		if len(list) > 0 {
			orpath := c.Ctx.Request.URL.Path
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
					fmt.Println(orpath, path)
					c.Ctx.Request.Header.Set("Content-Type", "application/txt")
					http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, path[1:])
					c.StopRun()
				}
			}
		}
	}
	c.Ctx.WriteString("not found")
}
