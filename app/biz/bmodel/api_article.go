package bmodel

import "time"

// ApiArticleModel 文章模型
type ApiCategoryFindModel struct {
	ChannelName  string `gorm:"column:channel_name;type:varchar(128)" json:"channel_name" form:"channel_name"`    // 频道名称
	ChannelTitle string `gorm:"column:channel_title;type:varchar(128)" json:"channel_title" form:"channel_title"` // 频道标题
	CategoryID   int64  `gorm:"column:category_id;type:bigint;" json:"category_id" form:"category_id"`            // 主键
	ParentID     int64  `gorm:"column:parent_id;type:bigint" json:"parent_id" form:"parent_id"`                   // 父节点
	SiteID       int64  `gorm:"column:site_id;type:bigint" json:"site_id" form:"site_id"`                         // 所属站点
	ChannelID    int64  `gorm:"column:channel_id;type:bigint" json:"channel_id" form:"channel_id"`                // 所属频道
	Title        string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                         // 类别标题
	CallIndex    string `gorm:"column:call_index;type:varchar(64)" json:"call_index" form:"call_index"`           // 调用别名
	ClassLayer   int32  `gorm:"column:class_layer;type:int;default:1" json:"class_layer" form:"class_layer"`      // 类别深度
	LinkURL      string `gorm:"column:link_url;type:varchar(256)" json:"link_url" form:"link_url"`                // 外部链接
	ImgURL       string `gorm:"column:img_url;type:varchar(256)" json:"img_url" form:"img_url"`                   // 图片地址
	SortID       int32  `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                            // 排序
	IsShow       bool   `gorm:"column:is_show;type:tinyint(1);default:1" json:"is_show" form:"is_show"`           // 是否显示:1显示，0隐藏
	IsSearch     bool   `gorm:"column:is_search;type:tinyint(1);default:1" json:"is_search" form:"is_search"`     // 允许检索:1允许，0禁止
	IsDeleted    bool   `gorm:"column:is_deleted;type:tinyint(1)" json:"is_deleted" form:"is_deleted"`            // 删除标识
}

// ApiArticleModel 文章模型
type ApiCategoryOneModel struct {
	ApiCategoryFindModel
	SeoTitle       string `gorm:"column:seo_title;type:varchar(128)" json:"seo_title" form:"seo_title"`                   // SEO标题
	SeoKeyword     string `gorm:"column:seo_keyword;type:varchar(128)" json:"seo_keyword" form:"seo_keyword"`             // SEO关健字
	SeoDescription string `gorm:"column:seo_description;type:varchar(128)" json:"seo_description" form:"seo_description"` // SEO描述
	Content        string `gorm:"column:content;type:text" json:"content" form:"content"`                                 // 内容介绍
}

// ApiArticleListModel 文章查询模型
type ApiArticleListModel struct {
	ArtilceID      int64     `gorm:"column:article_id;type:bigint;" json:"article_id" form:"article_id"`                     // 主键
	SiteID         int64     `gorm:"column:site_id;type:bigint" json:"site_id" form:"site_id"`                               // 所属站点
	ChannelID      int64     `gorm:"column:channel_id;type:bigint" json:"channel_id" form:"channel_id"`                      // 所属频道
	CategoryID     int64     `gorm:"column:category_id;type:bigint" json:"category_id" form:"category_id"`                   // 所属类别
	Title          string    `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                               // 文章标题
	SubTitle       string    `gorm:"column:sub_title;type:varchar(128)" json:"sub_title" form:"sub_title"`                   // 副标题
	IconURL        string    `gorm:"column:icon_url;type:varchar(256)" json:"icon_url" form:"icon_url"`                      // 图标地址
	CallIndex      string    `gorm:"column:call_index;type:varchar(64)" json:"call_index" form:"call_index"`                 // 调用别名
	Source         string    `gorm:"column:source;type:varchar(64)" json:"source" form:"source"`                             // 文章来源
	Author         string    `gorm:"column:author;type:varchar(64)" json:"author" form:"author"`                             // 作者
	LinkURL        string    `gorm:"column:link_url;type:varchar(256)" json:"link_url" form:"link_url"`                      // 外部链接
	ImgURL         string    `gorm:"column:img_url;type:varchar(256)" json:"img_url" form:"img_url"`                         // 图片地址
	SeoTitle       string    `gorm:"column:seo_title;type:varchar(128)" json:"seo_title" form:"seo_title"`                   // SEO标题
	SeoKeyword     string    `gorm:"column:seo_keyword;type:varchar(128)" json:"seo_keyword" form:"seo_keyword"`             // SEO关健字
	SeoDescription string    `gorm:"column:seo_description;type:varchar(128)" json:"seo_description" form:"seo_description"` // SEO描述
	Tags           string    `gorm:"column:tags;type:varchar(128)" json:"tags" form:"tags"`                                  // 标签
	Summary        string    `gorm:"column:summary;type:varchar(128)" json:"summary" form:"summary"`                         // 摘要
	Click          int32     `gorm:"column:click;type:int" json:"click" form:"click"`                                        // 点击数
	IsLock         bool      `gorm:"column:is_lock;type:tinyint(1)" json:"is_lock" form:"is_lock"`                           // 是否锁定
	IsComment      bool      `gorm:"column:is_comment;type:tinyint(1)" json:"is_comment" form:"is_comment"`                  // 是否允许评论
	LikeCount      int32     `gorm:"column:like_count;type:int" json:"like_count" form:"like_count"`                         // 点赞数
	IsTop          bool      `gorm:"column:is_top;type:tinyint(1)" json:"is_top" form:"is_top"`                              // 是否置顶
	IsHot          bool      `gorm:"column:is_hot;type:tinyint(1)" json:"is_hot" form:"is_hot"`                              // 是否热门
	IsSlide        bool      `gorm:"column:is_slide;type:tinyint(1)" json:"is_slide" form:"is_slide"`                        // 是否幻灯片
	StaticURL      string    `gorm:"column:static_url;type:varchar(256)" json:"static_url" form:"static_url"`                // 静态化地址
	PublishTime    time.Time `gorm:"column:publish_time;type:datetime" json:"publish_time" form:"publish_time"`              // 发布时间
}

// ApiCategoryNav 导航栏目
type ApiCategoryNav struct {
	Title     string `json:"title"`
	CallIndex string `json:"call_index"`
	LinkURL   string `json:"link_url"`
	NavType   string `json:"nav_type"`
}
