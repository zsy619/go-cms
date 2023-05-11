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

		funcs["IIF"] = IIF
		funcs["iif"] = IIF

		funcs["ArticleNew"] = ArticleNew
		funcs["ArticleNew"] = ArticleNew
		funcs["ArticleNewExtend"] = ArticleNewExtend
		funcs["articlenewextend"] = ArticleNewExtend
		funcs["article_new_extend"] = ArticleNewExtend
		funcs["ArticleTop"] = ArticleTop
		funcs["articletop"] = ArticleTop
		funcs["ArticleTopExtend"] = ArticleTopExtend
		funcs["articletopextend"] = ArticleTopExtend
		funcs["article_top_extend"] = ArticleTopExtend

		funcs["AdsNewExtend"] = AdsNewExtend
		funcs["adsnewextend"] = AdsNewExtend
		funcs["ads_new_extend"] = AdsNewExtend
		funcs["AdsNew"] = AdsNew
		funcs["adsnew"] = AdsNew
		funcs["AdsTopExtend"] = AdsTopExtend
		funcs["adstopextend"] = AdsTopExtend
		funcs["ads_top_extend"] = AdsTopExtend
		funcs["AdsTop"] = AdsTop
		funcs["adstop"] = AdsTop

		funcs["LinkNewExtend"] = LinkNewExtend
		funcs["linknewextend"] = LinkNewExtend
		funcs["LinkNew"] = LinkNew
		funcs["kinknew"] = LinkNew
		funcs["LinkTopExtend"] = LinkTopExtend
		funcs["linktopextend"] = LinkTopExtend
		funcs["LinkTop"] = LinkTop
		funcs["linktop"] = LinkTop

		funcs["TagNew"] = TagNew
		funcs["tagnew"] = TagNew
		funcs["TagNewExtend"] = TagNewExtend
		funcs["tagnewextend"] = TagNewExtend
		funcs["tag_new_extend"] = TagNewExtend
		funcs["TagTop"] = TagTop
		funcs["tagtop"] = TagTop
		funcs["TagTopExtend"] = TagTopExtend
		funcs["tagtopextend"] = TagTopExtend
		funcs["tag_top_extend"] = TagTopExtend
		funcs["TagArtilceTop"] = TagArtilceTop
		funcs["tagartilcetop"] = TagArtilceTop

		funcs["TopicNewExtend"] = TopicNewExtend
		funcs["topicnewextend"] = TopicNewExtend
		funcs["topic_new_extend"] = TopicNewExtend
		funcs["TopicNew"] = TopicNew
		funcs["topicnew"] = TopicNew
		funcs["TopicTopExtend"] = TopicTopExtend
		funcs["topictopextend"] = TopicTopExtend
		funcs["topic_top_extend"] = TopicTopExtend
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

		funcs["UrlForTopic"] = UrlForTopic
		funcs["urlfortopic"] = UrlForTopic

		funcs["UrlForTag"] = UrlForTag
		funcs["urlfortag"] = UrlForTag

		funcs["UrlForSearch"] = UrlForSearch
		funcs["urlforsearch"] = UrlForSearch
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
