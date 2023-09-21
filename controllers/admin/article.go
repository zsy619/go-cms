package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/net/html"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

type ArticleController struct{ BaseController }

func (c *ArticleController) Index() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	c.Data["channelId"] = channelId
	// 获取角色权限
	roleMap := c.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_article")
	c.Data["roleMap"] = roleMap
	c.display()
}

func (c *ArticleController) ArticlePaginate() {
	page, limit := c.GetPagingParameters()
	channelId, _ := c.GetInt64("channelId")
	categoryId, _ := c.GetInt64("categoryId")
	fmt.Println("categoryId", categoryId, "channelId", channelId)
	title := c.GetSafeString("title")
	callIndex := c.GetSafeString("callIndex")
	status, _ := c.GetInt32("status")
	list, count, err := biz.NewCmsArticle().ArticlePaginate(page, limit, channelId, categoryId, title, callIndex, status)
	if err != nil {
		logs.Error("ArticlePaginate", err.Error())
	}
	c.JSONPageSuccess(list, count)
}

func (c *ArticleController) ArticleEdit() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	articleId, _ := c.GetInt64("articleId")
	clone, _ := c.GetInt("clone")
	mdl, err := biz.NewCmsArticle().ArticleFind(articleId)
	if err != nil {
		mdl = &model.CmsArticle{
			SortID:      99,
			Author:      GlobalAdminName,
			ChannelID:   channelId,
			PublishTime: time.Now(),
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.ArticleID = 0
	}
	c.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := c.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_article")
	c.Data["roleMap"] = roleMap
	c.display()
}

