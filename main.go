package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/beego/beego/v2/server/web/session"
	"github.com/kardianos/service"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers"
	_ "haedu.gov.cn/cms/controllers"
	_ "haedu.gov.cn/cms/controllers/admin"
	_ "haedu.gov.cn/cms/controllers/funcs"
	_ "haedu.gov.cn/cms/controllers/www"
)

var globalSessions *session.Manager

// Session初始化
func InitSession() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      7200,
		Maxlifetime:     7200,
		Secure:          false,
		CookieLifeTime:  7200,
		ProviderConfig:  "./session",
	}

	globalSessions, _ = session.NewManager("memory", sessionConfig)
	go globalSessions.GC()
}

func init() {
	// InsertFilter是提供一个过滤函数
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		// 允许访问所有源
		AllowAllOrigins: true,
		// 可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 指的是允许的Header的种类
		AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		// 公开的HTTP标头列表
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		// 如果设置，则允许共享身份验证凭据，例如cookie
		AllowCredentials: true,
	}))

	err := os.MkdirAll("./logs", os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = logs.SetLogger(logs.AdapterFile, `{"filename":"logs/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	if err != nil {
		panic(err)
	}
	logs.Info("hello cms!")
	// InitSession()
}

var serviceConfig = &service.Config{
	Name:        "cms",
	DisplayName: "内容管理",
	Description: "内容管理",
}

func main() {
	// 构建服务对象
	prog := &Program{}
	s, err := service.New(prog, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 用于记录系统日志
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		err = s.Run()
		if err != nil {
			logger.Error(err)
		}
		return
	}

	cmd := os.Args[1]

	if cmd == "install" {
		err = s.Install()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("安装成功")
	}
	if cmd == "uninstall" {
		err = s.Uninstall()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("卸载成功")
	}
}

type Program struct{}

func (p *Program) Start(s service.Service) error {
	log.Println("开始服务")
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	log.Println("停止服务")
	return nil
}

func (p *Program) run() { // 此处编写具体的服务代码
	// https://www.cnblogs.com/ahfuzhang/p/16745742.html
	// 能够减少GC的频率，从而提升程序性能
	{
		// ballast := make([]byte, 10*1024*1024*1024)
		// runtime.KeepAlive(ballast)
	}

	// 例: /images/user/1.jpg 实际访问的是 static/images/user/1.jpg
	web.SetStaticPath("/images", "static/images")
	web.SetStaticPath("/static/www/images", "static/www/images")
	// 通过 /css/资源路径  可以访问static/css目录的内容
	web.SetStaticPath("/css", "static/css")
	// 通过 /js/资源路径  可以访问static/js目录的内容
	web.SetStaticPath("/js", "static/js")

	web.SetStaticPath("/Uploads", "Uploads")
	web.SetStaticPath("/Uploads", "uploads")
	web.SetStaticPath("/Download", "Download")
	web.SetStaticPath("/Download", "download")
	web.SetStaticPath("/Downloads", "Downloads")
	web.SetStaticPath("/Downloads", "downloads")
	web.SetStaticPath("/Public", "Public")
	web.SetStaticPath("/debug", "debug")

	web.ErrorController(&controllers.ErrorController{})

	gob.Register(&model.CmsAdmin{})

	// https://beego.me/docs/mvc/controller/config.md
	web.BConfig.RouterCaseSensitive = false                    // 是否路由忽略大小写匹配，默认是 true，区分大小写
	web.BConfig.WebConfig.Session.SessionOn = true             // 开启Session模块
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 86400 // 设置Session有效期,单位秒
	web.BConfig.RecoverPanic = true                            // 是否异常恢复，默认值为 true，即当应用出现异常的情况，通过 recover 恢复回来，而不会导致应用异常退出。
	web.BConfig.CopyRequestBody = true                         // 是否允许在 HTTP 请求时，返回原始请求体数据字节，默认为 false （GET or HEAD or 上传文件请求除外）。
	web.BConfig.EnableErrorsShow = true                        // 是否显示系统错误信息，默认为 true。
	web.BConfig.EnableGzip = true                              // 是否开启 gzip 支持，默认为 false 不支持 gzip，一旦开启了 gzip，那么在模板输出的内容会进行 gzip 或者 zlib 压缩，根据用户的 Accept-Encoding 来判断。
	web.BConfig.MaxMemory = 1 << 26                            // 文件上传默认内存缓存大小，默认值是 1 << 26(64M)。
	web.BConfig.WebConfig.AutoRender = true                    // 是否模板自动渲染，默认值为 true，对于 API 类型的应用，应用需要把该选项设置为 false，不需要渲染模板。
	web.BConfig.WebConfig.EnableDocs = true                    // 是否开启文档内置功能，默认是 false
	web.BConfig.WebConfig.FlashName = "BEEGO_FLASH"            // Flash 数据设置时 Cookie 的名称，默认是 BEEGO_FLASH
	web.BConfig.WebConfig.DirectoryIndex = true                // 是否开启静态目录的列表显示，默认不显示目录，返回 403 错误。

	// beego过滤器 https://beego.me/docs/mvc/controller/filter.md
	// if global.AllowLogin {
	// admin过滤器
	web.InsertFilter("/admin/*", web.BeforeRouter, controllers.FilterAdmin)
	// web.InsertFilter("/mkt/*", web.BeforeRouter, routers.FilterSSO)
	// web.InsertFilter("/cms/*", web.BeforeRouter, routers.FilterSSO)
	// // mobile过滤器
	// web.InsertFilter("/mobile/x/*", web.BeforeRouter, routers.FilterMobile)
	// }

	web.Run()
}

func RefreshToken() {
	// oauth := mp.NewOAuth2(&models.MpConfig{
	// 	AppId:          global.WxMpConfig.AppId,
	// 	AppSecret:      global.WxMpConfig.AppSecret,
	// 	EncodingAesKey: global.WxMpConfig.EncodingAesKey,
	// 	Token:          global.WxMpConfig.Token,
	// })
	// oauth.RefreshToken()
}
