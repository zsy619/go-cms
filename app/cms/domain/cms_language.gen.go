package domain

import (
	"time"
)

const TableNameCmsLanguage = "cms_language"

// CmsLanguage 语言设置表
type CmsLanguage struct {
	LanguageID  int64     `gorm:"column:language_id;type:bigint;primaryKey;autoIncrement:true" json:"language_id" form:"language_id"` // 主键
	Name        string    `gorm:"column:name;type:varchar(64);not null" json:"name" form:"name"`                                    // 语言名称
	Code        string    `gorm:"column:code;type:varchar(16);not null" json:"code" form:"code"`                                  // 语言代码 (如: zh-CN, en)
	Icon        string    `gorm:"column:icon;type:varchar(64)" json:"icon" form:"icon"`                                          // 图标/国旗
	SortID      int32     `gorm:"column:sort_id;type:int;default:99" json:"sort_id" form:"sort_id"`                             // 排序
	Description string    `gorm:"column:description;type:varchar(256)" json:"description" form:"description"`                       // 描述
	Status      int32     `gorm:"column:status;type:tinyint;default:1" json:"status" form:"status"`                              // 状态: 0禁用 1启用
	IsDefault   bool      `gorm:"column:is_default;type:tinyint(1);default:false" json:"is_default" form:"is_default"`             // 是否默认
	CreateID    int32     `gorm:"column:create_id;type:int" json:"create_id" form:"create_id"`                                   // 创建人ID
	CreateName  string    `gorm:"column:create_name;type:varchar(64)" json:"create_name" form:"create_name"`                       // 创建人姓名
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP" json:"create_time" form:"create_time"`   // 创建时间
	UpdateID    int32     `gorm:"column:update_id;type:int" json:"update_id" form:"update_id"`                                    // 更新人ID
	UpdateName  string    `gorm:"column:update_name;type:varchar(64)" json:"update_name" form:"update_name"`                       // 更新人姓名
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime" json:"update_time" form:"update_time"`                           // 修改时间
	TenantID    int64     `gorm:"column:tenant_id;type:bigint;not null;default:0;index" json:"tenant_id" form:"tenant_id"`        // 租户ID
	Deleted     bool      `gorm:"column:deleted;type:tinyint(1);not null;default:0;index" json:"deleted" form:"deleted"`           // 逻辑删除
}

// TableName CmsLanguage's table name
func (*CmsLanguage) TableName() string {
	return TableNameCmsLanguage
}
