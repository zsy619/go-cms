package plugin

import "github.com/beego/beego/v2/server/web"

func init() {
	webNs := web.NewNamespace("/plugin",
		web.NSNamespace("/xsbm",
			web.NSRouter("/index", &XsbmController{}, "*:Index"),
			web.NSRouter("/save", &XsbmController{}, "*:Save"),
		),
		web.NSNamespace("/job",
			web.NSRouter("/list", &JobController{}, "*:List"),
			web.NSRouter("/detail", &JobController{}, "*:Detail"),
		),
		web.NSNamespace("/company",
			web.NSRouter("/list", &CompanyController{}, "*:List"),
			web.NSRouter("/detail", &CompanyController{}, "*:Detail"),
		),
		web.NSNamespace("/airkeynote",
			web.NSRouter("/list", &AirkeynoteController{}, "*:List"),
			web.NSRouter("/detail", &AirkeynoteController{}, "*:Detail"),
		),
		web.NSNamespace("/jobfair",
			web.NSRouter("/list", &JobfairController{}, "*:List"),
			web.NSRouter("/detail", &JobfairController{}, "*:Detail"),
			web.NSRouter("/date/list", &JobfairController{}, "*:DateList"),
		),
	)
	web.AddNamespace(webNs)
}
