package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"time"
)

type SiteController struct{ BaseController }

// Index 站点管理列表
// @router /admin/site/index [get]
func (c *SiteController) Index() {
	c.display()
}

// SiteData 站点列表数据
// @router /admin/site/data [get]
func (this *SiteController) SiteData() {
	name := this.GetString("name", "")
	title := this.GetString("title", "")
	domain := this.GetString("domain", "")
	page, _ := this.GetInt("page")
	limit, _ := this.GetInt("limit")
	service := biz.NewCmsSiteModel()
	list, count, _ := service.List(name, title, domain, page, limit)
	this.JSONPaging(lib.CodeSuccess, "", list, count)
}

// Channel 站点栏目管理
// @router /admin/site/channel [get]
func (c *SiteController) Channel() {
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
	service := biz.NewCmsSiteModel()
	service.Delete(ids)
	data := lib.NewJSONResponse(lib.CodeSuccess, "")
	this.JSONData(data)
}

// SiteEdit 站点列表数据
// @router admin/site/edit [get]
func (this *SiteController) SiteEdit() {
	id, _ := this.GetInt64("id", 0)
	site := &model.CmsSite{}
	if id != 0 {
		service := biz.NewCmsSiteModel()
		domainService := biz.NewCmsSiteDomainModel()
		site = service.One(id)
		list := domainService.List(id)

		if list == nil || len(list) == 0 {
			list = append(list, &model.CmsSiteDomain{})
		}
		this.Data["listSize"] = len(list) - 1
		this.Data["domainList"] = list
	}
	if site == nil {
		site = &model.CmsSite{}
	}
	this.Data["site"] = site
	this.display()
}

// SiteEdit 站点列表数据
// @router admin/site/save [post]
func (this *SiteController) Save() {
	service := biz.NewCmsSiteModel()
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
	err := service.Save(&site, domains, remarks)
	if err != nil {
		result.SetResult(lib.CodeFatal, err.Error())
	}
	this.JSONData(result)
}
