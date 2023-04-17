package admin

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

type AdController struct{ BaseController }

// Index 广告管理
// @router /admin/ad/index [get]
func (c *AdController) Index() {
	list, _, _ := biz.NewCmsAd().CategoryPaginate(1, 99999, -1, -1, "", "")
	c.Data["categoryList"] = list
	c.display()
}

func (c *AdController) AdEdit() {
	adId, _ := c.GetInt64("adId")
	clone, _ := c.GetInt("clone")
	mdl, err := biz.NewCmsAd().AdFind(adId)
	if err != nil {
		mdl = &model.CmsAd{
			SortID:    99,
			BeginTime: time.Now(),
			EndTime:   time.Now().AddDate(0, 0, 7),
			Target:    "_blank",
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.AdID = 0
	}
	c.Data["mdl"] = mdl
	list, _, _ := biz.NewCmsAd().CategoryPaginate(1, 99999, -1, -1, "", "")
	c.Data["categoryList"] = list
	c.display()
}

func (c *AdController) AdSave() {
	mdl := model.CmsAd{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("AdSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsAd().AdSave(&mdl); err != nil {
		logs.Error("AdSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

// AdPaginate 保存排序
func (c *AdController) AdSaveSortId() {
	mdls := []vmodel.Ad_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AdSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AdSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsAd()
	for _, mdl := range mdls {
		if err := service.AdSaveSortId(mdl.AdId, int32(mdl.SortId)); err != nil {
			logs.Error("AdSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	// biz.NewApiAd().InitCache()
	c.JSONSuccess("保存成功", nil)
}

func (c *AdController) AdDestory() {
	adId, _ := c.GetInt64("adId")
	if err := biz.NewCmsAd().AdDestory(adId); err != nil {
		logs.Error("AdDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

func (c *AdController) AdChangeStatus() {
	var mdl vmodel.Ad_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("AdChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, Adid := range mdl.AdIds {
		if err := biz.NewCmsAd().AdChangeStatus(Adid, mdl.Status); err != nil {
			logs.Error("AdChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

// AdPaginate 列表
// @router /admin/ad/Adpaginate [get]
func (c *AdController) AdPaginate() {
	page, limit := c.GetPagingParameters()
	categoryId, _ := c.GetInt64("categoryId")
	status, _ := c.GetInt32("status")
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	list, count, _ := biz.NewCmsAd().AdPaginate(page, limit, -1, -1, categoryId, title, callIndex, status)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

// Category 链接分类
// @router /admin/Ad/category [get]
func (c *AdController) Category() {
	c.display()
}

func (c *AdController) CategoryEdit() {
	categoryId, _ := c.GetInt64("categoryId")
	mdl, err := biz.NewCmsAd().CategoryFind(categoryId)
	if err != nil {
		mdl = &model.CmsAdCategory{
			SortID: 99,
		}
	}
	c.Data["mdl"] = mdl
	c.display()
}

func (c *AdController) CategorySave() {
	mdl := model.CmsAdCategory{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsAd().CategorySave(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *AdController) CategorySaveSortId() {
	mdls := []vmodel.Category_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("CategorySaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("CategorySaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsAd().CategorySaveSortId(mdl.CategoryId, int32(mdl.SortId)); err != nil {
			logs.Error("CategorySaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *AdController) CategoryDestory() {
	categoryId, _ := c.GetInt64("categoryId")
	if err := biz.NewCmsAd().CategoryDestory(categoryId); err != nil {
		logs.Error("CategoryDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// CategoryPaginate 列表
// @router /admin/Ad/categorypaginate [get]
func (c *AdController) CategoryPaginate() {
	page, limit := c.GetPagingParameters()
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	list, count, _ := biz.NewCmsAd().CategoryPaginate(page, limit, -1, -1, title, callIndex)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}
