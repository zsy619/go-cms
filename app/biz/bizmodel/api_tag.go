package bizmodel

type ApiTagListModel struct {
	TagID          int64  `gorm:"column:tag_id;type:bigint;comment:标签ID" json:"tag_id" form:"tag_id"`
	SiteID         int64  `gorm:"column:site_id;type:bigint;comment:所属站点" json:"site_id" form:"site_id"`
	ChannelID      int64  `gorm:"column:channel_id;type:bigint;comment:所属频道" json:"channel_id" form:"channel_id"`
	Name           string `gorm:"column:name;type:varchar(64);comment:标签名称" json:"name" form:"name"`
	Title          string `gorm:"column:title;type:varchar(128);comment:标签标题" json:"title" form:"title"`
	ImgUrl1        string `gorm:"column:img_url1;type:varchar(256);comment:图片" json:"img_url1" form:"img_url1"`
	ImgUrl2        string `gorm:"column:img_url2;type:varchar(256);comment:图片" json:"img_url2" form:"img_url2"`
	SeoTitle       string `gorm:"column:seo_title;type:varchar(128);comment:SEO标题" json:"seo_title" form:"seo_title"`
	SeoKeyword     string `gorm:"column:seo_keyword;type:varchar(128);comment:SEO关健字" json:"seo_keyword" form:"seo_keyword"`
	SeoDescription string `gorm:"column:seo_description;type:varchar(128);comment:SEO描述" json:"seo_description" form:"seo_description"`
	SortID         int32  `gorm:"column:sort_id;type:int;comment:排序" json:"sort_id" form:"sort_id"`
}
