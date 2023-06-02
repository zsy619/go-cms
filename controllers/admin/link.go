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
	// 根据站点权限查询站点列表
	siteList, categoryList, _ := biz.NewCmsLink().SiteCategoryGet(GlobalRoleId, GlobalRoleType)
	c.Data["categoryList"] = categoryList
	c.Data["siteList"] = siteList
	c.display()
}

func (c *LinkController) LinkEdit() {
	siteId, _ := c.GetInt64("siteId")
	linkId, _ := c.GetInt64("linkId")
	clone, _ := c.GetInt("clone")
	mdl, err := biz.NewCmsLink().LinkFind(linkId)
	if err != nil {
		mdl = &model.CmsLink{
			SortID: 99,
			Target: "_blank",
			SiteID: siteId,
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.LinkID = 0
	}
	c.Data["mdl"] = mdl
	if GlobalRoleType == "super" {
		list, _, _ := biz.NewCmsLink().CategoryPaginate(1, 99999, -1, -1, "", "")
		c.Data["categoryList"] = list
	} else {
		siteIdList, _, _ := biz.NewCmsAdmin().RoleSiteFind(GlobalRoleId)
		categoryList := make([]interface{}, 0)
		if len(siteIdList) > 0 {
			for i := 0; i < len(siteIdList); i++ {
				list, _, _ := biz.NewCmsLink().CategoryPaginate(1, 99999, siteIdList[i].SiteID, -1, "", "")
				if len(list) > 0 {
					for j := 0; j < len(list); j++ {
						categoryList = append(categoryList, list[j])
					}
				}
			}
		}
		c.Data["categoryList"] = categoryList
	}

	c.display()
}

func (c *LinkController) LinkSave() {
	mdl := model.CmsLink{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("LinkSave", err.Error())
		c.JSONError(err.Error())
	}
	if mdl.LinkID <= 0 {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	} else {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	}
	if err := biz.NewCmsLink().LinkSave(&mdl); err != nil {
		logs.Error("LinkSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

// LinkPaginate 保存排序
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
	// biz.NewApiLink().InitCache()
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
	siteId, _ := c.GetInt64("siteId")
	categoryId, _ := c.GetInt64("categoryId")
	status, _ := c.GetInt32("status")
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	if GlobalRoleType == "super" || siteId != 0 {
		list, count, _ := biz.NewCmsLink().LinkPaginate(page, limit, siteId, -1, categoryId, title, callIndex, status)
		c.JSONPage(lib.CodeSuccess, "", list, count)
	} else {
		siteIdList, _, _ := biz.NewCmsAdmin().RoleSiteFind(GlobalRoleId)
		linkList := make([]interface{}, 0)
		var totalCount int64 = 0
		if len(siteIdList) > 0 {
			for i := 0; i < len(siteIdList); i++ {
				list, count, _ := biz.NewCmsLink().LinkPaginate(page, limit, siteIdList[i].SiteID, -1, categoryId, title, callIndex, status)
				totalCount += count
				if count > 0 {
					for j := 0; j < len(list); j++ {
						linkList = append(linkList, list[j])
					}
				}
			}
		}
		c.JSONPage(lib.CodeSuccess, "", linkList, totalCount)
	}
}

// Category 链接分类
// @router /admin/link/category [get]
func (c *LinkController) Category() {
	// 根据站点权限查询站点列表
	if GlobalRoleType == "super" {
		siteList, _, _ := biz.NewCmsSite().SitePaginate(1, 999999, "", "")
		c.Data["siteList"] = siteList
	} else {
		siteIdList, _, _ := biz.NewCmsAdmin().RoleSiteFind(GlobalRoleId)
		siteList := make([]interface{}, 0)
		if len(siteIdList) > 0 {
			for i := 0; i < len(siteIdList); i++ {
				item, _ := biz.NewCmsSite().SiteOne(siteIdList[i].SiteID)
				siteList = append(siteList, item)
			}
		}
		c.Data["siteList"] = siteList
	}
	c.display()
}

func (c *LinkController) CategoryEdit() {
	siteId, _ := c.GetInt64("siteId")
	categoryId, _ := c.GetInt64("categoryId")
	mdl, err := biz.NewCmsLink().CategoryFind(categoryId)
	if err != nil {
		mdl = &model.CmsLinkCategory{
			SortID: 99,
			SiteID: siteId,
		}
	}
	c.Data["mdl"] = mdl
	if GlobalRoleType == "super" {
		siteList, _, _ := biz.NewCmsSite().SitePaginate(1, 999999, "", "")
		c.Data["siteList"] = siteList
	} else {
		siteIdList, _, _ := biz.NewCmsAdmin().RoleSiteFind(GlobalRoleId)
		siteList := make([]interface{}, 0)
		if len(siteIdList) > 0 {
			for i := 0; i < len(siteIdList); i++ {
				item, _ := biz.NewCmsSite().SiteOne(siteIdList[i].SiteID)
				siteList = append(siteList, item)
			}
		}
		c.Data["siteList"] = siteList
	}
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
	siteId, _ := c.GetInt64("siteId")
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")

	if GlobalRoleType == "super" {
		list, count, _ := biz.NewCmsLink().CategoryPaginate(page, limit, siteId, -1, title, callIndex)
		c.JSONPage(lib.CodeSuccess, "", list, count)
	} else {
		siteIdList, _, _ := biz.NewCmsAdmin().RoleSiteFind(GlobalRoleId)
		categoryList := make([]interface{}, 0)
		var totalCount int64 = 0
		if len(siteIdList) > 0 {
			for i := 0; i < len(siteIdList); i++ {
				list, count, _ := biz.NewCmsLink().CategoryPaginate(1, 99999, siteIdList[i].SiteID, -1, "", "")
				totalCount += count
				if len(list) > 0 {
					for j := 0; j < len(list); j++ {
						categoryList = append(categoryList, list[j])
					}
				}
			}
		}
		c.JSONPage(lib.CodeSuccess, "", categoryList, totalCount)
	}
}
