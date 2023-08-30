package admin

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
	"strconv"
	"time"
)

func (c *ArticleController) PropertyEdit() {
	channelId, _ := c.GetInt64("channelId")
	propertyId, _ := c.GetInt64("propertyId")
	if propertyId <= 0 || channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	mdl, err := biz.NewCmsArticle().PropertyFind(propertyId)
	if err != nil {
		mdl = &model.CmsArticleProperty{
			SortID:     99,
			CreateName: GlobalAdminName,
			CreateTime: time.Now(),
		}
	}
	c.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := c.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_article")
	c.Data["roleMap"] = roleMap
	c.display()
}

func (c *ArticleController) PropertySave() {
	mdl := model.CmsArticleProperty{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("PropertySave", err.Error())
		c.JSONError(err.Error())
	}
	if mdl.PropertyID > 0 {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	} else {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	}

	if err := biz.NewCmsArticle().PropertySave(&mdl); err != nil {
		logs.Error("PropertySave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) PropertyPaginate() {
	page, limit := c.GetPagingParameters()
	limit = 99999
	parentId, _ := c.GetInt64("parentId")
	articleId, _ := c.GetInt64("articleId")
	status, _ := c.GetInt32("status")
	callIndex := c.GetString("callIndex")
	title := c.GetString("title")

	list, count, err := biz.NewCmsArticle().PropertyPaginate(page, limit, parentId, articleId, status, callIndex, title)
	if err != nil {
		logs.Error("PropertyPaginate", err.Error())
	}
	c.JSONPageSuccess(list, count)
}

func (c *ArticleController) PropertyDestroy() {
	propertyId, _ := c.GetInt64("propertyId")
	if err := biz.NewCmsArticle().PropertyDestroy(propertyId); err != nil {
		logs.Error("PropertyDestroy", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

func (c *ArticleController) PropertyChangeStatus() {
	var mdl vmodel.Property_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("PropertyChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, propertyId := range mdl.PropertyIds {
		if err := biz.NewCmsArticle().PropertyChangeStatus(propertyId, mdl.Status); err != nil {
			logs.Error("PropertyChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

func (c *ArticleController) PropertySaveSortId() {
	mdls := []vmodel.Property_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("PropertySaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("PropertySaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsArticle().PropertySaveSortId(mdl.PropertyId, int32(mdl.SortId)); err != nil {
			logs.Error("PropertySaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}
