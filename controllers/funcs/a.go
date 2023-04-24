package funcs

import (
	"fmt"
	"html/template"
	"sync"

	"github.com/beego/beego/v2/server/web"
)

// funcs 模板函数
var (
	funcs map[string]any
	mutex sync.Mutex
)

/**
 * @description: 初始化
 * @return {*}
 */
func init() {
	fmt.Println("初始化模板函数-------->开始")

	funcs = make(map[string]any)
	{
		funcs["ConcatStr"] = ConcatStr
		funcs["concatstr"] = ConcatStr
		funcs["SubStr"] = SubStr
		funcs["substr"] = SubStr
		funcs["StrCheck"] = StrCheck
		funcs["strcheck"] = StrCheck
		funcs["UcWords"] = UcWords
		funcs["ucwords"] = UcWords
		funcs["ToUpper"] = ToUpper
		funcs["toupper"] = ToUpper

		funcs["Time2Str"] = Time2Str
		funcs["time2str"] = Time2Str

		funcs["ArticleNew"] = ArticleNew
		funcs["ArticleNew"] = ArticleNew
		funcs["ArticleNewExt"] = ArticleNewExt
		funcs["articlenewext"] = ArticleNewExt
		funcs["ArticleTop"] = ArticleTop
		funcs["articletop"] = ArticleTop
		funcs["ArticleTopExt"] = ArticleTopExt
		funcs["articletopext"] = ArticleTopExt

		funcs["AdsNewExt"] = AdsNewExt
		funcs["adsnewext"] = AdsNewExt
		funcs["AdsNew"] = AdsNew
		funcs["adsnew"] = AdsNew
		funcs["AdsTopExt"] = AdsTopExt
		funcs["adstopext"] = AdsTopExt
		funcs["AdsTop"] = AdsTop
		funcs["adstop"] = AdsTop

		funcs["TagNew"] = TagNew
		funcs["tagnew"] = TagNew
		funcs["TagNewExt"] = TagNewExt
		funcs["tagnewext"] = TagNewExt
		funcs["TagTop"] = TagTop
		funcs["tagtop"] = TagTop
		funcs["TagTopExt"] = TagTopExt
		funcs["tagtopext"] = TagTopExt
		funcs["TagArtilceTop"] = TagArtilceTop
		funcs["tagartilcetop"] = TagArtilceTop

		funcs["TopicNewExt"] = TopicNewExt
		funcs["topicnewext"] = TopicNewExt
		funcs["TopicNew"] = TopicNew
		funcs["topicnew"] = TopicNew
		funcs["TopicTopExt"] = TopicTopExt
		funcs["topictopext"] = TopicTopExt
		funcs["TopicTop"] = TopicTop
		funcs["topictop"] = TopicTop
		funcs["TopicArtilceTop"] = TopicArtilceTop
		funcs["topicartilcetop"] = TopicArtilceTop

		funcs["UnixTimeFormat"] = UnixTimeFormat
		funcs["unixtimeformat"] = UnixTimeFormat

		funcs["SizeFormat"] = SizeFormat
		funcs["sizeformat"] = SizeFormat

		funcs["TemplateTheme"] = TemplateTheme
		funcs["templatetheme"] = TemplateTheme
		funcs["TemplateView"] = TemplateView
		funcs["templateview"] = TemplateView
		funcs["UrlForView"] = UrlForView
		funcs["urlforview"] = UrlForView

		funcs["SiteDefault"] = SiteDefault
		funcs["sitedefault"] = SiteDefault
		funcs["SiteMenu"] = SiteMenu
		funcs["sitemenu"] = SiteMenu
		funcs["SiteMenuFlag"] = SiteMenuFlag
		funcs["sitemenuflag"] = SiteMenuFlag

		funcs["IIF"] = IIF
		funcs["iif"] = IIF

		funcs["UrlForArticle"] = UrlForArticle
		funcs["urlforarticle"] = UrlForArticle
		funcs["UrlForArticleExt"] = UrlForArticleExt
		funcs["urlforarticleext"] = UrlForArticleExt

		funcs["UrlForCategory"] = UrlForCategory
		funcs["urlforcategory"] = UrlForCategory
		funcs["UrlForCategoryExt"] = UrlForCategoryExt
		funcs["urlforcategoryext"] = UrlForCategoryExt

		funcs["UrlForChannel"] = UrlForChannel
		funcs["urlforchannel"] = UrlForChannel
	}

	// 注册模板函数
	for k, v := range funcs {
		web.AddFuncMap(k, v)
	}

	fmt.Println("初始化模板函数-------->结束")
}

/**
 * @description: 注册模板函数
 * @param {*template.Template} tmpl 模板对象
 * @return {*}
 */
func InitFuncs(tmpl *template.Template) {
	mutex.Lock()
	defer mutex.Unlock()
	for k, v := range funcs {
		tmpl = tmpl.Funcs(template.FuncMap{k: v})
	}
}
