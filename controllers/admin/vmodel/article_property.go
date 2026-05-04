package vmodel

type Article_PropertySaveModel struct {
	PropertyID int64  `gorm:"column:property_id;type:bigint;primaryKey;autoIncrement:true" json:"property_id" form:"property_id"` // 主键
	ParentID   int64  `gorm:"column:parent_id;type:bigint" json:"parent_id" form:"parent_id"`                                     // 父ID
	ArticleID  int64  `gorm:"column:article_id;type:bigint" json:"article_id" form:"article_id"`                                  // 所属文章
	Title      string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                                           // 属性标题
	CallIndex  string `gorm:"column:call_index;type:varchar(64)" json:"call_index" form:"call_index"`                             // 调用别名
	Value      string `gorm:"column:value;type:varchar(128)" json:"value" form:"value"`                                           // 属性值
	SortID     int32  `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                                              // 排序
	Status     int32  `gorm:"column:status;type:tinyint" json:"status" form:"status"`                                             // 状态0草稿1提交2审核通过3审核未通过4驳回
	IsDeleted  bool   `gorm:"column:is_deleted;type:tinyint(1)" json:"is_deleted" form:"is_deleted"`                              // 删除标识
	BelongTo   string `gorm:"column:belong_to;type:varchar(64)" json:"belong_to" form:"belong_to"`                                // 归属
	CreateID   int32  `gorm:"column:create_id;type:int" json:"create_id" form:"create_id"`                                        // 创建人ID
	CreateName string `gorm:"column:create_name;type:varchar(64)" json:"create_name" form:"create_name"`                          // 创建人姓名
	UpdateID   int32  `gorm:"column:update_id;type:int" json:"update_id" form:"update_id"`                                        // 更新人ID
	UpdateName string `gorm:"column:update_name;type:varchar(64)" json:"update_name" form:"update_name"`                          // 更新人姓名
}