func (c *ArticleController) ArticleSave() {
	mdl := model.CmsArticle{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("ArticleSave", err.Error())
		c.JSONError(err.Error())
	}
	propertyData := c.GetSafeString("propertyData", "")
	if mdl.ArticleID > 0 {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	} else {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	}
	if mdl.Summary == "" {
		// 解析 HTML 文本
		doc, err := html.Parse(strings.NewReader(mdl.Content))
		if err != nil {
			logs.Error(err)
		} else {
			// 拼接文本节点的内容
			var buf bytes.Buffer
			var traverse func(*html.Node)
			traverse = func(n *html.Node) {
				if n.Type == html.TextNode {
					buf.WriteString(n.Data)
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					traverse(c)
				}
			}
			traverse(doc)
			outStr := buf.String()
			// 替换制表符为4个空格
			outStr = strings.ReplaceAll(outStr, "\t", "")
			// 替换回车换行为换行符
			outStr = strings.ReplaceAll(outStr, "\r\n", "")
			outStr = strings.ReplaceAll(outStr, "\r", "")
			outStr = strings.ReplaceAll(outStr, "\n", "")
			// 替换多个连续空格为一个空格
			outStr = strings.ReplaceAll(outStr, "  ", "")
			outStr = strings.ReplaceAll(outStr, " ", "")
			fmt.Println(outStr)
			if len(outStr) > 200 {
				mdl.Summary = outStr[:255]
			} else {
				mdl.Summary = outStr[:]
			}
		}
	}
	if mdl.Tags != "" {
		mdl.Tags = strings.ReplaceAll(mdl.Tags, "，", ",")
	}
	if err := biz.NewCmsArticle().ArticleSave(&mdl); err != nil {
		logs.Error("ArticleSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	// 保存自定义属性
	if propertyData != "" {
		var propertyList []vmodel.Article_PropertySaveModel
		err := json.Unmarshal([]byte(propertyData), &propertyList)
		if err != nil {
			logs.Error("ArticleSave", err.Error())
			c.JSONError(err.Error())
		}
		for i := 0; i < len(propertyList); i++ {
			var item model.CmsArticleProperty
			item.PropertyID = propertyList[i].PropertyID
			item.ParentID = propertyList[i].ParentID
			item.ArticleID = propertyList[i].ArticleID
			item.Title = propertyList[i].Title
			item.CallIndex = propertyList[i].CallIndex
			item.Value = propertyList[i].Value
			item.SortID = propertyList[i].SortID
			item.Status = propertyList[i].Status
			item.IsDeleted = propertyList[i].IsDeleted
			item.BelongTo = propertyList[i].BelongTo
			item.CreateID = propertyList[i].CreateID
			item.CreateName = propertyList[i].CreateName
			item.UpdateID = propertyList[i].UpdateID
			item.UpdateName = propertyList[i].UpdateName

			if err := biz.NewCmsArticle().PropertySave(&item); err != nil {
				logs.Error("PropertySave", err.Error())
				c.JSONError(err.Error())
				return
			}
		}

	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) ArticleClone() {
	articleId, _ := c.GetInt64("articleId")
	if _, err := biz.NewCmsArticle().ArticleClone(articleId); err != nil {
		logs.Error("ArticleClone", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("复制成功", nil)
}

func (c *ArticleController) ArticleDestory() {
	articleId, _ := c.GetInt64("articleId")
	if err := biz.NewCmsArticle().ArticleDestory(articleId); err != nil {
		logs.Error("ArticleDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

func (c *ArticleController) ArticleChangeStatus() {
	var mdl vmodel.Article_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("ArticleChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, articleId := range mdl.ArticleIds {
		if err := biz.NewCmsArticle().ArticleChangeStatus(articleId, mdl.Status); err != nil {
			logs.Error("ArticleChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

func (c *ArticleController) ArticleSaveSortId() {
	mdls := []vmodel.Article_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("ArticleSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("ArticleSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsArticle().ArticleSaveSortId(mdl.ArticleId, int32(mdl.SortId)); err != nil {
			logs.Error("ArticleSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) Category() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	c.Data["channelId"] = channelId
	// 获取角色权限
	roleMap := c.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_category")
	c.Data["roleMap"] = roleMap
	c.display()
}

func (c *ArticleController) CategoryFind() {
	channelId, _ := c.GetInt64("channelId")
	list, count, err := biz.NewCmsArticle().CategoryPaginate(1, 99999, channelId, "", "")
	if err != nil {
		logs.Error("CategoryFind", err.Error())
	}
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

func (c *ArticleController) CategoryEdit() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	c.Data["channelId"] = channelId
	categoryId, _ := c.GetInt64("categoryId")
	parentId, _ := c.GetInt64("parentId")
	mdl, err := biz.NewCmsArticle().CategoryFind(categoryId)
	if err != nil {
		mdl = &model.CmsArticleCategory{
			IsShow:    true,
			IsSearch:  true,
			ParentID:  parentId,
			ChannelID: channelId,
			SortID:    99,
		}
	}
	c.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := c.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_category")
	c.Data["roleMap"] = roleMap
	c.display()
}

func (c *ArticleController) CategorySave() {
	mdl := model.CmsArticleCategory{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsArticle().CategorySave(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) CategoryAutoUrl() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.JSONError("频道参数错误")
		return
	}
	if err := biz.NewCmsArticle().CategoryAutoUrl(channelId); err != nil {
		logs.Error("CategoryAutoUrl", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) CategorySaveSortId() {
	mdls := []vmodel.Category_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("CategorySaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("CategorySaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsArticle().CategorySaveSortId(mdl.CategoryId, int32(mdl.SortId)); err != nil {
			logs.Error("CategorySaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	// biz.Cache_ApiArticleCategoryFind = make(map[string][]map[string]interface{})
	c.JSONSuccess("保存成功", nil)
}

func (c *ArticleController) CategoryDestory() {
	categoryId, _ := c.GetInt64("categoryId")
	if err := biz.NewCmsArticle().CategoryDestory(categoryId); err != nil {
		logs.Error("CategoryDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

func (c *ArticleController) CategoryTree() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	categoryId, _ := c.GetInt64("categoryId")
	tree, err := biz.NewCmsArticle().CategoryTree(channelId, categoryId)
	if err != nil {
		logs.Error("CategoryTree", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.Data["json"] = tree
	c.ServeJSON()
	c.StopRun()
}
