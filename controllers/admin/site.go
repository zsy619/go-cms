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

type SiteController struct{ BaseController }

// Index 站点管理列表
// @router /admin/site/index [get]
func (ctrl *SiteController) Index() {
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("site_index")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// SiteData 站点列表数据
// @router /admin/site/data [get]
func (ctrl *SiteController) SiteData() {
	name := ctrl.GetString("name", "")
	title := ctrl.GetString("title", "")
	page, _ := ctrl.GetInt("page")
	limit, _ := ctrl.GetInt("limit")
	service := biz.NewCmsSite()
	list, count, _ := service.SitePaginate(page, limit, name, title)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

// Channel 站点栏目管理
// @router /admin/site/channel [get]
func (ctrl *SiteController) Channel() {
	service := biz.NewCmsSite()
	list, _, _ := service.SitePaginate(1, 999999, "", "")
	ctrl.Data["siteList"] = list
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("site_channel")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// SiteData 站点列表数据
// @router /admin/site/delete [get]
func (ctrl *SiteController) Delete() {
	ids := ctrl.GetString("ids", "")
	if ids == "" {
		ctrl.JSONError("参数丢失")
		return
	}
	service := biz.NewCmsSite()
	service.SiteDelete(ids)
	data := lib.NewJSONResponse(lib.CodeSuccess, "")
	ctrl.JSONData(data)
}

// SiteEdit 站点列表数据
// @router admin/site/edit [get]
func (ctrl *SiteController) SiteEdit() {
	id, _ := ctrl.GetInt64("id", 0)
	site := &model.CmsSite{}
	var list []*model.CmsSiteDomain
	if id != 0 {
		service := biz.NewCmsSite()
		domainService := biz.NewCmsSiteDomainModel()
		site, _ = service.SiteOne(id)
		list = domainService.List(id)
	} else {
		site = &model.CmsSite{
			SortID: 99,
		}
	}
	if len(list) == 0 {
		list = append(list, &model.CmsSiteDomain{})
	}
	ctrl.Data["listSize"] = len(list) - 1
	ctrl.Data["domainList"] = list
	ctrl.Data["site"] = site
	ctrl.display()
}

// SiteEdit 站点列表数据
// @router admin/site/save [post]
func (ctrl *SiteController) Save() {
	service := biz.NewCmsSite()
	adminSerice := biz.NewCmsAdmin()

	site := model.CmsSite{}
	result := lib.NewJSONResponse(lib.CodeSuccess, "保存成功")
	if err := ctrl.ParseForm(&site); err != nil {
		logs.Error("Save", err)
		result.SetResult(lib.CodeParamError, err.Error())
		ctrl.JSONData(result)
	}
	//if err := this.ParseForm(&domainList); err != nil {
	//	logs.Error("Save", err)
	//	result.SetResult(lib.CodeParamError, err.Error())
	//	this.JSONData(result)
	//}
	domains := ctrl.GetStrings("domain", nil)
	remarks := ctrl.GetStrings("remark", nil)
	if domains == nil {
		domains = []string{}
	}
	if remarks == nil {
		remarks = []string{}
	}

	user := adminSerice.OneByUserId(ctrl.IsLogin())
	if site.SiteID <= 0 {
		site.CreateTime = time.Now()
		site.UpdateTime = time.Now()
		site.CreateID = int32(ctrl.IsLogin())
		site.CreateName = user.RealName
	} else {
		site.UpdateTime = time.Now()
		site.UpdateID = int32(ctrl.IsLogin())
		site.UpdateName = user.RealName
	}
	err := service.SiteSave(&site, domains, remarks)
	if err != nil {
		result.SetResult(lib.CodeFatal, err.Error())
	}
	// 清除相关缓存
	if site.IsDefault {
		lib.SiteCache.Reset()
		lib.SiteFindCache.Reset()
	}
	ctrl.JSONData(result)
}

func (ctrl *SiteController) SiteSaveSortId() {
	mdls := []vmodel.Site_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("SiteSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("SiteSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsSite().SiteSaveSortId(mdl.SiteId, mdl.SortId); err != nil {
			logs.Error("SiteSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	// biz.Cache_ApiSiteDefault = nil
	// biz.Cache_ApiSiteGet = make(map[int64]*model.CmsSite, 0)
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *SiteController) ChannelFind() {
	siteId, _ := ctrl.GetInt64("siteId")
	list, count, err := biz.NewCmsSite().ChannelPaginate(1, 99999, siteId, "", "")
	if err != nil {
		logs.Error("ChannelFind", err.Error())
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

func (ctrl *SiteController) ChannelEdit() {
	service := biz.NewCmsSite()
	// list, _, _ := service.SitePaginate(1, 999999, "", "")
	// ctrl.Data["siteList"] = list
	channelId, _ := ctrl.GetInt64("channelId")
	parentId, _ := ctrl.GetInt64("parentId")
	siteId, _ := ctrl.GetInt64("siteId")
	mdl, err := service.ChannelFind(channelId)
	if err != nil {
		mdl = &model.CmsSiteChannel{
			ParentID: parentId,
			SiteID:   siteId,
			SortID:   99,
		}
	}
	ctrl.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("site_channel")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *SiteController) ChannelSave() {
	mdl := model.CmsSiteChannel{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("ChannelSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if err := biz.NewCmsSite().ChannelSave(&mdl); err != nil {
		logs.Error("ChannelSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *SiteController) ChannelSaveSortId() {
	mdls := []vmodel.Channel_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("ChannelSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("ChannelSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsSite().ChannelSaveSortId(mdl.ChannelID, mdl.SortId); err != nil {
			logs.Error("ChannelSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *SiteController) ChannelDestory() {
	channelId, _ := ctrl.GetInt64("channelId")
	siteId, _ := ctrl.GetInt64("siteId")
	if err := biz.NewCmsSite().ChannelDestory(siteId, channelId); err != nil {
		logs.Error("ChannelDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

func (ctrl *SiteController) ChannelTree() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	ChannelId, _ := ctrl.GetInt64("ChannelId")
	tree, err := biz.NewCmsSite().ChannelTree(channelId, ChannelId)
	if err != nil {
		logs.Error("ChannelTree", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.Data["json"] = tree
	ctrl.ServeJSON()
	ctrl.StopRun()
}
