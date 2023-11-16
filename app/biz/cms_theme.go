package biz

import (
	"errors"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type CmsTheme struct{}

func NewCmsTheme() *CmsTheme {
	return &CmsTheme{}
}

// ThemeClone 克隆
func (svc *CmsTheme) ThemeClone(themeId int64) (int64, error) {
	mdl, do := query.CmsThemeDo()
	art, err := do.Where(mdl.ThemeID.Eq(themeId)).First()
	if err != nil {
		return 0, err
	}
	art.ThemeID = 0
	art.CreateTime = time.Now()
	art.UpdateTime = time.Now()
	err = do.Create(art)
	return art.ThemeID, err
}

// ThemePaginate 分页查询
func (svc *CmsTheme) ThemePaginate(page, limit int, name, title string) ([]*model.CmsTheme, int64, error) {
	mdl, do := query.CmsThemeDo()
	if name != "" {
		do = do.Where(mdl.Name.Like("%" + name + "%"))
	}
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// ThemeSetDefault 设置默认主题
func (svc *CmsTheme) ThemeSetDefault(name string) error {
	mdl, do := query.CmsThemeDo()
	now := time.Now()
	if _, err := do.Where(mdl.Name.Eq(name)).UpdateColumns(map[string]interface{}{
		mdl.IsDefault.ColumnName().String():  1,
		mdl.UpdateTime.ColumnName().String(): now,
	}); err != nil {
		return err
	}
	if _, err := do.Where(mdl.Name.Neq(name)).UpdateColumns(map[string]interface{}{
		mdl.IsDefault.ColumnName().String():  0,
		mdl.UpdateTime.ColumnName().String(): now,
	}); err != nil {
		return err
	}
	return nil
}

// ThemeFind 获取
func (svc *CmsTheme) ThemeFind(themeId int64) (*model.CmsTheme, error) {
	mdl, do := query.CmsThemeDo()
	return do.Where(mdl.ThemeID.Eq(themeId)).First()
}

// ThemeSave 保存或更新
func (svc *CmsTheme) ThemeSave(input *model.CmsTheme) error {
	mdl, do := query.CmsThemeDo()
	if input.Name != "" {
		if count, _ := do.Where(mdl.ThemeID.Neq(input.ThemeID), mdl.Name.Eq(input.Name)).Count(); count > 0 {
			return errors.New("模板名称重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	if input.ThemeID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.ThemeID.Eq(input.ThemeID)).Updates(map[string]interface{}{
			mdl.Name.ColumnName().String():       input.Name,
			mdl.Title.ColumnName().String():      input.Title,
			mdl.Thumb.ColumnName().String():      input.Thumb,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.IsDefault.ColumnName().String():  input.IsDefault,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
			mdl.UpdateID.ColumnName().String():   input.UpdateID,
			mdl.UpdateName.ColumnName().String(): input.UpdateName,
		})
	}
	return err
}

/**
 * @description: ThemeDestory 删除
 * @param {int64} ThemeId ID
 * @return {*}
 */
func (svc *CmsTheme) ThemeDestory(ThemeId int64) error {
	mdl, do := query.CmsThemeDo()
	if _, err := do.Where(mdl.ThemeID.Eq(ThemeId)).Delete(); err != nil {
		return err
	}
	return nil
}

/**
 * @description: ThemeSaveSortId 保存排序
 * @param {int64} ThemeId ID
 * @param {int32} sortId 排序
 * @return {*}
 */
func (svc *CmsTheme) ThemeSaveSortId(ThemeId int64, sortId int32) error {
	mdl, do := query.CmsThemeDo()
	_, err := do.Where(mdl.ThemeID.Eq(ThemeId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

/**
 * @description: 判断模板是否存在
 * @param {string} name 模板名称
 * @return {*}
 */
func (svc *CmsTheme) ThemeExists(name string) bool {
	mdl, do := query.CmsThemeDo()
	if count, _ := do.Where(mdl.Name.Eq(name)).Count(); count > 0 {
		return true
	}
	return false
}
