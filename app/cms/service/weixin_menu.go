package service

import (
	"errors"
	"time"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
)

type WeixinMenu struct{}

func NewWeixinMenu() *WeixinMenu {
	return &WeixinMenu{}
}

// MenuPaginate 分页
func (svc *WeixinMenu) MenuPaginate(page, limit int, accountId int64) ([]*domain.WeixinMenu, int64, error) {
	mdl, do := mapper.WeixinMenuDo()
	if accountId > 0 {
		do = do.Where(mdl.AccountID.Eq(accountId))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// MenuFindByParentId 根据父级ID获取
func (svc *WeixinMenu) MenuFindByParentId(accountId, parentId int64) ([]*domain.WeixinMenu, error) {
	mdl, do := mapper.WeixinMenuDo()
	return do.Where(mdl.AccountID.Eq(accountId), mdl.ParentID.Eq(parentId)).Order(mdl.SortID).Find()
}

// MenuFind 获取
func (svc *WeixinMenu) MenuFind(menuId int64) (*domain.WeixinMenu, error) {
	mdl, do := mapper.WeixinMenuDo()
	return do.Where(mdl.MenuID.Eq(menuId)).First()
}

// MenuSave 保存或更新
func (svc *WeixinMenu) MenuSave(input *domain.WeixinMenu) error {
	mdl, do := mapper.WeixinMenuDo()
	if input.Name != "" {
		if count, _ := do.Where(mdl.MenuID.Neq(input.MenuID), mdl.AccountID.Eq(input.AccountID), mdl.ParentID.Eq(input.ParentID), mdl.Name.Eq(input.Name)).Count(); count > 0 {
			return errors.New("菜单名称重复")
		}
	}
	if input.Key != "" {
		if count, _ := do.Where(mdl.MenuID.Neq(input.MenuID), mdl.Key.Eq(input.Key)).Count(); count > 0 {
			return errors.New("菜单Key重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	if input.MenuID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.MenuID.Eq(input.MenuID)).Updates(map[string]interface{}{
			mdl.Name.ColumnName().String():       input.Name,
			mdl.Type.ColumnName().String():       input.Type,
			mdl.Key.ColumnName().String():        input.Key,
			mdl.URL.ColumnName().String():        input.URL,
			mdl.MediaID.ColumnName().String():    input.MediaID,
			mdl.AppID.ColumnName().String():      input.AppID,
			mdl.PagePath.ColumnName().String():   input.PagePath,
			mdl.ArticleID.ColumnName().String():  input.ArticleID,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
			mdl.UpdateID.ColumnName().String():   input.UpdateID,
			mdl.UpdateName.ColumnName().String(): input.UpdateName,
		})
	}
	return err
}

// MenuSaveSortId 保存排序
func (svc *WeixinMenu) MenuSaveSortId(menuId int64, sortId int32) error {
	mdl, do := mapper.WeixinMenuDo()
	_, err := do.Where(mdl.MenuID.Eq(menuId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// MenuDestory 删除
func (svc *WeixinMenu) MenuDestory(menuId int64) error {
	mdl, do := mapper.WeixinMenuDo()
	if count, _ := do.Where(mdl.ParentID.Eq(menuId)).Count(); count > 0 {
		return errors.New("请先删除子菜单")
	}
	if _, err := do.Where(mdl.MenuID.Eq(menuId)).Delete(); err != nil {
		return err
	}
	return nil
}
