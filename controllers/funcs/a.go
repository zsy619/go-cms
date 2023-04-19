package funcs

import "github.com/beego/beego/v2/server/web"

func init() {
	web.AddFuncMap("ConcatStr", ConcatStr)
	web.AddFuncMap("concatstr", ConcatStr)
	web.AddFuncMap("SubStr", SubStr)
	web.AddFuncMap("substr", SubStr)
	web.AddFuncMap("StrCheck", StrCheck)
	web.AddFuncMap("strcheck", StrCheck)

	web.AddFuncMap("Time2Str", Time2Str)
	web.AddFuncMap("time2str", Time2Str)

	web.AddFuncMap("ArticleNew", ArticleNew)
	web.AddFuncMap("articlenew", ArticleNew)
	web.AddFuncMap("ArticleTop", ArticleTop)
	web.AddFuncMap("articletop", ArticleTop)

	web.AddFuncMap("TagNew", TagNew)
	web.AddFuncMap("tagnew", TagNew)

	web.AddFuncMap("TopicNew", TopicNew)
	web.AddFuncMap("topicnew", TopicNew)
}
