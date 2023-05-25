package biz

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xgeneric"
)

type CmsSite struct{}

func NewCmsSite() *CmsSite {
	return &CmsSite{}
}

func (this *CmsSite) SiteSaveSortId(siteId int64, sortId int32) error {
	site, siteDo := query.CmsSiteDo()
	_, err := siteDo.Where(site.SiteID.Eq(siteId)).UpdateColumns(
		map[string]interface{}{
			site.SortID.ColumnName().String():     sortId,
			site.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

func (this *CmsSite) SitePaginate(page, limit int, name string, title string) ([]*model.CmsSite, int64, error) {
	site, siteDo := query.CmsSiteDo()
	siteDo = siteDo.Where(site.IsDeleted.Is(false))
	if name != "" {
		siteDo = siteDo.Where(site.Name.Like("%" + name + "%"))
	}
	if title != "" {
		siteDo = siteDo.Where(site.Title.Like("%" + title + "%"))
	}
	return siteDo.Order(site.IsDefault.Desc(), site.SortID).FindByPage((page-1)*limit, limit)
}

func (this *CmsSite) SiteDelete(ids string) {
	if ids == "" {
		return
	}
	site, siteDo := query.CmsSiteDo()
	idarr := strings.Split(ids, ",")
	for i := 0; i < len(idarr); i++ {
		id := idarr[i]
		id64, _ := strconv.ParseInt(id, 0, 64)
		if id != "" {
			siteDo.Where(site.SiteID.Eq(id64)).Update(site.IsDeleted, true)
		}
	}
}

func (this *CmsSite) SiteOne(id int64) (*model.CmsSite, error) {
	site, siteDo := query.CmsSiteDo()
	return siteDo.Where(site.SiteID.Eq(id)).First()
}

func (this *CmsSite) SiteSave(mdl *model.CmsSite, domains []string, remarks []string) error {
	if mdl.Title == "" {
		return errors.New("站点名称不能为空")
	}
	// if mdl.DirPath == "" {
	// 	return errors.New("生成目录名不能为空")
	// }

	site, siteDo := query.CmsSiteDo()
	if mdl.IsDefault {
		if count, _ := siteDo.Where(site.IsDefault.Is(true), site.IsDeleted.Is(false), site.SiteID.Neq(mdl.SiteID)).Count(); count > 0 {
			return errors.New("默认站点只能有一个，请修改后重试")
		}
	}
	if count, _ := siteDo.Where(site.Name.Eq(mdl.Name), site.IsDeleted.Is(false), site.SiteID.Neq(mdl.SiteID)).Count(); count > 0 {
		return errors.New("网站名称已存在，请修改后重试")
	}
	if count, _ := siteDo.Where(site.Flag.Eq(mdl.Flag), site.IsDeleted.Is(false), site.SiteID.Neq(mdl.SiteID)).Count(); count > 0 {
		return errors.New("网站标识已存在，请修改后重试")
	}
	mdl.Logo2 = xgeneric.IFF(mdl.Logo1 == "", "", mdl.Logo2)
	mdl.Icon2 = xgeneric.IFF(mdl.Icon1 == "", "", mdl.Icon2)
	if mdl.SiteID <= 0 {
		if err := siteDo.Save(mdl); err != nil {
			return err
		}
	} else {
		// 修改 cms_admin_nav
		nav, navDo := query.CmsAdminNavDo()
		navDo.Where(nav.SiteID.Eq(mdl.SiteID), nav.Type.Eq("Site")).UpdateColumns(map[string]interface{}{
			nav.Title.ColumnName().String():      mdl.Title,
			nav.SortID.ColumnName().String():     mdl.SortID,
			nav.UpdateTime.ColumnName().String(): time.Now(),
		})
		if _, err := siteDo.Where(site.SiteID.Eq(mdl.SiteID)).UpdateColumns(map[string]interface{}{
			site.Name.ColumnName().String():            mdl.Name,
			site.Flag.ColumnName().String():            mdl.Flag,
			site.Title.ColumnName().String():           mdl.Title,
			site.Template.ColumnName().String():        mdl.Template,
			site.IsDefault.ColumnName().String():       mdl.IsDefault,
			site.IsMobile.ColumnName().String():        mdl.IsMobile,
			site.Logo1.ColumnName().String():           mdl.Logo1,
			site.Logo2.ColumnName().String():           mdl.Logo2,
			site.Icon1.ColumnName().String():           mdl.Icon1,
			site.Icon2.ColumnName().String():           mdl.Icon2,
			site.Company.ColumnName().String():         mdl.Company,
			site.Address.ColumnName().String():         mdl.Address,
			site.Telphone.ColumnName().String():        mdl.Telphone,
			site.Fax.ColumnName().String():             mdl.Fax,
			site.Email.ColumnName().String():           mdl.Email,
			site.Crod.ColumnName().String():            mdl.Crod,
			site.HomeTitle.ColumnName().String():       mdl.HomeTitle,
			site.Copyright.ColumnName().String():       mdl.Copyright,
			site.Statcode.ColumnName().String():        mdl.Statcode,
			site.Robots.ColumnName().String():          mdl.Robots,
			site.MetaKeyword.ColumnName().String():     mdl.MetaKeyword,
			site.MetaDescription.ColumnName().String(): mdl.MetaDescription,
			site.UpdateID.ColumnName().String():        mdl.UpdateID,
			site.UpdateName.ColumnName().String():      mdl.UpdateName,
			site.UpdateTime.ColumnName().String():      time.Now(),
		}); err != nil {
			return err
		}
	}
	domain, domainDo := query.CmsSiteDomainDo()
	domainDo.Where(domain.SiteID.Eq(mdl.SiteID)).Delete()
	if domains != nil {
		domainLen := len(domains)
		for i := 0; i < domainLen; i++ {
			if domains[i] == "" {
				continue
			}
			domain := &model.CmsSiteDomain{}
			domain.Domain = domains[i]
			domain.Remark = remarks[i]
			domain.SiteID = mdl.SiteID
			domainDo.Save(domain)
		}
	}
	return nil
}

/**
 * @description: 频道分页查询
 * @param {*} page 页码
 * @param {int} limit 每页条数
 * @param {int64} siteId 站点ID
 * @param {*} name 频道名称
 * @param {string} title 频道标题
 * @return {*}
 */
func (this *CmsSite) ChannelPaginate(page, limit int, siteId int64, name, title string) ([]*model.CmsSiteChannel, int64, error) {
	mdl, do := query.CmsSiteChannelDo()
	if siteId > 0 {
		do = do.Where(mdl.SiteID.Eq(siteId))
	}
	if name != "" {
		do = do.Where(mdl.Name.Like("%" + name + "%"))
	}
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// ChannelFind 获取
func (this *CmsSite) ChannelFind(channelId int64) (*model.CmsSiteChannel, error) {
	mdl, do := query.CmsSiteChannelDo()
	return do.Where(mdl.ChannelID.Eq(channelId)).First()
}

// ChannelSave 保存或更新
func (this *CmsSite) ChannelSave(input *model.CmsSiteChannel) error {
	mdl, do := query.CmsSiteChannelDo()
	if input.Name != "" {
		if count, _ := do.Where(mdl.ChannelID.Neq(input.ChannelID), mdl.Name.Eq(input.Name)).Count(); count > 0 {
			return errors.New("频道名称不能重复，确保唯一")
		}
	}
	{
		var classLayer int32
		do.Where(mdl.ChannelID.Eq(input.ParentID)).Pluck(mdl.ClassLayer, &classLayer)
		classLayer++
		input.ClassLayer = classLayer
	}
	var err error
	input.UpdateTime = time.Now()
	input.ImgUrl2 = xgeneric.IFF(input.ImgUrl1 == "", "", input.ImgUrl2)
	if input.ChannelID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.ChannelID.Eq(input.ChannelID)).Updates(map[string]interface{}{
			mdl.ChannelID.ColumnName().String():  input.ChannelID,
			mdl.SiteID.ColumnName().String():     input.SiteID,
			mdl.Name.ColumnName().String():       input.Name,
			mdl.Title.ColumnName().String():      input.Title,
			mdl.Kind.ColumnName().String():       input.Kind,
			mdl.ClassLayer.ColumnName().String(): input.ClassLayer,
			mdl.LinkURL.ColumnName().String():    input.LinkURL,
			mdl.ImgUrl1.ColumnName().String():    input.ImgUrl1,
			mdl.ImgUrl2.ColumnName().String():    input.ImgUrl2,
			mdl.IsComment.ColumnName().String():  input.IsComment,
			mdl.IsAlbum.ColumnName().String():    input.IsAlbum,
			mdl.IsAttach.ColumnName().String():   input.IsAttach,
			mdl.IsSpec.ColumnName().String():     input.IsSpec,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.Status.ColumnName().String():     input.Status,
			mdl.TmplChnl.ColumnName().String():   input.TmplChnl,
			mdl.TmplCat.ColumnName().String():    input.TmplCat,
			mdl.TmplLst.ColumnName().String():    input.TmplLst,
			mdl.TmplDtl.ColumnName().String():    input.TmplDtl,
			mdl.IsShow.ColumnName().String():     input.IsShow,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	if err == nil {
		this.ChannelNav(input)
	}
	return err
}

/**
 * @description: ChannelNav 导航
 * @param {*model.CmsSiteChannel} input 频道
 * @return {*}
 */
func (this *CmsSite) ChannelNav(input *model.CmsSiteChannel) error {
	dt, _ := time.Parse("2006-01-02 15:04:05", "2023-03-20 00:00:00")
	// cms_admin_nav
	var siteNavId int64
	navMdl, navDo := query.CmsAdminNavDo()
	navDo.Where(navMdl.SiteID.Eq(input.SiteID), navMdl.Type.Eq("Site")).Pluck(navMdl.NavID, &siteNavId)
	if siteNavId <= 0 {
		siteMdl, siteDo := query.CmsSiteDo()
		site, _ := siteDo.Where(siteMdl.SiteID.Eq(input.SiteID)).First()
		// 创建站点导航
		nav := &model.CmsAdminNav{
			SiteID:     input.SiteID,
			ParentID:   100000,
			Type:       "Site",
			Name:       fmt.Sprintf("site_%d", site.SiteID),
			Title:      site.Title,
			SubTitle:   "",
			SortID:     site.SortID,
			Action:     "Show",
			IconURL:    "fa fa-gears",
			IsHide:     0,
			IsSys:      1,
			CreateTime: dt,
			UpdateTime: dt,
		}
		navDo.Create(nav)
		siteNavId = nav.NavID
	}
	if siteNavId > 0 {
		var channelNavId int64
		// 创建频道导航
		navDo.Where(navMdl.ChannelID.Eq(input.ChannelID), navMdl.Type.Eq("Channel")).Pluck(navMdl.NavID, &channelNavId)
		if channelNavId <= 0 {
			nav := &model.CmsAdminNav{
				SiteID:     input.SiteID,
				ChannelID:  input.ChannelID,
				ParentID:   siteNavId,
				Type:       "Channel",
				Name:       fmt.Sprintf("channel_%d", input.ChannelID),
				Title:      input.Title,
				SubTitle:   "",
				SortID:     input.SortID,
				Action:     "Show",
				LinkURL:    "",
				IconURL:    "fa fa-navicon",
				IsHide:     0,
				IsSys:      1,
				CreateTime: dt,
				UpdateTime: dt,
			}
			navDo.Create(nav)
			channelNavId = nav.NavID
		} else {
			navDo.Where(navMdl.NavID.Eq(channelNavId)).Updates(map[string]interface{}{
				navMdl.Title.ColumnName().String():      input.Title,
				navMdl.SortID.ColumnName().String():     input.SortID,
				navMdl.UpdateTime.ColumnName().String(): input.UpdateTime,
			})
		}
		if channelNavId > 0 {
			if count, _ := navDo.Where(navMdl.ParentID.Eq(channelNavId)).Count(); count == 0 {
				// 创建文章导航
				navArticle := &model.CmsAdminNav{
					SiteID:     input.SiteID,
					ChannelID:  input.ChannelID,
					ParentID:   channelNavId,
					Type:       "Article",
					Name:       fmt.Sprintf("channel_%d_%s", input.ChannelID, "article"),
					Title:      "内容管理",
					SubTitle:   "内容管理",
					SortID:     1,
					Action:     "Show,View,Add,Edit,Delete,Audit",
					LinkURL:    fmt.Sprintf("/admin/article/index?channelId=%d", input.ChannelID),
					IconURL:    "fa fa-tachometer",
					IsHide:     0,
					IsSys:      1,
					CreateTime: dt,
					UpdateTime: dt,
				}
				navDo.Create(navArticle)
				navCategory := &model.CmsAdminNav{
					SiteID:     input.SiteID,
					ChannelID:  input.ChannelID,
					ParentID:   channelNavId,
					Type:       "Article",
					Name:       fmt.Sprintf("channel_%d_%s", input.ChannelID, "category"),
					Title:      "栏目管理",
					SubTitle:   "栏目管理",
					SortID:     2,
					Action:     "Show,View,Add,Edit,Delete",
					LinkURL:    fmt.Sprintf("/admin/article/category?channelId=%d", input.ChannelID),
					IconURL:    "fa fa-tachometer",
					IsHide:     0,
					IsSys:      1,
					CreateTime: dt,
					UpdateTime: dt,
				}
				navDo.Create(navCategory)
				navComment := &model.CmsAdminNav{
					SiteID:     input.SiteID,
					ChannelID:  input.ChannelID,
					ParentID:   channelNavId,
					Type:       "Article",
					Name:       fmt.Sprintf("channel_%d_%s", input.ChannelID, "comment"),
					Title:      "评论管理",
					SubTitle:   "评论管理",
					SortID:     3,
					Action:     "Show,View,Add,Edit,Delete,Audit",
					LinkURL:    fmt.Sprintf("/admin/article/comment?channelId=%d", input.ChannelID),
					IconURL:    "fa fa-tachometer",
					IsHide:     0,
					IsSys:      1,
					CreateTime: dt,
					UpdateTime: dt,
				}
				navDo.Create(navComment)
			}
		}
	}
	return nil
}

func (this *CmsSite) ChannelSaveSortId(channelId int64, sortId int32) error {
	mdl, do := query.CmsSiteChannelDo()
	_, err := do.Where(mdl.ChannelID.Eq(channelId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

/**
 * @description: ChannelDestory 删除
 * @param {int64} siteId 站点ID
 * @param {int64} channelId 频道ID
 * @return {*}
 */
func (this *CmsSite) ChannelDestory(siteId, channelId int64) error {
	mdl, do := query.CmsSiteChannelDo()
	if count, _ := do.Where(mdl.ParentID.Eq(channelId)).Count(); count > 0 {
		return errors.New("请先删除子分类")
	}
	// 删除频道
	if _, err := do.Where(mdl.ChannelID.Eq(channelId)).Delete(); err != nil {
		return err
	}
	return nil
}

/**
 * @description: ChannelTree 获取频道树
 * @param {*} siteId 站点ID
 * @param {int64} channelId 频道ID
 * @return {*}
 */
func (this *CmsSite) ChannelTree(siteId, channelId int64) ([]*bizmodel.TreeNode, error) {
	out := make([]*bizmodel.TreeNode, 0)
	mdl, do := query.CmsSiteChannelDo()
	list, err := do.Where(mdl.SiteID.Eq(siteId), mdl.ParentID.Eq(0)).Order(mdl.SortID).Find()
	if err != nil {
		return out, err
	}
	for _, item := range list {
		child := &bizmodel.TreeNode{
			Id:       item.ChannelID,
			Name:     item.Title,
			Open:     true,
			Checked:  item.ChannelID == channelId,
			Selected: item.ChannelID == channelId,
			Children: nil,
		}
		children, _ := this.ChannelTreeByParentId(item.ChannelID, channelId)
		if children != nil {
			child.Children = children
		}
		out = append(out, child)
	}
	return out, nil
}

/**
 * @description: ChannelTreeByParentId 获取子分类
 * @param {*} parentId 父级ID
 * @param {int64} channelId 频道ID
 * @return {*}
 */
func (this *CmsSite) ChannelTreeByParentId(parentId, channelId int64) ([]*bizmodel.TreeNode, error) {
	out := make([]*bizmodel.TreeNode, 0)
	mdl, do := query.CmsSiteChannelDo()
	list, err := do.Where(mdl.ParentID.Eq(parentId)).Order(mdl.SortID).Find()
	if err != nil {
		return out, err
	}
	for _, item := range list {
		child := &bizmodel.TreeNode{
			Id:       item.ChannelID,
			Name:     item.Title,
			Checked:  item.ChannelID == channelId,
			Selected: item.ChannelID == channelId,
			Children: nil,
		}
		children, _ := this.ChannelTreeByParentId(item.ChannelID, channelId)
		if children != nil {
			child.Children = children
		}
		out = append(out, child)
	}
	return out, nil
}
