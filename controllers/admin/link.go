package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
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

// LinkPaginate 列表
// @router /admin/link/linkpaginate [get]
func (c *LinkController) LinkPaginate() {
	page, limit := c.GetPagingParameters()
	categoryId, _ := c.GetInt64("categoryId")
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	list, count, _ := biz.NewCmsLink().LinkPaginate(page, limit, -1, -1, categoryId, title, callIndex)
	c.JSONPaging(lib.CodeSuccess, "", list, count)
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

type CategorySaveSortIdModel struct {
	CategoryId int64 `json:"category_id"`
	SortId     int   `json:"sort_id"`
}

func (c *LinkController) CategorySaveSortId() {
	mdls := []CategorySaveSortIdModel{}
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
	c.JSONPaging(lib.CodeSuccess, "", list, count)
}
