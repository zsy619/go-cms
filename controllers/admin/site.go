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

type SiteController struct{ BaseController }

// Index 站点管理列表
// @router /admin/site/index [get]
func (c *SiteController) Index() {
	// 获取角色权限
	roleMap := c.RolePowerGet("site_index")
	c.Data["roleMap"] = roleMap
	c.display()
}

// SiteData 站点列表数据
// @router /admin/site/data [get]
func (this *SiteController) SiteData() {
	name := this.GetString("name", "")
	title := this.GetString("title", "")
	page, _ := this.GetInt("page")
	limit, _ := this.GetInt("limit")
	service := biz.NewCmsSite()
	list, count, _ := service.SitePaginate(page, limit, name, title)
	this.JSONPage(lib.CodeSuccess, "", list, count)
}

// Channel 站点栏目管理
// @router /admin/site/channel [get]
func (c *SiteController) Channel() {
	service := biz.NewCmsSite()
	list, _, _ := service.SitePaginate(1, 999999, "", "")
	c.Data["siteList"] = list
	// 获取角色权限
	roleMap := c.RolePowerGet("site_channel")
	c.Data["roleMap"] = roleMap
	c.display()
}

// SiteData 站点列表数据
// @router /admin/site/delete [get]
func (this *SiteController) Delete() {
	ids := this.GetString("ids", "")
	if ids == "" {
		this.JSONError("参数丢失")
		return
	}
	service := biz.NewCmsSite()
	service.SiteDelete(ids)
	data := lib.NewJSONResponse(lib.CodeSuccess, "")
	this.JSONData(data)
}

// SiteEdit 站点列表数据
// @router admin/site/edit [get]
func (this *SiteController) SiteEdit() {
	id, _ := this.GetInt64("id", 0)
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
	if list == nil || len(list) == 0 {
		list = append(list, &model.CmsSiteDomain{})
	}
	this.Data["listSize"] = len(list) - 1
	this.Data["domainList"] = list
	this.Data["site"] = site
	this.display()
}

// SiteEdit 站点列表数据
// @router admin/site/save [post]
func (this *SiteController) Save() {
	service := biz.NewCmsSite()
	adminSerice := biz.NewCmsAdmin()

	site := model.CmsSite{}
	result := lib.NewJSONResponse(lib.CodeSuccess, "保存成功")
	if err := this.ParseForm(&site); err != nil {
		logs.Error("Save", err)
		result.SetResult(lib.CodeParamError, err.Error())
		this.JSONData(result)
	}
	//if err := this.ParseForm(&domainList); err != nil {
	//	logs.Error("Save", err)
	//	result.SetResult(lib.CodeParamError, err.Error())
	//	this.JSONData(result)
	//}
	domains := this.GetStrings("domain", nil)
	remarks := this.GetStrings("remark", nil)
	if domains == nil {
		domains = []string{}
	}
	if remarks == nil {
		remarks = []string{}
	}

	user := adminSerice.OneByUserId(this.IsLogin())
	if site.SiteID <= 0 {
		site.CreateTime = time.Now()
		site.CreateID = int32(this.IsLogin())
		site.CreateName = user.RealName
	} else {
		site.UpdateTime = time.Now()
		site.UpdateID = int32(this.IsLogin())
		site.UpdateName = user.RealName
	}
	err := service.SiteSave(&site, domains, remarks)
	if err != nil {
		result.SetResult(lib.CodeFatal, err.Error())
	}
	this.JSONData(result)
}

func (c *SiteController) SiteSaveSortId() {
	mdls := []vmodel.Site_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("SiteSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("SiteSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsSite().SiteSaveSortId(mdl.SiteId, mdl.SortId); err != nil {
			logs.Error("SiteSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	// biz.Cache_ApiSiteDefault = nil
	// biz.Cache_ApiSiteGet = make(map[int64]*model.CmsSite, 0)
	c.JSONSuccess("保存成功", nil)
}

func (c *SiteController) ChannelFind() {
	siteId, _ := c.GetInt64("siteId")
	list, count, err := biz.NewCmsSite().ChannelPaginate(1, 99999, siteId, "", "")
	if err != nil {
		logs.Error("ChannelFind", err.Error())
	}
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

func (c *SiteController) ChannelEdit() {
	service := biz.NewCmsSite()
	// list, _, _ := service.SitePaginate(1, 999999, "", "")
	// c.Data["siteList"] = list
	channelId, _ := c.GetInt64("channelId")
	parentId, _ := c.GetInt64("parentId")
	siteId, _ := c.GetInt64("siteId")
	mdl, err := service.ChannelFind(channelId)
	if err != nil {
		mdl = &model.CmsSiteChannel{
			ParentID: parentId,
			SiteID:   siteId,
			SortID:   99,
		}
	}
	c.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := c.RolePowerGet("site_channel")
	c.Data["roleMap"] = roleMap
	c.display()
}

func (c *SiteController) ChannelSave() {
	mdl := model.CmsSiteChannel{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("ChannelSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsSite().ChannelSave(&mdl); err != nil {
		logs.Error("ChannelSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *SiteController) ChannelSaveSortId() {
	mdls := []vmodel.Channel_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("ChannelSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("ChannelSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsSite().ChannelSaveSortId(mdl.ChannelID, mdl.SortId); err != nil {
			logs.Error("ChannelSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *SiteController) ChannelDestory() {
	channelId, _ := c.GetInt64("channelId")
	siteId, _ := c.GetInt64("siteId")
	if err := biz.NewCmsSite().ChannelDestory(siteId, channelId); err != nil {
		logs.Error("ChannelDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

func (c *SiteController) ChannelTree() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	ChannelId, _ := c.GetInt64("ChannelId")
	tree, err := biz.NewCmsSite().ChannelTree(channelId, ChannelId)
	if err != nil {
		logs.Error("ChannelTree", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.Data["json"] = tree
	c.ServeJSON()
	c.StopRun()
}
