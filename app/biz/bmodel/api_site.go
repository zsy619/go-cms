package bmodel

type ApiChannelFindModel struct {
	ChannelID  int64  `gorm:"column:channel_id;type:bigint;primaryKey;" json:"channel_id" form:"channel_id"` // 主键
	ParentID   int64  `gorm:"column:parent_id;type:bigint" json:"parent_id" form:"parent_id"`                // 父级ID
	Title      string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                      // 标题
	Name       string `gorm:"column:name;type:varchar(128)" json:"name" form:"name"`                         // 频道名称
	Kind       int32  `gorm:"column:kind;type:tinyint" json:"kind" form:"kind"`                              // 频道类型
	ClassLayer int32  `gorm:"column:class_layer;type:tinyint" json:"class_layer" form:"class_layer"`         // 层级
	SortID     int32  `gorm:"column:sort_id;type:tinyint" json:"sort_id" form:"sort_id"`                     // 排序
	IsAlbum    int32  `gorm:"column:is_album;type:tinyint" json:"is_album" form:"is_album"`                  // 是否相册
	IsAttach   int32  `gorm:"column:is_attach;type:tinyint" json:"is_attach" form:"is_attach"`               // 是否附件
	IsSpec     int32  `gorm:"column:is_spec;type:tinyint" json:"is_spec" form:"is_spec"`                     // 是否规格
}
