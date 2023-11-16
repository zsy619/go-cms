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
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type ArticleController struct{ BaseController }

func (ctrl *ArticleController) Index() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	ctrl.Data["channelId"] = channelId
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_article")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *ArticleController) ArticlePaginate() {
	page, limit := ctrl.GetPagingParameters()
	channelId, _ := ctrl.GetInt64("channelId")
	categoryId, _ := ctrl.GetInt64("categoryId")
	fmt.Println("categoryId", categoryId, "channelId", channelId)
	title := ctrl.GetSafeString("title")
	callIndex := ctrl.GetSafeString("callIndex")
	status, _ := ctrl.GetInt32("status")
	list, count, err := biz.NewCmsArticle().ArticlePaginate(page, limit, channelId, categoryId, title, callIndex, status)
	if err != nil {
		logs.Error("ArticlePaginate", err.Error())
	}
	ctrl.JSONPageSuccess(list, count)
}

func (ctrl *ArticleController) ArticleEdit() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	articleId, _ := ctrl.GetInt64("articleId")
	clone, _ := ctrl.GetInt("clone")
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
	ctrl.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_article")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *ArticleController) ArticleSave() {
	mdl := model.CmsArticle{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("ArticleSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	propertyData := ctrl.GetSafeString("propertyData", "")
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
		ctrl.JSONError(err.Error())
		return
	}
	// 保存自定义属性
	if propertyData != "" {
		var propertyList []vmodel.Article_PropertySaveModel
		err := json.Unmarshal([]byte(propertyData), &propertyList)
		if err != nil {
			logs.Error("ArticleSave", err.Error())
			ctrl.JSONError(err.Error())
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
				ctrl.JSONError(err.Error())
				return
			}
		}

	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *ArticleController) ArticleClone() {
	articleId, _ := ctrl.GetInt64("articleId")
	if _, err := biz.NewCmsArticle().ArticleClone(articleId); err != nil {
		logs.Error("ArticleClone", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("复制成功", nil)
}

func (ctrl *ArticleController) ArticleDestory() {
	articleId, _ := ctrl.GetInt64("articleId")
	if err := biz.NewCmsArticle().ArticleDestory(articleId); err != nil {
		logs.Error("ArticleDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

func (ctrl *ArticleController) ArticleChangeStatus() {
	var mdl vmodel.Article_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("ArticleChangeStatus", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, articleId := range mdl.ArticleIds {
		if err := biz.NewCmsArticle().ArticleChangeStatus(articleId, mdl.Status); err != nil {
			logs.Error("ArticleChangeStatus", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

func (ctrl *ArticleController) ArticleSaveSortId() {
	mdls := []vmodel.Article_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("ArticleSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("ArticleSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsArticle().ArticleSaveSortId(mdl.ArticleId, int32(mdl.SortId)); err != nil {
			logs.Error("ArticleSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *ArticleController) Category() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	ctrl.Data["channelId"] = channelId
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_category")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *ArticleController) CategoryFind() {
	channelId, _ := ctrl.GetInt64("channelId")
	list, count, err := biz.NewCmsArticle().CategoryPaginate(1, 99999, channelId, "", "")
	if err != nil {
		logs.Error("CategoryFind", err.Error())
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

func (ctrl *ArticleController) CategoryEdit() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	ctrl.Data["channelId"] = channelId
	categoryId, _ := ctrl.GetInt64("categoryId")
	parentId, _ := ctrl.GetInt64("parentId")
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
	ctrl.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_category")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *ArticleController) CategorySave() {
	mdl := model.CmsArticleCategory{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if err := biz.NewCmsArticle().CategorySave(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *ArticleController) CategoryAutoUrl() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.JSONError("频道参数错误")
		return
	}
	if err := biz.NewCmsArticle().CategoryAutoUrl(channelId); err != nil {
		logs.Error("CategoryAutoUrl", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *ArticleController) CategorySaveSortId() {
	mdls := []vmodel.Category_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("CategorySaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("CategorySaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsArticle().CategorySaveSortId(mdl.CategoryId, int32(mdl.SortId)); err != nil {
			logs.Error("CategorySaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	// biz.Cache_ApiArticleCategoryFind = make(map[string][]map[string]interface{})
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *ArticleController) CategoryDestory() {
	categoryId, _ := ctrl.GetInt64("categoryId")
	if err := biz.NewCmsArticle().CategoryDestory(categoryId); err != nil {
		logs.Error("CategoryDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

func (ctrl *ArticleController) CategoryTree() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	categoryId, _ := ctrl.GetInt64("categoryId")
	tree, err := biz.NewCmsArticle().CategoryTree(channelId, categoryId)
	if err != nil {
		logs.Error("CategoryTree", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.Data["json"] = tree
	ctrl.ServeJSON()
	ctrl.StopRun()
}
