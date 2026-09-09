package admin

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xjson"
	"golang.org/x/net/html"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

// ArticleController 文章管理控制器
type ArticleController struct{ BaseController }

// Index 文章列表页面
func (ctrl *ArticleController) Index() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	ctrl.Data["channelId"] = channelId
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_article")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// ArticlePaginate 获取文章分页列表
// @router /admin/article/ArticlePaginate [get]
func (ctrl *ArticleController) ArticlePaginate() {
	page, limit := ctrl.GetPagingParameters()
	channelId, _ := ctrl.GetInt64("channelId")
	categoryId, _ := ctrl.GetInt64("categoryId")
	title := ctrl.GetSafeString("title")
	callIndex := ctrl.GetSafeString("callIndex")
	status, _ := ctrl.GetInt32("status")

	logs.Debug("文章分页查询: channelId=%d, categoryId=%d, status=%d", channelId, categoryId, status)

	list, count, err := service.NewCmsArticle().ArticlePaginate(page, limit, channelId, categoryId, title, callIndex, status)
	if err != nil {
		logs.Error("ArticlePaginate查询失败: %v", err)
	}
	ctrl.JSONPageSuccess(list, count)
}

// ArticleEdit 文章编辑页面
func (ctrl *ArticleController) ArticleEdit() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	articleId, _ := ctrl.GetInt64("articleId")
	clone, _ := ctrl.GetInt("clone")
	mdl, err := service.NewCmsArticle().ArticleFind(articleId)
	if err != nil {
		mdl = &domain.CmsArticle{
			SortID:      99,
			Author:      GlobalAdminName,
			ChannelID:   channelId,
			PublishTime: time.Now(),
		}
	}
	if clone == 1 {
		mdl.ArticleID = 0
	}
	ctrl.Data["mdl"] = mdl
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_article")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// ArticleSave 保存文章
// @router /admin/article/ArticleSave [post]
func (ctrl *ArticleController) ArticleSave() {
	mdl := domain.CmsArticle{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("ArticleSave解析表单失败: %v", err)
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
		doc, err := html.Parse(strings.NewReader(mdl.Content))
		if err != nil {
			logs.Error("解析HTML内容失败: %v", err)
		} else {
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
			outStr = strings.ReplaceAll(outStr, "\t", "")
			outStr = strings.ReplaceAll(outStr, "\r\n", "")
			outStr = strings.ReplaceAll(outStr, "\r", "")
			outStr = strings.ReplaceAll(outStr, "\n", "")
			outStr = strings.ReplaceAll(outStr, "  ", "")
			outStr = strings.ReplaceAll(outStr, " ", "")
			logs.Debug("提取文章摘要: length=%d", len(outStr))
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
	if err := service.NewCmsArticle().ArticleSave(&mdl); err != nil {
		logs.Error("保存文章失败: articleId=%d, error=%v", mdl.ArticleID, err)
		ctrl.JSONError(err.Error())
		return
	}
	if propertyData != "" {
		var propertyList []vmodel.Article_PropertySaveModel
		err := json.Unmarshal([]byte(propertyData), &propertyList)
		if err != nil {
			logs.Error("解析属性数据失败: %v", err)
			ctrl.JSONError(err.Error())
		}
		for i := 0; i < len(propertyList); i++ {
			var item domain.CmsArticleProperty
			item.PropertyID = propertyList[i].PropertyID
			item.ParentID = propertyList[i].ParentID
			item.ArticleID = propertyList[i].ArticleID
			item.Title = propertyList[i].Title
			item.CallIndex = propertyList[i].CallIndex
			item.Value = propertyList[i].Value
			item.SortID = propertyList[i].SortID
			item.Status = propertyList[i].Status
			item.Deleted = propertyList[i].Deleted
			item.BelongTo = propertyList[i].BelongTo
			item.CreateID = propertyList[i].CreateID
			item.CreateName = propertyList[i].CreateName
			item.UpdateID = propertyList[i].UpdateID
			item.UpdateName = propertyList[i].UpdateName

			if err := service.NewCmsArticle().PropertySave(&item); err != nil {
				logs.Error("保存文章属性失败: propertyId=%d, error=%v", propertyList[i].PropertyID, err)
				ctrl.JSONError(err.Error())
				return
			}
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// ArticleClone 克隆文章
// @router /admin/article/ArticleClone [post]
func (ctrl *ArticleController) ArticleClone() {
	articleId, _ := ctrl.GetInt64("articleId")
	if _, err := service.NewCmsArticle().ArticleClone(articleId); err != nil {
		logs.Error("克隆文章失败: articleId=%d, error=%v", articleId, err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("复制成功", nil)
}

// ArticleDestory 删除文章
// @router /admin/article/ArticleDestory [post]
func (ctrl *ArticleController) ArticleDestory() {
	articleId, _ := ctrl.GetInt64("articleId")
	if err := service.NewCmsArticle().ArticleDestory(articleId); err != nil {
		logs.Error("删除文章失败: articleId=%d, error=%v", articleId, err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// ArticleChangeStatus 批量修改文章状态
// @router /admin/article/ArticleChangeStatus [post]
func (ctrl *ArticleController) ArticleChangeStatus() {
	var mdl vmodel.Article_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("解析状态修改请求失败: %v", err)
		ctrl.JSONError(err.Error())
		return
	}
	for _, articleId := range mdl.ArticleIds {
		if err := service.NewCmsArticle().ArticleChangeStatus(articleId, mdl.Status); err != nil {
			logs.Error("修改文章状态失败: articleId=%d, status=%d, error=%v", articleId, mdl.Status, err)
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

// ArticleSaveSortId 保存文章排序
// @router /admin/article/ArticleSaveSortId [post]
func (ctrl *ArticleController) ArticleSaveSortId() {
	mdls := []vmodel.Article_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("ArticleSaveSortId request: %s", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("ArticleSaveSortId解析失败: %v", err)
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := service.NewCmsArticle().ArticleSaveSortId(mdl.ArticleId, int32(mdl.SortId)); err != nil {
			logs.Error("保存文章排序失败: articleId=%d, sortId=%d, error=%v", mdl.ArticleId, mdl.SortId, err)
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// Category 栏目管理页面
func (ctrl *ArticleController) Category() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	ctrl.Data["channelId"] = channelId
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_category")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// CategoryFind 获取栏目列表
// @router /admin/article/CategoryFind [get]
func (ctrl *ArticleController) CategoryFind() {
	channelId, _ := ctrl.GetInt64("channelId")
	list, count, err := service.NewCmsArticle().CategoryPaginate(1, 99999, channelId, "", "")
	if err != nil {
		logs.Error("CategoryFind查询失败: %v", err)
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

// CategoryEdit 栏目编辑页面
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
	mdl, err := service.NewCmsArticle().CategoryFind(categoryId)
	if err != nil {
		mdl = &domain.CmsArticleCategory{
			IsShow:    true,
			IsSearch:  true,
			ParentID:  parentId,
			ChannelID: channelId,
			SortID:    99,
		}
	}
	ctrl.Data["mdl"] = mdl
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_category")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// CategorySave 保存栏目
// @router /admin/article/CategorySave [post]
func (ctrl *ArticleController) CategorySave() {
	mdl := domain.CmsArticleCategory{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("CategorySave解析表单失败: %v", err)
		ctrl.JSONError(err.Error())
	}
	if err := service.NewCmsArticle().CategorySave(&mdl); err != nil {
		logs.Error("保存栏目失败: categoryId=%d, error=%v", mdl.CategoryID, err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// CategoryAutoUrl 自动生成栏目URL
// @router /admin/article/CategoryAutoUrl [post]
func (ctrl *ArticleController) CategoryAutoUrl() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.JSONError("频道参数错误")
		return
	}
	if err := service.NewCmsArticle().CategoryAutoUrl(channelId); err != nil {
		logs.Error("自动生成栏目URL失败: channelId=%d, error=%v", channelId, err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// CategorySaveSortId 保存栏目排序
// @router /admin/article/CategorySaveSortId [post]
func (ctrl *ArticleController) CategorySaveSortId() {
	mdls := []vmodel.Category_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("CategorySaveSortId request: %s", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("CategorySaveSortId解析失败: %v", err)
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := service.NewCmsArticle().CategorySaveSortId(mdl.CategoryId, int32(mdl.SortId)); err != nil {
			logs.Error("保存栏目排序失败: categoryId=%d, sortId=%d, error=%v", mdl.CategoryId, mdl.SortId, err)
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// CategoryDestory 删除栏目
// @router /admin/article/CategoryDestory [post]
func (ctrl *ArticleController) CategoryDestory() {
	categoryId, _ := ctrl.GetInt64("categoryId")
	if err := service.NewCmsArticle().CategoryDestory(categoryId); err != nil {
		logs.Error("删除栏目失败: categoryId=%d, error=%v", categoryId, err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// CategoryTree 获取栏目树结构
// @router /admin/article/CategoryTree [get]
func (ctrl *ArticleController) CategoryTree() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	categoryId, _ := ctrl.GetInt64("categoryId")
	tree, err := service.NewCmsArticle().CategoryTree(channelId, categoryId)
	if err != nil {
		logs.Error("获取栏目树失败: channelId=%d, categoryId=%d, error=%v", channelId, categoryId, err)
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.Data["json"] = tree
	ctrl.ServeJSON()
	ctrl.StopRun()
}