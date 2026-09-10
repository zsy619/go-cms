package admin

import (
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

// SiteController 站点管理控制器
type SiteController struct{ BaseController }

// Index 站点管理首页
// @router /admin/site/index [get]
func (ctrl *SiteController) Index() {
	roleMap := ctrl.RolePowerGet("site_index")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// SiteData 获取站点列表数据
// @router /admin/site/data [get]
//
// query 参数:
//   - tree: 任意非空值 → 返回所有未删除站点(不分页),供前端 treeTable 渲染
//   - page/limit: 分页参数(tree=空时生效)
//   - name/title: 模糊搜索关键字
func (ctrl *SiteController) SiteData() {
	name := ctrl.GetStringTrim("name", "")
	title := ctrl.GetStringTrim("title", "")
	siteService := service.NewCmsSite()

	// 树形表格模式: 返回全部数据(不分页),前端 treeTable 自行渲染树
	if ctrl.GetStringTrim("tree", "") != "" {
		list, count, _ := siteService.SiteListForTree(name, title)
		ctrl.JSONPage(lib.CodeSuccess, "", list, count)
		return
	}

	// 普通模式: 分页
	page, _ := ctrl.GetInt("page")
	limit, _ := ctrl.GetInt("limit")
	list, count, _ := siteService.SitePaginate(page, limit, name, title)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

// Channel 站点栏目管理页面
// @router /admin/site/channel [get]
func (ctrl *SiteController) Channel() {
	siteService := service.NewCmsSite()
	list, _, _ := siteService.SitePaginate(1, 999999, "", "")
	ctrl.Data["siteList"] = list
	roleMap := ctrl.RolePowerGet("site_channel")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// Delete 删除站点
// @router /admin/site/delete [get]
func (ctrl *SiteController) Delete() {
	ids := ctrl.GetStringTrim("ids", "")
	if ids == "" {
		ctrl.JSONError("参数丢失")
		return
	}
	siteService := service.NewCmsSite()
	siteService.SiteDelete(ids)
	data := lib.NewJSONResponse(lib.CodeSuccess, "")
	ctrl.JSONData(data)
}

// SiteEdit 站点编辑页面
// @router /admin/site/edit [get]
//
// query 参数:
//   - id: 编辑时为站点ID,新增时为 0
//   - parent_id: 新增子站时默认选中的父级ID(可选)
func (ctrl *SiteController) SiteEdit() {
	id, _ := ctrl.GetInt64("id", 0)
	parentIDQuery, _ := ctrl.GetInt64("parent_id", 0)
	site := &domain.CmsSite{SortID: 99} // 默认值
	var list []*domain.CmsSiteDomain
	if id != 0 {
		siteService := service.NewCmsSite()
		domainService := service.NewCmsSiteDomainModel()
		if s, _ := siteService.SiteOne(id); s != nil {
			site = s
		}
		list = domainService.List(id)
	}
	if len(list) == 0 {
		list = append(list, &domain.CmsSiteDomain{})
	}

	// 新增子站: 用 query 参数设置默认父级
	if id == 0 && parentIDQuery > 0 && site.ParentID == 0 {
		site.ParentID = parentIDQuery
	}

	// 加载所有可用站点作为父级选项(排除自身及其所有子站,防止循环引用)
	allSites, _, _ := service.NewCmsSite().SiteListForTree("", "")
	siteList := make([]*domain.CmsSite, 0)
	if id == 0 {
		// 新增时,全部可用
		siteList = allSites
	} else {
		// 编辑时,排除自身及子站(防止父级选自己形成环)
		descendantIDs := collectDescendantIDs(allSites, id)
		descendantIDs[id] = true
		for _, s := range allSites {
			if !descendantIDs[s.SiteID] {
				siteList = append(siteList, s)
			}
		}
	}

	ctrl.Data["listSize"] = len(list) - 1
	ctrl.Data["domainList"] = list
	ctrl.Data["site"] = site
	ctrl.Data["siteList"] = siteList
	ctrl.display()
}

// collectDescendantIDs 递归收集 siteID 的所有后代ID(含间接子站),返回 map[id]bool
func collectDescendantIDs(all []*domain.CmsSite, rootID int64) map[int64]bool {
	out := make(map[int64]bool)
	for _, s := range all {
		if s.ParentID == rootID {
			out[s.SiteID] = true
			for k := range collectDescendantIDs(all, s.SiteID) {
				out[k] = true
			}
		}
	}
	return out
}

// Save 保存站点
// @router /admin/site/save [post]
func (ctrl *SiteController) Save() {
	siteService := service.NewCmsSite()
	adminService := service.NewCmsAdmin()

	site := domain.CmsSite{}
	result := lib.NewJSONResponse(lib.CodeSuccess, "保存成功")
	if err := ctrl.ParseForm(&site); err != nil {
		logs.Error("解析站点表单失败: %v", err)
		result.SetResult(lib.CodeParamError, err.Error())
		ctrl.JSONData(result)
		return
	}

	domains := ctrl.GetStrings("domain", nil)
	remarks := ctrl.GetStrings("remark", nil)
	if domains == nil {
		domains = []string{}
	}
	if remarks == nil {
		remarks = []string{}
	}

	user := adminService.OneByUserId(ctrl.IsLogin())
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
	err := siteService.SiteSave(&site, domains, remarks)
	if err != nil {
		result.SetResult(lib.CodeFatal, err.Error())
	}

	if site.IsDefault {
		lib.SiteCache.Reset()
		lib.SiteFindCache.Reset()
	}
	ctrl.JSONData(result)
}

// SiteSaveSortId 保存站点排序
// @router /admin/site/SiteSaveSortId [post]
func (ctrl *SiteController) SiteSaveSortId() {
	mdls := []vmodel.Site_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("SiteSaveSortId请求: %s", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("SiteSaveSortId解析失败: %v", err)
		ctrl.JSONError(err.Error())
		return
	}
	for _, mdl := range mdls {
		if err := service.NewCmsSite().SiteSaveSortId(mdl.SiteId, mdl.SortId); err != nil {
			logs.Error("SiteSaveSortId保存失败: siteId=%d, sortId=%d, error=%v", mdl.SiteId, mdl.SortId, err)
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// ChannelFind 获取频道列表
// @router /admin/site/channel/find [get]
func (ctrl *SiteController) ChannelFind() {
	siteId, _ := ctrl.GetInt64("siteId")
	list, count, err := service.NewCmsSite().ChannelPaginate(1, 99999, siteId, "", "")
	if err != nil {
		logs.Error("ChannelFind查询失败: %v", err)
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

// ChannelEdit 频道编辑页面
// @router /admin/site/channel/edit [get]
func (ctrl *SiteController) ChannelEdit() {
	siteService := service.NewCmsSite()
	channelId, _ := ctrl.GetInt64("channelId")
	parentId, _ := ctrl.GetInt64("parentId")
	siteId, _ := ctrl.GetInt64("siteId")
	mdl, err := siteService.ChannelFind(channelId)
	if err != nil {
		mdl = &domain.CmsSiteChannel{
			ParentID: parentId,
			SiteID:   siteId,
			SortID:   99,
		}
	}
	ctrl.Data["mdl"] = mdl
	roleMap := ctrl.RolePowerGet("site_channel")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// ChannelSave 保存频道
// @router /admin/site/channel/save [post]
func (ctrl *SiteController) ChannelSave() {
	mdl := domain.CmsSiteChannel{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("ChannelSave解析表单失败: %v", err)
		ctrl.JSONError(err.Error())
		return
	}
	if err := service.NewCmsSite().ChannelSave(&mdl); err != nil {
		logs.Error("ChannelSave保存失败: channelId=%d, error=%v", mdl.ChannelID, err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// ChannelSaveSortId 保存频道排序
// @router /admin/site/channel/sort [post]
func (ctrl *SiteController) ChannelSaveSortId() {
	mdls := []vmodel.Channel_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("ChannelSaveSortId请求: %s", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("ChannelSaveSortId解析失败: %v", err)
		ctrl.JSONError(err.Error())
		return
	}
	for _, mdl := range mdls {
		if err := service.NewCmsSite().ChannelSaveSortId(mdl.ChannelID, mdl.SortId); err != nil {
			logs.Error("ChannelSaveSortId保存失败: channelId=%d, sortId=%d, error=%v", mdl.ChannelID, mdl.SortId, err)
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// ChannelDestory 删除频道
// @router /admin/site/channel/destroy [post]
func (ctrl *SiteController) ChannelDestory() {
	channelId, _ := ctrl.GetInt64("channelId")
	siteId, _ := ctrl.GetInt64("siteId")
	if err := service.NewCmsSite().ChannelDestory(siteId, channelId); err != nil {
		logs.Error("ChannelDestory删除失败: siteId=%d, channelId=%d, error=%v", siteId, channelId, err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// ChannelTree 获取频道树结构
// @router /admin/site/channel/tree [get]
func (ctrl *SiteController) ChannelTree() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	parentId, _ := ctrl.GetInt64("parentId")
	tree, err := service.NewCmsSite().ChannelTree(channelId, parentId)
	if err != nil {
		logs.Error("ChannelTree查询失败: channelId=%d, parentId=%d, error=%v", channelId, parentId, err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.Data["json"] = tree
	if err := ctrl.ServeJSON(); err != nil {
		logs.Error("ServeJSON失败: error=%v", err)
	}
	ctrl.StopRun()
}
