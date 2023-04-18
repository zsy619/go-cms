package bizmodel

type ApiLinkListModel struct {
	LinkID     int64  `gorm:"column:link_id;type:bigint;primaryKey;" json:"link_id" form:"link_id"`      // 主键
	SiteID     int64  `gorm:"column:site_id;type:bigint" json:"site_id" form:"site_id"`                  // 所属站点
	ChannelID  int64  `gorm:"column:channel_id;type:bigint" json:"channel_id" form:"channel_id"`         // 所属频道
	CategoryID int64  `gorm:"column:category_id;type:bigint" json:"category_id" form:"category_id"`      // 类别ID
	Title      string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                  // 标题
	LinkURL    string `gorm:"column:link_url;type:varchar(256)" json:"link_url" form:"link_url"`         // 外部链接
	Target     string `gorm:"column:target;type:varchar(16);default:_blank" json:"target" form:"target"` // 是否开启浏览器新窗口
	ImgUrl1    string `gorm:"column:img_url1;type:varchar(256)" json:"img_url1" form:"img_url1"`         // 网站logo地址
	ImgUrl2    string `gorm:"column:img_url2;type:varchar(256)" json:"img_url2" form:"img_url2"`         // 网站logo地址
	IsLock     int32  `gorm:"column:is_lock;type:tinyint" json:"is_lock" form:"is_lock"`                 // 是否锁定（不允许编辑）
	IsTop      int32  `gorm:"column:is_top;type:tinyint" json:"is_top" form:"is_top"`                    // 是否置顶
	IsRed      int32  `gorm:"column:is_red;type:tinyint" json:"is_red" form:"is_red"`                    // 是否推荐
	IsHot      int32  `gorm:"column:is_hot;type:tinyint" json:"is_hot" form:"is_hot"`                    // 是否热门
	IsSlide    int32  `gorm:"column:is_slide;type:tinyint" json:"is_slide" form:"is_slide"`              // 是否幻灯片
}
