package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

type ArticleController struct{ BaseController }

func (c *ArticleController) Index() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	c.Data["channelId"] = channelId
	c.display()
}

func (c *ArticleController) Category() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	c.Data["channelId"] = channelId
	c.display()
}

func (c *ArticleController) CategoryFind() {
	channelId, _ := c.GetInt64("channelId")
	list, count, err := biz.NewCmsArticle().CategoryPagination(0, 99999, channelId, "", "")
	if err != nil {
		logs.Error("CategoryFind", err.Error())
	}
	c.JSONPaging(lib.CodeSuccess, "", list, count)
}

func (c *ArticleController) CategoryEdit() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	c.Data["channelId"] = channelId
	categoryId, _ := c.GetInt64("categoryId")
	parentId, _ := c.GetInt64("parentId")
	mdl, err := biz.NewCmsArticle().CategoryFind(categoryId)
	if err != nil {
		mdl = &model.CmsArticleCategory{
			ParentID:  parentId,
			ChannelID: channelId,
			SortID:    99,
		}
	}
	c.Data["mdl"] = mdl
	c.display()
}

func (c *ArticleController) CategorySave() {
	mdl := model.CmsArticleCategory{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsArticle().CategorySave(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) CategorySaveSortId() {
	mdls := []vmodel.Category_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("CategorySaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("CategorySaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsArticle().CategorySaveSortId(mdl.CategoryId, int32(mdl.SortId)); err != nil {
			logs.Error("CategorySaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) CategoryDestory() {
	categoryId, _ := c.GetInt64("categoryId")
	if err := biz.NewCmsArticle().CategoryDestory(categoryId); err != nil {
		logs.Error("CategoryDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}
