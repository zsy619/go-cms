package bizmodel

type ApiTopicModel struct {
	TopicID        int64  `gorm:"column:topic_id;type:bigint;comment:专题ID" json:"topic_id" form:"topic_id"`
	SiteID         int64  `gorm:"column:site_id;type:bigint;comment:所属站点" json:"site_id" form:"site_id"`
	ChannelID      int64  `gorm:"column:channel_id;type:bigint;comment:所属频道" json:"channel_id" form:"channel_id"`
	Name           string `gorm:"column:name;type:varchar(64);comment:专题名称" json:"name" form:"name"`
	Title          string `gorm:"column:title;type:varchar(128);comment:专题标题" json:"title" form:"title"`
	ImgUrl1        string `gorm:"column:img_url1;type:varchar(256);comment:图片" json:"img_url1" form:"img_url1"`
	ImgUrl2        string `gorm:"column:img_url2;type:varchar(256);comment:图片" json:"img_url2" form:"img_url2"`
	SeoTitle       string `gorm:"column:seo_title;type:varchar(128);comment:SEO标题" json:"seo_title" form:"seo_title"`
	SeoKeyword     string `gorm:"column:seo_keyword;type:varchar(128);comment:SEO关健字" json:"seo_keyword" form:"seo_keyword"`
	SeoDescription string `gorm:"column:seo_description;type:varchar(128);comment:SEO描述" json:"seo_description" form:"seo_description"`
	Remark         string `gorm:"column:remark;type:varchar(256);comment:备注" json:"remark" form:"remark"`
	SortID         int32  `gorm:"column:sort_id;type:int;comment:排序" json:"sort_id" form:"sort_id"`
	Click          int32  `gorm:"column:click;type:int;comment:浏览次数" json:"click" form:"click"`
	Template       string `gorm:"column:template;type:varchar(256);comment:模板路径" json:"template" form:"template"`
	SiteFlag       string `gorm:"column:site_flag;type:varchar(64);comment:站点标识" json:"site_flag" form:"site_flag"`
}
