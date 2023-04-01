package bmodel

type CategoryOneModel struct {
	ChannelName    string `gorm:"column:channel_name;type:varchar(128)" json:"channel_name" form:"channel_name"`          // 频道名称
	ChannelTitle   string `gorm:"column:channel_title;type:varchar(128)" json:"channel_title" form:"channel_title"`       // 频道标题
	CategoryID     int64  `gorm:"column:category_id;type:bigint;" json:"category_id" form:"category_id"`                  // 主键
	ParentID       int64  `gorm:"column:parent_id;type:bigint" json:"parent_id" form:"parent_id"`                         // 父节点
	SiteID         int64  `gorm:"column:site_id;type:bigint" json:"site_id" form:"site_id"`                               // 所属站点
	ChannelID      int64  `gorm:"column:channel_id;type:bigint" json:"channel_id" form:"channel_id"`                      // 所属频道
	Title          string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                               // 类别标题
	CallIndex      string `gorm:"column:call_index;type:varchar(64)" json:"call_index" form:"call_index"`                 // 调用别名
	ClassLayer     int32  `gorm:"column:class_layer;type:int;default:1" json:"class_layer" form:"class_layer"`            // 类别深度
	LinkURL        string `gorm:"column:link_url;type:varchar(256)" json:"link_url" form:"link_url"`                      // 外部链接
	ImgURL         string `gorm:"column:img_url;type:varchar(256)" json:"img_url" form:"img_url"`                         // 图片地址
	SeoTitle       string `gorm:"column:seo_title;type:varchar(128)" json:"seo_title" form:"seo_title"`                   // SEO标题
	SeoKeyword     string `gorm:"column:seo_keyword;type:varchar(128)" json:"seo_keyword" form:"seo_keyword"`             // SEO关健字
	SeoDescription string `gorm:"column:seo_description;type:varchar(128)" json:"seo_description" form:"seo_description"` // SEO描述
	Content        string `gorm:"column:content;type:text" json:"content" form:"content"`                                 // 内容介绍
	SortID         int32  `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                                  // 排序
	IsShow         bool   `gorm:"column:is_show;type:tinyint(1);default:1" json:"is_show" form:"is_show"`                 // 是否显示:1显示，0隐藏
	IsSearch       bool   `gorm:"column:is_search;type:tinyint(1);default:1" json:"is_search" form:"is_search"`           // 允许检索:1允许，0禁止
	IsDeleted      bool   `gorm:"column:is_deleted;type:tinyint(1)" json:"is_deleted" form:"is_deleted"`                  // 删除标识
}
