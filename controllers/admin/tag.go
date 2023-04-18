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

type TagController struct{ BaseController }

// Index 标签管理
// @router /admin/tag/index [get]
func (c *TagController) Index() {
	list, _, _ := biz.NewCmsSite().SitePaginate(1, 9999, "", "")
	c.Data["site"] = list
	c.display()
}

// TagEdit 编辑
// @router /admin/tag/TagEdit [get]
func (c *TagController) TagEdit() {
	tagId, _ := c.GetInt64("tagId")
	clone, _ := c.GetInt("clone")
	mdl, err := biz.NewCmsTag().TagFind(tagId)
	if err != nil {
		mdl = &model.CmsTag{
			SortID: 99,
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.TagID = 0
	}
	c.Data["mdl"] = mdl
	list, _, _ := biz.NewCmsSite().SitePaginate(1, 9999, "", "")
	c.Data["site"] = list
	c.display()
}

// tagSave 保存
// @router /admin/tag/tagSave [post]
func (c *TagController) TagSave() {
	mdl := model.CmsTag{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("TagSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsTag().TagSave(&mdl); err != nil {
		logs.Error("TagSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

// tagSaveSortId 保存排序
// @router /admin/tag/tagSaveSortId [post]
func (c *TagController) TagSaveSortId() {
	mdls := []vmodel.Tag_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("TagSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("TagSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsTag()
	for _, mdl := range mdls {
		if err := service.TagSaveSortId(mdl.TagId, int32(mdl.SortId)); err != nil {
			logs.Error("TagSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

// tagDestory 删除
// @router /admin/tag/tagDestory [post]
func (c *TagController) TagDestory() {
	tagId, _ := c.GetInt64("tagId")
	if err := biz.NewCmsTag().TagDestory(tagId); err != nil {
		logs.Error("TagDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// tagChangeStatus 更改状态
// @router /admin/tag/tagChangeStatus [post]
func (c *TagController) TagChangeStatus() {
	var mdl vmodel.Tag_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("TagChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, tagId := range mdl.TagIds {
		if err := biz.NewCmsTag().TagChangeStatus(tagId, mdl.Status); err != nil {
			logs.Error("TagChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

// tagPaginate 列表
// @router /admin/tag/tagPaginate [get]
func (c *TagController) TagPaginate() {
	page, limit := c.GetPagingParameters()
	siteId, _ := c.GetInt64("siteId")
	status, _ := c.GetInt32("status")
	title := c.GetString("title")
	name := c.GetString("name")
	list, count, _ := biz.NewCmsTag().TagPaginate(page, limit, siteId, -1, name, title, status)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}
