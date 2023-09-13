package plugin

import "github.com/beego/beego/v2/server/web"

func init() {
	web.Router("/plugin/xsbm/index", &XsbmController{}, "*:Index")
	web.Router("/plugin/xsbm/save", &XsbmController{}, "*:Save")

	web.Router("/plugin/job/list", &JobController{}, "*:List")
	web.Router("/plugin/job/detail", &JobController{}, "*:Detail")

	web.Router("/plugin/company/list", &CompanyController{}, "*:List")
	web.Router("/plugin/company/detail", &CompanyController{}, "*:Detail")

	web.Router("/plugin/airkeynote/list", &AirkeynoteController{}, "*:List")
	web.Router("/plugin/airkeynote/detail", &AirkeynoteController{}, "*:Detail")

	web.Router("/plugin/jobfair/list", &JobfairController{}, "*:List")
	web.Router("/plugin/jobfair/detail", &JobfairController{}, "*:Detail")
	web.Router("/plugin/jobfair/date/list", &JobfairController{}, "*:DateList")
}
