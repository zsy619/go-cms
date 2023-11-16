package biz

import (
	"errors"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type PlgOnlineRegister struct{}

func NewPlgOnlineRegister() *PlgOnlineRegister {
	return &PlgOnlineRegister{}
}

// Paginate 分页查询
func (svc *PlgOnlineRegister) Paginate(page, limit int, realName, special, degree, tags, remark string, isRead int32) ([]*model.PlgOnlineRegister, int64, error) {
	mdl, do := query.PlgOnlineRegisterDo()
	if realName != "" {
		do.Where(mdl.RealName.Like("%" + realName + "%"))
	}
	if special != "" {
		do.Where(mdl.Special.Like("%" + special + "%"))
	}
	if degree != "" {
		do.Where(mdl.Degree.Like("%" + degree + "%"))
	}
	if tags != "" {
		do.Where(mdl.Tags.Like("%" + tags + "%"))
	}
	if remark != "" {
		do.Where(mdl.Remark.Like("%" + remark + "%"))
	}
	if isRead == 0 {
		do.Where(mdl.IsRead.Eq(0))
	}
	if isRead == 1 {
		do.Where(mdl.IsRead.Eq(1))
	}
	return do.Order(mdl.CreateTime.Desc()).FindByPage((page-1)*limit, limit)
}

// Find 获取
func (svc *PlgOnlineRegister) Find(registerId int64) (*model.PlgOnlineRegister, error) {
	mdl, do := query.PlgOnlineRegisterDo()
	return do.Where(mdl.RegisterID.Eq(registerId)).First()
}

// ChangeRead 设置已读或未读
func (svc *PlgOnlineRegister) ChangeRead(registerId int64, isRead int32, updateId int32, updateName string) error {
	mdl, do := query.PlgOnlineRegisterDo()
	_, err := do.Where(mdl.RegisterID.Eq(registerId)).Updates(map[string]interface{}{
		mdl.IsRead.ColumnName().String():     isRead,
		mdl.UpdateID.ColumnName().String():   updateId,
		mdl.UpdateName.ColumnName().String(): updateName,
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return err
}

// LinkSave 保存或更新
func (svc *PlgOnlineRegister) Save(input *model.PlgOnlineRegister) error {
	mdl, do := query.PlgOnlineRegisterDo()
	if input.Mobile != "" {
		if count, _ := do.Where(mdl.Mobile.Eq(input.Mobile), mdl.RegisterID.Neq(input.RegisterID)).Count(); count > 0 {
			return errors.New("手机号已注册~")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	if input.RegisterID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.RegisterID.Eq(input.RegisterID)).Updates(map[string]interface{}{
			mdl.RealName.ColumnName().String():   input.RealName,
			mdl.Mobile.ColumnName().String():     input.Mobile,
			mdl.Email.ColumnName().String():      input.Email,
			mdl.Special.ColumnName().String():    input.Special,
			mdl.Degree.ColumnName().String():     input.Degree,
			mdl.Sex.ColumnName().String():        input.Sex,
			mdl.Year.ColumnName().String():       input.Year,
			mdl.Wechat.ColumnName().String():     input.Wechat,
			mdl.Address.ColumnName().String():    input.Address,
			mdl.Content.ColumnName().String():    input.Content,
			mdl.Tags.ColumnName().String():       input.Tags,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.UpdateID.ColumnName().String():   input.UpdateID,
			mdl.UpdateName.ColumnName().String(): input.UpdateName,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

// Destory 删除
func (svc *PlgOnlineRegister) Destory(registerId int64) error {
	mdl, do := query.PlgOnlineRegisterDo()
	if _, err := do.Where(mdl.RegisterID.Eq(registerId)).Delete(); err != nil {
		return err
	}
	return nil
}
