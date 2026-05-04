package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xcrypto"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
	"haedu.gov.cn/cms/global"
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
func (svc *CmsAdmin) Login(userKey, password string, userType int, loginType LoginType) (*domain.CmsAdmin, error) {
	mdl, do := mapper.CmsAdminDo()
	ctx := context.Background()
	switch loginType {
	case LoginName:
		do = do.Where(mdl.IsDeleted.Is(false), mdl.UserName.Eq(userKey))
	case LoginMobile:
		do = do.Where(mdl.IsDeleted.Is(false), mdl.Mobile.Eq(userKey))
	case LoginEmail:
		do = do.Where(mdl.IsDeleted.Is(false), mdl.Email.Eq(userKey))
	case LoginNameMobile:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.UserName.Eq(userKey)).Or(mdl.Mobile.Eq(userKey)),
		)
	case LoginNameEmail:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.UserName.Eq(userKey)).Or(mdl.Email.Eq(userKey)),
		)
	case LoginMobileEmail:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.Mobile.Eq(userKey)).Or(mdl.Email.Eq(userKey)),
		)
	case LoginAll:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.UserName.Eq(userKey)).Or(mdl.Mobile.Eq(userKey)).Or(mdl.Email.Eq(userKey)),
		)
	}
	find, err := do.First()
	if err != nil {
		return nil, err
	}
	if find == nil {
		return nil, errors.New("用户不存在")
	}
	if !find.Enabled {
		return nil, errors.New("用户已禁用")
	}
	// 密码加密格式 0不加密 1默认加密 2MD5类型
	key := global.ReverseLowerString(userKey)
	password, _ = xcrypto.Sm4Encrypt(password, key)
	logs.Debug("Login --> ", find.Password, password)
	if find.Password == password {
		return find, nil
	}
	return nil, errors.New("账号密码错误")
}

// LoginLog 写入登录日志
func (*CmsAdmin) LoginLog(userId int64, userName string, method, path, queryx, statusCode, ip string) {
	_, logDo := mapper.CmsAdminLogDo()
	log := &domain.CmsAdminLog{
		UserID:     userId,
		UserName:   userName,
		Method:     method,
		Path:       path,
		Query:      queryx,
		StatusCode: statusCode,
		IP:         ip,
		CreateTime: time.Now(),
	}
	_ = logDo.Create(log)
	mdl, do := mapper.CmsAdminDo()
	_, _ = do.Where(mdl.UserID.Eq(userId)).UpdateColumns(
		map[string]interface{}{
			mdl.LastIP.ColumnName().String():   ip,
			mdl.LastTime.ColumnName().String(): time.Now().Format("2006-01-02 15:04:05"),
		},
	)
}

