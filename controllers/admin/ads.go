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

type AdsController struct{ BaseController }

// Index 广告管理
// @router /admin/ads/index [get]
func (c *AdsController) Index() {
	list, _, _ := biz.NewCmsAds().CategoryPaginate(1, 99999, -1, -1, "", "")
	c.Data["categoryList"] = list
	c.display()
}

// AdsEdit 广告编辑
// @router /admin/ads/adsEdit [get]
func (c *AdsController) AdsEdit() {
	adsId, _ := c.GetInt64("adsId")
	clone, _ := c.GetInt("clone")
	mdl, err := biz.NewCmsAds().AdsFind(adsId)
	if err != nil {
		mdl = &model.CmsAds{
			SortID:    99,
			BeginTime: time.Now(),
			EndTime:   time.Now().AddDate(0, 0, 7),
			Target:    "_blank",
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.AdsID = 0
	}
	c.Data["mdl"] = mdl
	list, _, _ := biz.NewCmsAds().CategoryPaginate(1, 99999, -1, -1, "", "")
	c.Data["categoryList"] = list
	c.display()
}

// AdsSave 广告保存
// @router /admin/ads/adsSave [post]
func (c *AdsController) AdsSave() {
	mdl := model.CmsAds{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("AdsSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsAds().AdsSave(&mdl); err != nil {
		logs.Error("AdsSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

// AdsSaveSortId 保存排序
// @router /admin/ad/AdsSaveSortId [post]
func (c *AdsController) AdsSaveSortId() {
	mdls := []vmodel.Ads_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AdsSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AdsSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsAds()
	for _, mdl := range mdls {
		if err := service.AdsSaveSortId(mdl.AdsId, int32(mdl.SortId)); err != nil {
			logs.Error("AdsSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

// AdsDestory 删除
// @router /admin/ads/AdsDestory [post]
func (c *AdsController) AdsDestory() {
	adsId, _ := c.GetInt64("adsId")
	if err := biz.NewCmsAds().AdsDestory(adsId); err != nil {
		logs.Error("AdsDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// AdsChangeStatus 更改状态
// @router /admin/ads/AdsChangeStatus [post]
func (c *AdsController) AdsChangeStatus() {
	var mdl vmodel.Ads_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("AdsChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, adsId := range mdl.AdsIds {
		if err := biz.NewCmsAds().AdsChangeStatus(adsId, mdl.Status); err != nil {
			logs.Error("AdsChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

// AdsPaginate 列表
// @router /admin/ads/AdsPaginate [get]
func (c *AdsController) AdsPaginate() {
	page, limit := c.GetPagingParameters()
	categoryId, _ := c.GetInt64("categoryId")
	status, _ := c.GetInt32("status")
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	list, count, _ := biz.NewCmsAds().AdsPaginate(page, limit, -1, -1, categoryId, title, callIndex, status)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

// Category 链接分类
// @router /admin/ads/category [get]
func (c *AdsController) Category() {
	c.display()
}

func (c *AdsController) CategoryEdit() {
	categoryId, _ := c.GetInt64("categoryId")
	mdl, err := biz.NewCmsAds().CategoryFind(categoryId)
	if err != nil {
		mdl = &model.CmsAdsCategory{
			SortID: 99,
		}
	}
	c.Data["mdl"] = mdl
	c.display()
}

func (c *AdsController) CategorySave() {
	mdl := model.CmsAdsCategory{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsAds().CategorySave(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *AdsController) CategorySaveSortId() {
	mdls := []vmodel.Category_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("CategorySaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("CategorySaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsAds().CategorySaveSortId(mdl.CategoryId, int32(mdl.SortId)); err != nil {
			logs.Error("CategorySaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *AdsController) CategoryDestory() {
	categoryId, _ := c.GetInt64("categoryId")
	if err := biz.NewCmsAds().CategoryDestory(categoryId); err != nil {
		logs.Error("CategoryDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// CategoryPaginate 列表
// @router /admin/ads/categorypaginate [get]
func (c *AdsController) CategoryPaginate() {
	page, limit := c.GetPagingParameters()
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	list, count, _ := biz.NewCmsAds().CategoryPaginate(page, limit, -1, -1, title, callIndex)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}
