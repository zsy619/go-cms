package plugin

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// init 注册就业相关插件路由
// 路由前缀: /plugin
func init() {
	logs.Debug("开始注册就业插件路由")

	webNs := web.NewNamespace("/plugin",
		// 线上报名插件
		web.NSNamespace("/xsbm",
			web.NSRouter("/index", &XsbmController{}, "*:Index"),
			web.NSRouter("/save", &XsbmController{}, "*:Save"),
		),
		// 就业岗位插件
		web.NSNamespace("/job",
			web.NSRouter("/list", &JobController{}, "*:List"),
			web.NSRouter("/detail", &JobController{}, "*:Detail"),
		),
		// 企业管理插件
		web.NSNamespace("/company",
			web.NSRouter("/list", &CompanyController{}, "*:List"),
			web.NSRouter("/detail", &CompanyController{}, "*:Detail"),
		),
		// 空中宣讲会插件
		web.NSNamespace("/airkeynote",
			web.NSRouter("/list", &AirkeynoteController{}, "*:List"),
			web.NSRouter("/detail", &AirkeynoteController{}, "*:Detail"),
		),
		// 招聘会插件
		web.NSNamespace("/jobfair",
			web.NSRouter("/list", &JobfairController{}, "*:List"),
			web.NSRouter("/detail", &JobfairController{}, "*:Detail"),
			web.NSRouter("/date/list", &JobfairController{}, "*:DateList"),
		),
	)
	web.AddNamespace(webNs)

	logs.Info("就业插件路由注册完成")
}