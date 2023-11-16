package admin

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type AdsController struct{ BaseController }

// Index 广告管理
// @router /admin/ads/index [get]
func (ctrl *AdsController) Index() {
	siteList, categoryList, _ := biz.NewCmsAds().SiteCategoryGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["categoryList"] = categoryList
	ctrl.Data["siteList"] = siteList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("ads_index")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// AdsEdit 广告编辑
// @router /admin/ads/adsEdit [get]
func (ctrl *AdsController) AdsEdit() {
	adsId, _ := ctrl.GetInt64("adsId")
	clone, _ := ctrl.GetInt("clone")
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
	ctrl.Data["mdl"] = mdl
	_, categoryList, _ := biz.NewCmsAds().SiteCategoryGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["categoryList"] = categoryList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("ads_index")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// AdsSave 广告保存
// @router /admin/ads/adsSave [post]
func (ctrl *AdsController) AdsSave() {
	mdl := model.CmsAds{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("AdsSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.AdsID <= 0 {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	} else {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	}
	if err := biz.NewCmsAds().AdsSave(&mdl); err != nil {
		logs.Error("AdsSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// AdsSaveSortId 保存排序
// @router /admin/ad/AdsSaveSortId [post]
func (ctrl *AdsController) AdsSaveSortId() {
	mdls := []vmodel.Ads_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("AdsSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AdsSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	service := biz.NewCmsAds()
	for _, mdl := range mdls {
		if err := service.AdsSaveSortId(mdl.AdsId, int32(mdl.SortId)); err != nil {
			logs.Error("AdsSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// AdsDestory 删除
// @router /admin/ads/AdsDestory [post]
func (ctrl *AdsController) AdsDestory() {
	adsId, _ := ctrl.GetInt64("adsId")
	if err := biz.NewCmsAds().AdsDestory(adsId); err != nil {
		logs.Error("AdsDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// AdsChangeStatus 更改状态
// @router /admin/ads/AdsChangeStatus [post]
func (ctrl *AdsController) AdsChangeStatus() {
	var mdl vmodel.Ads_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("AdsChangeStatus", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, adsId := range mdl.AdsIds {
		if err := biz.NewCmsAds().AdsChangeStatus(adsId, mdl.Status); err != nil {
			logs.Error("AdsChangeStatus", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

// AdsPaginate 列表
// @router /admin/ads/AdsPaginate [get]
func (ctrl *AdsController) AdsPaginate() {
	page, limit := ctrl.GetPagingParameters()
	siteId, _ := ctrl.GetInt64("siteId")
	categoryId, _ := ctrl.GetInt64("categoryId")
	status, _ := ctrl.GetInt32("status")
	title := ctrl.GetSafeString("title")
	callIndex := ctrl.GetSafeString("callIndex")

	siteIds := biz.NewCmsAds().SiteIdsGet(siteId, GlobalRoleType, GlobalRoleId)
	list, count, _ := biz.NewCmsAds().AdsPaginate(page, limit, -1, categoryId, title, callIndex, status, siteIds...)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

// Category 链接分类
// @router /admin/ads/category [get]
func (ctrl *AdsController) Category() {
	siteList, _, _ := biz.NewCmsAds().SiteCategoryGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["siteList"] = siteList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("ads_category")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *AdsController) CategoryEdit() {
	categoryId, _ := ctrl.GetInt64("categoryId")
	mdl, err := biz.NewCmsAds().CategoryFind(categoryId)
	if err != nil {
		mdl = &model.CmsAdsCategory{
			SortID: 99,
		}
	}
	ctrl.Data["mdl"] = mdl
	siteList, _, _ := biz.NewCmsAds().SiteCategoryGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["siteList"] = siteList
	ctrl.display()
}

func (ctrl *AdsController) CategorySave() {
	mdl := model.CmsAdsCategory{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if err := biz.NewCmsAds().CategorySave(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *AdsController) CategorySaveSortId() {
	mdls := []vmodel.Category_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("CategorySaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("CategorySaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsAds().CategorySaveSortId(mdl.CategoryId, int32(mdl.SortId)); err != nil {
			logs.Error("CategorySaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *AdsController) CategoryDestory() {
	categoryId, _ := ctrl.GetInt64("categoryId")
	if err := biz.NewCmsAds().CategoryDestory(categoryId); err != nil {
		logs.Error("CategoryDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// CategoryPaginate 列表
// @router /admin/ads/categorypaginate [get]
func (ctrl *AdsController) CategoryPaginate() {
	page, limit := ctrl.GetPagingParameters()
	siteId, _ := ctrl.GetInt64("siteId")
	title := ctrl.GetSafeString("title")
	callIndex := ctrl.GetSafeString("callIndex")

	siteIds := biz.NewCmsAds().SiteIdsGet(siteId, GlobalRoleType, GlobalRoleId)
	categoryList, count, _ := biz.NewCmsAds().CategoryPaginate(page, limit, -1, title, callIndex, siteIds...)
	ctrl.JSONPage(lib.CodeSuccess, "", categoryList, count)
}
