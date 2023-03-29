package www

import (
	"strconv"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/lib"
)

type ApiArticleController struct{ BaseController }

// @router /api/category/find [get]
func (this *ApiArticleController) CategoryFind() {
	channel_name := this.GetString("channel_name")
	if channel_name == "" {
		this.JSONErrorOfData("频道名称不能为空", nil)
	}
	outChannel, count, err := biz.NewApiArticle().CategoryFind(channel_name)
	if err != nil {
		this.JSONErrorOfData(err.Error(), outChannel)
	}
	this.JSONSuccess(strconv.FormatInt(count, 10), outChannel)
}

// @router /api/article/find [get]
func (this *ApiArticleController) Find() {
	limit, _ := this.GetInt("limit", 6)
	call_index := this.GetString("call_index")
	order_by := this.GetString("order_by", "sort_id")
	is_cache, _ := this.GetBool("is_cache", true)
	channel_id, _ := this.GetInt64("channel_id", 0)
	outArticle, count, err := biz.NewApiArticle().Find(limit, channel_id, call_index, order_by, is_cache)
	if err != nil {
		this.JSONErrorOfData(err.Error(), outArticle)
	}
	this.JSONSuccess(strconv.FormatInt(count, 10), outArticle)
}

// @router /api/article/paginate [get]
func (this *ApiArticleController) Paginate() {
	limit, _ := this.GetInt("limit", 6)
	page, _ := this.GetInt("page", 1)
	order_by := this.GetString("order_by", "sort_id")
	call_index := this.GetString("call_index")
	keyword := this.GetString("keyword")
	channel_id, _ := this.GetInt64("channel_id", 0)
	outArticle, count, err := biz.NewApiArticle().Paginate(page, limit, channel_id, call_index, keyword, order_by)
	if err != nil {
		this.JSONPage(lib.CodeError, err.Error(), outArticle, count)
	}
	this.JSONPageSuccess(outArticle, count)
}
