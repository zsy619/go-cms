package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

// ArticleAlbum 相册管理
func (c *ArticleController) ArticleAlbum() {
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

func (c *ArticleController) AlbumPaginate() {
	// page, limit := c.GetPagingParameters()
	articleId, _ := c.GetInt64("articleId")
	list, count, err := biz.NewCmsArticle().AlbumPaginate(1, 99999, articleId)
	if err != nil {
		logs.Error("AlbumPaginate", err.Error())
	}
	c.JSONPageSuccess(list, count)
}

func (c *ArticleController) AlbumSave() {
	mdl := model.CmsArticleAlbum{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("AlbumSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsArticle().AlbumSave(&mdl); err != nil {
		logs.Error("AlbumSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) AlbumSaveShow() {
	mdls := vmodel.Article_AlbumSaveShowModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AlbumSaveShow", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AlbumSaveShow", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls.AlbumIds {
		if err := biz.NewCmsArticle().AlbumSaveShow(mdl, mdls.Show); err != nil {
			logs.Error("AlbumSaveShow", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) AlbumSaveBatch() {
	mdls := []vmodel.Article_AlbumSaveBatchdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AlbumSaveBatch", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AlbumSaveBatch", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsArticle()
	for _, mdl := range mdls {
		if err := service.AlbumSaveInfo(mdl.AlbumID, mdl.Title, mdl.LinkURL, mdl.Click, mdl.SortID, mdl.Remark); err != nil {
			logs.Error("AlbumSaveBatch", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) AlbumDestory() {
	albumId, _ := c.GetInt64("albumId")
	if err := biz.NewCmsArticle().AlbumDestory(albumId); err != nil {
		logs.Error("AlbumDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}
