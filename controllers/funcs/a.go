package funcs

import "github.com/beego/beego/v2/server/web"

func init() {
	web.AddFuncMap("ConcatStr", ConcatStr)
	web.AddFuncMap("concatstr", ConcatStr)
	web.AddFuncMap("SubStr", SubStr)
	web.AddFuncMap("substr", SubStr)
}
