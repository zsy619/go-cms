package model

type ApiTopicModel struct {
	TopicID        int64  `gorm:"column:topic_id;type:bigint;" json:"topic_id" form:"topic_id"`                            // 专题ID
	SiteID         int64  `gorm:"column:site_id;type:bigint;" json:"site_id" form:"site_id"`                               // 所属站点
	ChannelID      int64  `gorm:"column:channel_id;type:bigint;" json:"channel_id" form:"channel_id"`                      // 所属频道
	Name           string `gorm:"column:name;type:varchar(64);" json:"name" form:"name"`                                   // 专题名称
	Title          string `gorm:"column:title;type:varchar(128);" json:"title" form:"title"`                               // 专题标题
	ImgUrl1        string `gorm:"column:img_url1;type:varchar(256);" json:"img_url1" form:"img_url1"`                      // 图片
	ImgUrl2        string `gorm:"column:img_url2;type:varchar(256);" json:"img_url2" form:"img_url2"`                      // 图片
	SeoTitle       string `gorm:"column:seo_title;type:varchar(128);" json:"seo_title" form:"seo_title"`                   // SEO标题
	SeoKeyword     string `gorm:"column:seo_keyword;type:varchar(128);" json:"seo_keyword" form:"seo_keyword"`             // SEO关健字
	SeoDescription string `gorm:"column:seo_description;type:varchar(128);" json:"seo_description" form:"seo_description"` // SEO描述
	Remark         string `gorm:"column:remark;type:varchar(256);" json:"remark" form:"remark"`                            // 备注
	SortID         int32  `gorm:"column:sort_id;type:int;" json:"sort_id" form:"sort_id"`                                  // 排序
	Click          int32  `gorm:"column:click;type:int;" json:"click" form:"click"`                                        // 浏览次数
	Template       string `gorm:"column:template;type:varchar(256);" json:"template" form:"template"`                      // 模板路径
	SiteFlag       string `gorm:"column:site_flag;type:varchar(64);" json:"site_flag" form:"site_flag"`                    // 站点标识
}
