package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

// ArticleAttach 附件管理
func (c *ArticleController) ArticleAttach() {
	articleId, _ := c.GetInt64("articleId")
	if articleId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	channelId, _ := c.GetInt64("channelId")
	c.Data["articleId"] = articleId
	c.Data["channelId"] = channelId
	c.display()
}

func (c *ArticleController) AttachPaginate() {
	// page, limit := c.GetPagingParameters()
	articleId, _ := c.GetInt64("articleId")
	list, count, err := biz.NewCmsArticle().AttachPaginate(1, 99999, articleId)
	if err != nil {
		logs.Error("AttachPaginate", err.Error())
	}
	c.JSONPageSuccess(list, count)
}

func (c *ArticleController) AttachSave() {
	mdl := model.CmsArticleAttach{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("AttachSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsArticle().AttachSave(&mdl); err != nil {
		logs.Error("AttachSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) AttachSaveShow() {
	mdls := vmodel.Article_AttachSaveShowModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AttachSaveShow", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AttachSaveShow", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls.AttachIds {
		if err := biz.NewCmsArticle().AttachSaveShow(mdl, mdls.Show); err != nil {
			logs.Error("AttachSaveShow", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) AttachSaveBatch() {
	mdls := []vmodel.Article_AttachSaveBatchdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AttachSaveBatch", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AttachSaveBatch", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsArticle()
	for _, mdl := range mdls {
		if err := service.AttachSaveInfo(mdl.AttachID, mdl.Title, mdl.Point, mdl.Click, mdl.SortID, mdl.Remark); err != nil {
			logs.Error("AttachSaveBatch", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) AttachDestory() {
	attachId, _ := c.GetInt64("attachId")
	if err := biz.NewCmsArticle().AttachDestory(attachId); err != nil {
		logs.Error("AttachDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}
