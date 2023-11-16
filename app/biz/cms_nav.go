package biz

import (
	"errors"
	"fmt"
	"time"

	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type CmsAdminNav struct{}

func NewCmsAdminNav() *CmsAdminNav {
	return &CmsAdminNav{}
}

// NavSaveSortId 更新排序
func (svc *CmsAdminNav) NavSaveSortId(navId int64, updateId int32, updateName string, sortId int32) error {
	mdl, do := query.CmsAdminNavDo()
	_, err := do.Where(mdl.NavID.Eq(navId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
			mdl.UpdateID.ColumnName().String():   updateId,
			mdl.UpdateName.ColumnName().String(): updateName,
		},
	)
	return err
}

// NavFind 通过nav_id获取详情
func (svc *CmsAdminNav) NavFind(navId int64) (*model.CmsAdminNav, error) {
	mdl, do := query.CmsAdminNavDo()
	return do.Where(mdl.NavID.Eq(navId)).First()
}

// NavDestroy 根据nav_id删除
func (svc *CmsAdminNav) NavDestroy(navId int64) error {
	mdl, err := svc.NavFind(navId)
	if err != nil {
		return err
	}
	// 判断是否关联站点
	sitMdl, siteDo := query.CmsSiteDo()
	if sitCount, _ := siteDo.Where(sitMdl.SiteID.Eq(mdl.SiteID)).Count(); sitCount > 0 {
		return errors.New("该导航与站点有关联，无法删除")
	}
	// 判断是否关联频道
	channelMdl, channelDo := query.CmsSiteChannelDo()
	if channelCount, _ := channelDo.Where(channelMdl.ChannelID.Eq(mdl.ChannelID)).Count(); channelCount > 0 {
		return errors.New("该导航与频道有关联，无法删除")
	}
	// 判断是否有子节点
	navMdl, navDo := query.CmsAdminNavDo()
	if navCount, _ := navDo.Where(navMdl.ParentID.Eq(navId)).Count(); navCount > 0 {
		return errors.New("该导航下有子节点，无法删除")
	}
	// 根据nav_id删除
	if _, navErr := navDo.Where(navMdl.NavID.Eq(navId)).Delete(); navErr != nil {
		return navErr
	}
	return nil
}

// NavTree 获取树形结构列表
func (svc *CmsAdminNav) NavTree(navId int64) ([]*bizmodel.TreeNode, error) {
	out := make([]*bizmodel.TreeNode, 0)
	navMdl, navDo := query.CmsAdminNavDo()
	list, err := navDo.Where(navMdl.ParentID.Eq(0)).Order(navMdl.SortID).Find()
	if err != nil {
		return out, err
	}
	for _, item := range list {
		child := &bizmodel.TreeNode{
			Id:       item.NavID,
			Name:     item.Title,
			Open:     true,
			Checked:  item.NavID == navId,
			Selected: item.NavID == navId,
			Children: nil,
		}
		children, _ := svc.NavTreeByParentId(item.NavID, navId)
		if children != nil {
			child.Children = children
		}
		out = append(out, child)
	}
	return out, nil
}

// NavTreeByParentId 根据父ID获取分类树
func (svc *CmsAdminNav) NavTreeByParentId(parentId, navId int64) ([]*bizmodel.TreeNode, error) {
	out := make([]*bizmodel.TreeNode, 0)
	navMdl, navDo := query.CmsAdminNavDo()
	list, err := navDo.Where(navMdl.ParentID.Eq(parentId)).Order(navMdl.SortID).Find()
	if err != nil {
		return out, err
	}
	for _, item := range list {
		child := &bizmodel.TreeNode{
			Id:       item.NavID,
			Name:     item.Title,
			Checked:  item.NavID == navId,
			Selected: item.NavID == navId,
			Children: nil,
		}
		children, _ := svc.NavTreeByParentId(item.NavID, parentId)
		if children != nil {
			child.Children = children
		}
		out = append(out, child)
	}
	return out, nil
}

// NavSave 保存导航详情
func (svc *CmsAdminNav) NavSave(input *model.CmsAdminNav) error {
	if input.ParentID < 0 {
		return errors.New("归属节点不能为空")
	}
	if input.Type == "" {
		return errors.New("类型不能为空")
	}
	if input.Name == "" {
		return errors.New("导航标识不能为空")
	}
	if input.Title == "" {
		return errors.New("导航名称不能为空")
	}
	mdl, do := query.CmsAdminNavDo()
	var err error

	if input.NavID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		input.UpdateTime = time.Now()
		_, err = do.Where(mdl.NavID.Eq(input.NavID)).Updates(map[string]interface{}{
			mdl.ParentID.ColumnName().String():   input.ParentID,
			mdl.SiteID.ColumnName().String():     input.SiteID,
			mdl.ChannelID.ColumnName().String():  input.ChannelID,
			mdl.Type.ColumnName().String():       input.Type,
			mdl.Name.ColumnName().String():       input.Name,
			mdl.Title.ColumnName().String():      input.Title,
			mdl.SubTitle.ColumnName().String():   input.SubTitle,
			mdl.IconURL.ColumnName().String():    input.IconURL,
			mdl.LinkURL.ColumnName().String():    input.LinkURL,
			mdl.IsHide.ColumnName().String():     input.IsHide,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.Action.ColumnName().String():     input.Action,
			mdl.IsSys.ColumnName().String():      input.IsSys,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.UpdateID.ColumnName().String():   input.UpdateID,
			mdl.UpdateName.ColumnName().String(): input.UpdateName,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return err
}
