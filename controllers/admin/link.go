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

type LinkController struct{ BaseController }

// Index 链接管理
// @router /admin/link/index [get]
func (c *LinkController) Index() {
	list, _, _ := biz.NewCmsLink().CategoryPaginate(1, 99999, -1, -1, "", "")
	c.Data["categoryList"] = list
	c.display()
}

func (c *LinkController) LinkEdit() {
	linkId, _ := c.GetInt64("linkId")
	clone, _ := c.GetInt("clone")
	mdl, err := biz.NewCmsLink().LinkFind(linkId)
	if err != nil {
		mdl = &model.CmsLink{
			SortID: 99,
			Target: "_blank",
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.LinkID = 0
	}
	c.Data["mdl"] = mdl
	list, _, _ := biz.NewCmsLink().CategoryPaginate(1, 99999, -1, -1, "", "")
	c.Data["categoryList"] = list
	c.display()
}

func (c *LinkController) LinkSave() {
	mdl := model.CmsLink{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("LinkSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsLink().LinkSave(&mdl); err != nil {
		logs.Error("LinkSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *LinkController) LinkSaveSortId() {
	mdls := []vmodel.Link_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("LinkSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("LinkSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsLink()
	for _, mdl := range mdls {
		if err := service.LinkSaveSortId(mdl.LinkId, int32(mdl.SortId)); err != nil {
			logs.Error("LinkSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	biz.NewApiLink().InitCache()
	c.JSONSuccess("保存成功", nil)
}

func (c *LinkController) LinkDestory() {
	linkId, _ := c.GetInt64("linkId")
	if err := biz.NewCmsLink().LinkDestory(linkId); err != nil {
		logs.Error("LinkDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

func (c *LinkController) LinkChangeStatus() {
	var mdl vmodel.Link_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("LinkChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, linkid := range mdl.LinkIds {
		if err := biz.NewCmsLink().LinkChangeStatus(linkid, mdl.Status); err != nil {
			logs.Error("LinkChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

// LinkPaginate 列表
// @router /admin/link/linkpaginate [get]
func (c *LinkController) LinkPaginate() {
	page, limit := c.GetPagingParameters()
	categoryId, _ := c.GetInt64("categoryId")
	status, _ := c.GetInt32("status")
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	list, count, _ := biz.NewCmsLink().LinkPaginate(page, limit, -1, -1, categoryId, title, callIndex, status)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

// Category 链接分类
// @router /admin/link/category [get]
func (c *LinkController) Category() {
	c.display()
}

func (c *LinkController) CategoryEdit() {
	categoryId, _ := c.GetInt64("categoryId")
	mdl, err := biz.NewCmsLink().CategoryFind(categoryId)
	if err != nil {
		mdl = &model.CmsLinkCategory{
			SortID: 99,
		}
	}
	c.Data["mdl"] = mdl
	c.display()
}

func (c *LinkController) CategorySave() {
	mdl := model.CmsLinkCategory{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsLink().CategorySave(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *LinkController) CategorySaveSortId() {
	mdls := []vmodel.Category_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("CategorySaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("CategorySaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsLink().CategorySaveSortId(mdl.CategoryId, int32(mdl.SortId)); err != nil {
			logs.Error("CategorySaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *LinkController) CategoryDestory() {
	categoryId, _ := c.GetInt64("categoryId")
	if err := biz.NewCmsLink().CategoryDestory(categoryId); err != nil {
		logs.Error("CategoryDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// CategoryPaginate 列表
// @router /admin/link/categorypaginate [get]
func (c *LinkController) CategoryPaginate() {
	page, limit := c.GetPagingParameters()
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	list, count, _ := biz.NewCmsLink().CategoryPaginate(page, limit, -1, -1, title, callIndex)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}
