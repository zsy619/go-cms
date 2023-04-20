package funcs

import "github.com/beego/beego/v2/server/web"

func init() {
	web.AddFuncMap("ConcatStr", ConcatStr)
	web.AddFuncMap("concatstr", ConcatStr)
	web.AddFuncMap("SubStr", SubStr)
	web.AddFuncMap("substr", SubStr)
	web.AddFuncMap("StrCheck", StrCheck)
	web.AddFuncMap("strcheck", StrCheck)
	web.AddFuncMap("UcWords", UcWords)
	web.AddFuncMap("ucwords", UcWords)
	web.AddFuncMap("ToUpper", ToUpper)
	web.AddFuncMap("toupper", ToUpper)

	web.AddFuncMap("Time2Str", Time2Str)
	web.AddFuncMap("time2str", Time2Str)

	web.AddFuncMap("ArticleNew", ArticleNew)
	web.AddFuncMap("ArticleNew", ArticleNew)
	web.AddFuncMap("ArticleNewExt", ArticleNewExt)
	web.AddFuncMap("articlenewext", ArticleNewExt)
	web.AddFuncMap("ArticleTop", ArticleTop)
	web.AddFuncMap("articletop", ArticleTop)
	web.AddFuncMap("ArticleTopExt", ArticleTopExt)
	web.AddFuncMap("articletopext", ArticleTopExt)

	web.AddFuncMap("AdsNewExt", AdsNewExt)
	web.AddFuncMap("adsnewext", AdsNewExt)
	web.AddFuncMap("AdsNew", AdsNew)
	web.AddFuncMap("adsnew", AdsNew)
	web.AddFuncMap("AdsTopExt", AdsTopExt)
	web.AddFuncMap("adstopext", AdsTopExt)
	web.AddFuncMap("AdsTop", AdsTop)
	web.AddFuncMap("adstop", AdsTop)

	web.AddFuncMap("TagNew", TagNew)
	web.AddFuncMap("tagnew", TagNew)
	web.AddFuncMap("TagNewExt", TagNewExt)
	web.AddFuncMap("tagnewext", TagNewExt)
	web.AddFuncMap("TagTop", TagTop)
	web.AddFuncMap("tagtop", TagTop)
	web.AddFuncMap("TagTopExt", TagTopExt)
	web.AddFuncMap("tagtopext", TagTopExt)
	web.AddFuncMap("TagArtilceTop", TagArtilceTop)
	web.AddFuncMap("tagartilcetop", TagArtilceTop)

	web.AddFuncMap("TopicNewExt", TopicNewExt)
	web.AddFuncMap("topicnewext", TopicNewExt)
	web.AddFuncMap("TopicNew", TopicNew)
	web.AddFuncMap("topicnew", TopicNew)
	web.AddFuncMap("TopicTopExt", TopicTopExt)
	web.AddFuncMap("topictopext", TopicTopExt)
	web.AddFuncMap("TopicTop", TopicTop)
	web.AddFuncMap("topictop", TopicTop)
	web.AddFuncMap("TopicArtilceTop", TopicArtilceTop)
	web.AddFuncMap("topicartilcetop", TopicArtilceTop)

	web.AddFuncMap("UnixTimeFormat", UnixTimeFormat)
	web.AddFuncMap("unixtimeformat", UnixTimeFormat)

	web.AddFuncMap("SizeFormat", SizeFormat)
	web.AddFuncMap("sizeformat", SizeFormat)
}
