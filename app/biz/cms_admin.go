package biz

import (
	"context"
	"errors"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xcrypto"
)

type LoginType int

const (
	LoginName LoginType = iota + 1
	LoginMobile
	LoginEmail
	LoginNameMobile
	LoginNameEmail
	LoginMobileEmail
	LoginAll
)

type CmsAdmin struct{}

func NewCmsAdmin() *CmsAdmin {
	return &CmsAdmin{}
}

// Login 登录
func (m *CmsAdmin) Login(login_key, password string, user_type int, login_type LoginType) (*model.CmsAdmin, error) {
	mdl, do := query.CmsAdminDo()
	ctx := context.Background()
	switch login_type {
	case LoginName:
		do = do.Where(mdl.IsDeleted.Is(false), mdl.UserName.Eq(login_key))
	case LoginMobile:
		do = do.Where(mdl.IsDeleted.Is(false), mdl.Mobile.Eq(login_key))
	case LoginEmail:
		do = do.Where(mdl.IsDeleted.Is(false), mdl.Email.Eq(login_key))
	case LoginNameMobile:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.UserName.Eq(login_key)).Or(mdl.Mobile.Eq(login_key)),
		)
	case LoginNameEmail:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.UserName.Eq(login_key)).Or(mdl.Email.Eq(login_key)),
		)
	case LoginMobileEmail:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.Mobile.Eq(login_key)).Or(mdl.Email.Eq(login_key)),
		)
	case LoginAll:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.UserName.Eq(login_key)).Or(mdl.Mobile.Eq(login_key)).Or(mdl.Email.Eq(login_key)),
		)
	}
	find, err := do.First()
	if err != nil {
		return nil, err
	}
	if find == nil {
		return nil, errors.New("用户不存在")
	}
	// 密码加密格式 0不加密 1默认加密 2MD5类型
	switch find.PasswordFormat {
	case 1: // 默认加密

		break
	case 2: // MD5加密
		password = xcrypto.GetMD5Hash(password)
		break
	default:
		break
	}
	logs.Debug("Login --> ", find.Password, password)
	if find.Password == password {
		return find, nil
	}
	return nil, errors.New("账号密码错误")
}

// LoginLog 写入登录日志
func (m *CmsAdmin) LoginLog(userId int64, userName string, method, path, queryx, statusCode, ip string) {
	_, logDo := query.CmsAdminLogDo()
	log := &model.CmsAdminLog{
		UserID:     userId,
		UserName:   userName,
		Method:     method,
		Path:       path,
		Query:      queryx,
		StatusCode: statusCode,
		IP:         ip,
		CreateTime: time.Now(),
	}
	logDo.Create(log)
	mdl, do := query.CmsAdminDo()
	do.Where(mdl.UserID.Eq(userId)).UpdateColumns(
		map[string]interface{}{
			mdl.LastIP.ColumnName().String():   ip,
			mdl.LastTime.ColumnName().String(): time.Now().Format("2006-01-02 15:04:05"),
		},
	)
}

func (m *CmsAdmin) LogPaginate(page, limit int, userId int64, userName string) ([]*model.CmsAdminLog, int64, error) {
	mdl, do := query.CmsAdminLogDo()
	if userId > 0 {
		do = do.Where(mdl.UserID.Eq(userId))
	}
	return do.Where(mdl.UserName.Like("%"+userName+"%")).FindByPage((page-1)*limit, limit)
}

func (m *CmsAdmin) RolePaginate(page, limit int, name string) ([]*model.CmsAdminRole, int64, error) {
	mdl, do := query.CmsAdminRoleDo()
	return do.Where(mdl.Name.Like("%"+name+"%")).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

func (m *CmsAdmin) RoleFind(roleId int64) (*model.CmsAdminRole, error) {
	mdl, do := query.CmsAdminRoleDo()
	return do.Where(mdl.RoleID.Eq(roleId)).First()
}

// RoleSave 保存或更新
func (this *CmsAdmin) RoleSave(input *model.CmsAdminRole) error {
	mdl, do := query.CmsAdminRoleDo()
	if input.Name != "" {
		if count, _ := do.Where(mdl.RoleID.Neq(input.RoleID), mdl.Name.Eq(input.Name)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	if input.RoleID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.RoleID.Eq(input.RoleID)).Updates(map[string]interface{}{
			mdl.Name.ColumnName().String():       input.Name,
			mdl.Type.ColumnName().String():       input.Type,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

func (this *CmsAdmin) RoleValSave(roleId int64, values map[string]string) error {
	mdl, do := query.CmsAdminRoleValueDo()
	_, err := do.Where(mdl.RoleID.Eq(roleId)).Delete()
	for k, v := range values {
		err = do.Save(&model.CmsAdminRoleValue{
			RoleID:  roleId,
			NavName: k,
			Action:  v,
		})
	}
	return err
}

// RoleDestory 删除
func (this *CmsAdmin) RoleDestory(roleId int64) error {
	adminMdl, adminDo := query.CmsAdminDo()
	if count, _ := adminDo.Where(adminMdl.RoleID.Eq(roleId)).Count(); count > 0 {
		return errors.New("该角色下有用户，无法删除")
	}
	valueMdl, valueDo := query.CmsAdminRoleValueDo()
	if _, err := valueDo.Where(valueMdl.RoleID.Eq(roleId)).Delete(); err != nil {
		return err
	}
	mdl, do := query.CmsAdminRoleDo()
	if _, err := do.Where(mdl.RoleID.Eq(roleId)).Delete(); err != nil {
		return err
	}
	return nil
}

func (this *CmsAdmin) RoleSaveSortId(roleId int64, sortId int32) error {
	mdl, do := query.CmsAdminRoleDo()
	_, err := do.Where(mdl.RoleID.Eq(roleId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

func (this *CmsAdmin) RoleValueFind(roleId int64) ([]*model.CmsAdminRoleValue, int64, error) {
	mdl, do := query.CmsAdminRoleValueDo()
	return do.Where(mdl.RoleID.Eq(roleId)).FindByPage(0, 99999)
}

func (this *CmsAdmin) NavFind(roleId int64) ([]*model.CmsAdminNav, int64, error) {
	mdl, do := query.CmsAdminNavDo()
	return do.Order(mdl.SortID).FindByPage(0, 99999)
}
