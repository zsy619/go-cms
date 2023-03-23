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

type AdminController struct{ BaseController }

func (c *AdminController) Log() {
	c.display()
}

func (c *AdminController) LogPaginate() {
	page, limit := c.GetPagingParameters()
	userName := c.GetString("userName")
	list, count, _ := biz.NewCmsAdmin().LogPaginate(page, limit, 0, userName)
	c.JSONPaging(lib.CodeSuccess, "", list, count)
}

func (c *AdminController) Role() {
	c.display()
}

func (c *AdminController) RolePaginate() {
	page, limit := c.GetPagingParameters()
	name := c.GetString("name")
	list, count, _ := biz.NewCmsAdmin().RolePaginate(page, limit, name)
	c.JSONPaging(lib.CodeSuccess, "", list, count)
}

func (c *AdminController) RoleEdit() {
	roleId, _ := c.GetInt64("roleId")
	mdl, err := biz.NewCmsAdmin().RoleFind(roleId)
	if err != nil {
		mdl = &model.CmsAdminRole{
			SortID: 99,
		}
	}
	c.Data["mdl"] = mdl
	c.display()
}

func (c *AdminController) RoleSave() {
	mdl := model.CmsAdminRole{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("RoleSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsAdmin().RoleSave(&mdl); err != nil {
		logs.Error("RoleSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *AdminController) RoleSaveSortId() {
	mdls := []vmodel.Role_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("RoleSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("RoleSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsAdmin().RoleSaveSortId(mdl.RoleId, int32(mdl.SortId)); err != nil {
			logs.Error("RoleSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *AdminController) RoleDestory() {
	roleId, _ := c.GetInt64("roleId")
	if err := biz.NewCmsAdmin().RoleDestory(roleId); err != nil {
		logs.Error("RoleDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}