func (*CmsAdmin) AdminPaginate(page, limit int, roleId int64, realName, userName string) ([]*domain.CmsAdmin, int64, error) {
	mdl, do := mapper.CmsAdminDo()
	if roleId > 0 {
		do = do.Where(mdl.RoleID.Eq(roleId))
	}
	if realName != "" {
		do = do.Where(mdl.RealName.Like("%" + realName + "%"))
	}
	if userName != "" {
		do = do.Where(mdl.UserName.Like("%" + userName + "%"))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// AdminDestory 删除
func (svc *CmsAdmin) AdminDestory(userId int64) error {
	mdl, do := mapper.CmsAdminDo()
	if _, err := do.Where(mdl.UserID.Eq(userId)).Delete(); err != nil {
		return err
	}
	return nil
}

// AdminSave 保存或更新
func (svc *CmsAdmin) AdminSave(input *domain.CmsAdmin) error {
	mdl, do := mapper.CmsAdminDo()
	if input.NickName != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.NickName.Eq(input.NickName)).Count(); count > 0 {
			return errors.New("昵称重复")
		}
	}
	if input.UserName != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.UserName.Eq(input.UserName)).Count(); count > 0 {
			return errors.New("登录名重复")
		}
	}
	if input.Email != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.Email.Eq(input.Email)).Count(); count > 0 {
			return errors.New("邮箱重复")
		}
	}
	if input.Mobile != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.Mobile.Eq(input.Mobile)).Count(); count > 0 {
			return errors.New("手机号重复")
		}
	}
	if input.IDCard != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.IDCard.Eq(input.IDCard)).Count(); count > 0 {
			return errors.New("身份证号重复")
		}
	}
	if input.UserNumber != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.UserNumber.Eq(input.UserNumber)).Count(); count > 0 {
			return errors.New("特征标识重复")
		}
	}
	roleMdl, roleDo := mapper.CmsAdminRoleDo()
	var roleType string
	if err := roleDo.Where(roleMdl.RoleID.Eq(input.RoleID)).Pluck(roleMdl.Type, &roleType); err != nil {
		return err
	}
	input.RoleType = roleType
	var err error
	input.UpdateTime = time.Now()
	if input.UserID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.UserID.Eq(input.UserID)).Updates(map[string]interface{}{
			mdl.RoleID.ColumnName().String():     input.RoleID,
			mdl.RoleType.ColumnName().String():   input.RoleType,
			mdl.RealName.ColumnName().String():   input.RealName,
			mdl.NickName.ColumnName().String():   input.NickName,
			mdl.IDCard.ColumnName().String():     input.IDCard,
			mdl.Sex.ColumnName().String():        input.Sex,
			mdl.UserName.ColumnName().String():   input.UserName,
			mdl.UserNumber.ColumnName().String(): input.UserNumber,
			mdl.Email.ColumnName().String():      input.Email,
			mdl.Telphone.ColumnName().String():   input.Telphone,
			mdl.Mobile.ColumnName().String():     input.Mobile,
			mdl.MobilePin.ColumnName().String():  input.MobilePin,
			mdl.UserType.ColumnName().String():   input.UserType,
			mdl.Enabled.ColumnName().String():    input.Enabled,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.Avatar.ColumnName().String():     input.Avatar,
			mdl.IsAudit.ColumnName().String():    input.IsAudit,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
		if input.Password != "" {
			_, err = do.Where(mdl.RoleID.Eq(input.RoleID)).Updates(map[string]interface{}{
				mdl.Password.ColumnName().String():       input.Password,
				mdl.PasswordSalt.ColumnName().String():   input.PasswordSalt,
				mdl.PasswordFormat.ColumnName().String(): input.PasswordFormat,
			})
		}
	}
	return err
}

