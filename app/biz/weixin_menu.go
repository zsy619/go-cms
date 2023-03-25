package biz

import (
	"errors"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type WeixinMenu struct{}

func NewWeixinMenu() *WeixinMenu {
	return &WeixinMenu{}
}

func (this *WeixinMenu) MenuPaginate(page, limit int, accountId int64) ([]*model.WeixinMenu, int64, error) {
	mdl, do := query.WeixinMenuDo()
	if accountId > 0 {
		do = do.Where(mdl.AccountID.Eq(accountId))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// MenuFind 获取
func (this *WeixinMenu) MenuFind(menuId int64) (*model.WeixinMenu, error) {
	mdl, do := query.WeixinMenuDo()
	return do.Where(mdl.MenuID.Eq(menuId)).First()
}

// MenuSave 保存或更新
func (this *WeixinMenu) MenuSave(input *model.WeixinMenu) error {
	mdl, do := query.WeixinMenuDo()
	if input.Name != "" {
		if count, _ := do.Where(mdl.MenuID.Neq(input.MenuID), mdl.Name.Eq(input.Name)).Count(); count > 0 {
			return errors.New("菜单名称重复")
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
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

func (this *WeixinMenu) MenuSaveSortId(menuId int64, sortId int32) error {
	mdl, do := query.WeixinMenuDo()
	_, err := do.Where(mdl.MenuID.Eq(menuId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// MenuDestory 删除
func (this *WeixinMenu) MenuDestory(menuId int64) error {
	mdl, do := query.WeixinMenuDo()
	if count, _ := do.Where(mdl.ParentID.Eq(menuId)).Count(); count > 0 {
		return errors.New("请先删除子菜单")
	}
	if _, err := do.Where(mdl.MenuID.Eq(menuId)).Delete(); err != nil {
		return err
	}
	return nil
}
