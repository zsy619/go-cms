package plugin

import "github.com/beego/beego/v2/server/web"

func init() {
	web.Router("/plugin/xsbm/index", &XsbmController{}, "*:Index")
	web.Router("/plugin/xsbm/save", &XsbmController{}, "*:Save")
}
