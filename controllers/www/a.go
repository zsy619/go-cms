package www

import "github.com/beego/beego/v2/server/web"

func init() {
	web.AutoPrefix("/", &IndexController{})
	web.AutoPrefix("/", &WeixinController{})
}