func (svc *CmsAdmin) AdminSaveSortId(userId int64, sortId int32) error {
	mdl, do := mapper.CmsAdminDo()
	_, err := do.Where(mdl.UserID.Eq(userId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// FindByAccount 根据账号查找
func (svc *CmsAdmin) FindByAccount(userName string) (*domain.CmsAdmin, error) {
	mdl, do := mapper.CmsAdminDo()
	return do.Where(mdl.UserName.Eq(userName)).First()
}

func (m *CmsAdmin) AdminFind(userId int64) (*domain.CmsAdmin, error) {
	mdl, do := mapper.CmsAdminDo()
	return do.Where(mdl.UserID.Eq(userId)).First()
}

func (m *CmsAdmin) LogPaginate(page, limit int, userId int64, userName string) ([]*domain.CmsAdminLog, int64, error) {
	mdl, do := mapper.CmsAdminLogDo()
	if userId > 0 {
		do = do.Where(mdl.UserID.Eq(userId))
	}
	return do.Where(mdl.UserName.Like("%"+userName+"%")).Order(mdl.LogID.Desc()).FindByPage((page-1)*limit, limit)
}

func (m *CmsAdmin) RolePaginate(page, limit int, name string) ([]*domain.CmsAdminRole, int64, error) {
	mdl, do := mapper.CmsAdminRoleDo()
	return do.Where(mdl.Name.Like("%"+name+"%")).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

func (m *CmsAdmin) RoleFind(roleId int64) (*domain.CmsAdminRole, error) {
	mdl, do := mapper.CmsAdminRoleDo()
	return do.Where(mdl.RoleID.Eq(roleId)).First()
}

// RoleSave 保存或更新
func (svc *CmsAdmin) RoleSave(input *domain.CmsAdminRole) error {
	mdl, do := mapper.CmsAdminRoleDo()
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

func (svc *CmsAdmin) RoleValSave(roleId int64, values map[string]string) error {
	mdl, do := mapper.CmsAdminRoleValueDo()
	_, err := do.Where(mdl.RoleID.Eq(roleId)).Delete()
	for k, v := range values {
		err = do.Save(&domain.CmsAdminRoleValue{
			RoleID:  roleId,
			NavName: k,
			Action:  v,
		})
	}
	return err
}

// RoleSiteFind 站点权限-根据roleId获取列表
func (svc *CmsAdmin) RoleSiteFind(roleId int64) ([]*domain.CmsAdminRoleSite, int64, error) {
	mdl, do := mapper.CmsAdminRoleSiteDo()
	return do.Where(mdl.RoleID.Eq(roleId)).FindByPage(0, 99999)
}

// RoleSiteSave 站点权限-保存
func (svc *CmsAdmin) RoleSiteSave(roleId int64, values []string) error {
	mdl, do := mapper.CmsAdminRoleSiteDo()
	_, err := do.Where(mdl.RoleID.Eq(roleId)).Delete()
	for i := 0; i < len(values); i++ {
		item, _ := strconv.ParseInt(values[i], 0, 64)
		err = do.Save(&domain.CmsAdminRoleSite{
			RoleID: roleId,
			SiteID: item,
		})
	}
	return err
}

// RoleDestory 删除
func (svc *CmsAdmin) RoleDestory(roleId int64) error {
	adminMdl, adminDo := mapper.CmsAdminDo()
	if count, _ := adminDo.Where(adminMdl.RoleID.Eq(roleId)).Count(); count > 0 {
		return errors.New("该角色下有用户，无法删除")
	}
	valueMdl, valueDo := mapper.CmsAdminRoleValueDo()
	if _, err := valueDo.Where(valueMdl.RoleID.Eq(roleId)).Delete(); err != nil {
		return err
	}
	siteMdl, siteDo := mapper.CmsAdminRoleSiteDo()
	if _, err := siteDo.Where(siteMdl.RoleID.Eq(roleId)).Delete(); err != nil {
		return err
	}
	mdl, do := mapper.CmsAdminRoleDo()
	if _, err := do.Where(mdl.RoleID.Eq(roleId)).Delete(); err != nil {
		return err
	}
	return nil
}

func (svc *CmsAdmin) RoleSaveSortId(roleId int64, sortId int32) error {
	mdl, do := mapper.CmsAdminRoleDo()
	_, err := do.Where(mdl.RoleID.Eq(roleId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

func (svc *CmsAdmin) RoleValueFind(roleId int64) ([]*domain.CmsAdminRoleValue, int64, error) {
	mdl, do := mapper.CmsAdminRoleValueDo()
	return do.Where(mdl.RoleID.Eq(roleId)).FindByPage(0, 99999)
}

func (svc *CmsAdmin) RolePower(roleId int64, navName string) (*domain.CmsAdminRoleValue, error) {
	mdl, do := mapper.CmsAdminRoleValueDo()
	return do.Where(mdl.RoleID.Eq(roleId), mdl.NavName.Eq(navName)).First()
}

func (svc *CmsAdmin) NavFind(roleId int64) ([]*domain.CmsAdminNav, int64, error) {
	mdl, do := mapper.CmsAdminNavDo()
	return do.Order(mdl.SortID).FindByPage(0, 99999)
}

func (svc *CmsAdmin) OneByUserId(userId int64) *domain.CmsAdmin {
	mdl, do := mapper.CmsAdminDo()
	first, err := do.Where(mdl.UserID.Eq(userId)).First()
	if err != nil {
		return nil
	}
	return first
}

func (svc *CmsAdmin) OneByUserName(userName string) *domain.CmsAdmin {
	mdl, do := mapper.CmsAdminDo()
	first, err := do.Where(mdl.UserName.Eq(userName)).First()
	if err != nil {
		return nil
	}
	return first
}

func (svc *CmsAdmin) ModifyPassword(userId int64, oldPassword, newPassword string) error {
	mdl, do := mapper.CmsAdminDo()
	admin, err := do.Where(mdl.UserID.Eq(userId)).First()
	if err != nil {
		return err
	}

	key := global.ReverseLowerString(admin.UserName)
	oldPassword, _ = xcrypto.Sm4Encrypt(oldPassword, key)
	logs.Debug("ModifyPassword --> ", admin.Password, oldPassword)
	if admin.Password != oldPassword {
		return errors.New("旧密码错误")
	}

	password := newPassword
	passwordSalt := ""
	password, _ = xcrypto.Sm4Encrypt(password, key)
	_, err = do.Where(mdl.UserID.Eq(userId)).UpdateColumns(map[string]interface{}{
		mdl.Password.ColumnName().String():       password,
		mdl.PasswordSalt.ColumnName().String():   passwordSalt,
		mdl.PasswordFormat.ColumnName().String(): admin.PasswordFormat,
		mdl.UpdateTime.ColumnName().String():     time.Now(),
	})
	return err
}
