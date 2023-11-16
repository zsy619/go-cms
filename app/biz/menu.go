package biz

import (
	"fmt"

	linq "github.com/ahmetb/go-linq/v3"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type HomeInfo struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

type LogoInfo struct {
	Title string `json:"title"`
	Image string `json:"image"`
	Href  string `json:"href"`
}

type MenuInfo struct {
	Title  string      `json:"title"`
	Icon   string      `json:"icon"`
	Href   string      `json:"href"`
	Target string      `json:"target"`
	Child  []*MenuInfo `json:"child"`
}

type MenuOuter struct {
	HomeInfo HomeInfo    `json:"homeInfo"`
	LogoInfo LogoInfo    `json:"logoInfo"`
	MenuInfo []*MenuInfo `json:"menuInfo"`
}

type Menu struct{}

func NewMenu() *Menu {
	return &Menu{}
}

/**
 * @description: MenuList 获取菜单列表
 * @param {int64} adminId 管理员ID
 * @return {*}
 */
func (svc *Menu) MenuList(adminId int64) *MenuOuter {
	outerMenu := &MenuOuter{
		HomeInfo: HomeInfo{
			Title: "首页",
			Href:  "/admin/index/welcome",
		},
		LogoInfo: LogoInfo{
			Title: "CMS管理系统",
			Image: "/static/admin/images/logo.png",
			Href:  "",
		},
		MenuInfo: []*MenuInfo{},
	}

	// 角色获取
	sqlRole := fmt.Sprintf("SELECT a.role_id,a.type,a.is_sys FROM cms_admin_role a LEFT JOIN cms_admin b ON a.role_id = b.role_id WHERE b.user_id = %d", adminId)
	_, adminRoleDo := query.CmsAdminRoleDo()
	role := model.CmsAdminRole{}
	if err := adminRoleDo.UnderlyingDB().Raw(sqlRole).Scan(&role).Error; err != nil {
		return outerMenu
	}
	if role.RoleID == 0 {
		return outerMenu
	}
	// 是否为超级管理员
	hasSupper := role.Type == "supper" && role.IsSys == 1
	sqlNav := ""
	if hasSupper {
		sqlNav = "SELECT a.* FROM cms_admin_nav a WHERE a.is_hide = 0 ORDER BY a.parent_id ASC,a.sort_id ASC"
	} else {
		sqlNav = fmt.Sprintf(`SELECT a.nav_id,a.parent_id,a.site_id,a.channel_id,a.type,a.name,a.title,a.sub_title,a.icon_url,a.link_url,a.sort_id,b.action FROM cms_admin_nav a
LEFT JOIN cms_admin_role_value b ON a.name=b.nav_name
WHERE a.is_hide=0 AND b.role_id=%d
ORDER BY a.sort_id`, role.RoleID)
	}
	fmt.Println("sqlNav---------->", sqlNav)
	_, adminNavDo := query.CmsAdminNavDo()
	navs := []*model.CmsAdminNav{}
	if err := adminNavDo.UnderlyingDB().Raw(sqlNav).Scan(&navs).Error; err != nil {
		return outerMenu
	}

	parentNavs := []*model.CmsAdminNav{}
	linq.From(navs).WhereT(func(s *model.CmsAdminNav) bool {
		return s.ParentID == 0 && s.IsHide == 0
	}).OrderByT(func(s *model.CmsAdminNav) int32 {
		return s.SortID
	}).ToSlice(&parentNavs)

	for _, parentNav := range parentNavs {
		parentMenuInfo := &MenuInfo{
			Title:  parentNav.Title,
			Icon:   parentNav.IconURL,
			Href:   parentNav.LinkURL,
			Target: "_self",
			Child:  []*MenuInfo{},
		}
		childMenuInfo, _ := svc.ChildMenu(navs, parentNav.NavID)
		if childMenuInfo != nil {
			parentMenuInfo.Child = append(parentMenuInfo.Child, childMenuInfo...)
		}
		outerMenu.MenuInfo = append(outerMenu.MenuInfo, parentMenuInfo)
	}

	return outerMenu
}

func (svc *Menu) ChildMenu(navs []*model.CmsAdminNav, parentId int64) ([]*MenuInfo, error) {
	if len(navs) == 0 {
		return nil, nil
	}
	outMenu := []*MenuInfo{}

	childNavs := []*model.CmsAdminNav{}
	linq.From(navs).WhereT(func(s *model.CmsAdminNav) bool {
		return s.ParentID == parentId && s.IsHide == 0
	}).OrderByT(func(s *model.CmsAdminNav) int32 {
		return s.SortID
	}).ToSlice(&childNavs)

	for _, childNav := range childNavs {
		fmt.Println(childNav.ParentID, childNav.Title, "-------->", childNav.SortID)
		childMenuInfo := &MenuInfo{
			Title:  childNav.Title,
			Icon:   childNav.IconURL,
			Href:   childNav.LinkURL,
			Target: "_self",
			Child:  []*MenuInfo{},
		}
		childMenuInfo.Child, _ = svc.ChildMenu(navs, childNav.NavID)
		outMenu = append(outMenu, childMenuInfo)
	}
	return outMenu, nil
}
