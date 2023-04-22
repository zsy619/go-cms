package bizmodel

import "time"

const ApiArticleListModel_Field = `a.article_id,a.site_id,a.channel_id,a.category_id,a.title,a.sub_title,a.ico_url1,a.ico_url2,a.call_index,a.source,a.author,a.link_url,a.img_url1,a.img_url2,a.seo_title,a.seo_keyword,a.seo_description,a.tags,a.summary,a.click,a.is_lock,a.is_comment,a.like_count,a.is_top,a.is_hot,a.is_slide,a.static_url,a.publish_time` +
	`,b.call_index as category_call_index,b.title as category_title,b.link_url as category_link_url` +
	`,c.name as channel_name,c.title as channel_title`

const ApiArticleListModel_Table = `SELECT ` + ApiArticleListModel_Field + ` FROM cms_article a` +
	` LEFT JOIN cms_article_category b ON a.category_id=b.category_id` +
	` LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id` +
	` WHERE a.status=2 AND b.status=2`

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
	ImgUrl1      string `gorm:"column:img_url1;type:varchar(256)" json:"img_url1" form:"img_url1"`                // 图片地址
	ImgUrl2      string `gorm:"column:img_url2;type:varchar(256)" json:"img_url2" form:"img_url2"`                // 图片地址
	SortID       int32  `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                            // 排序
	IsShow       bool   `gorm:"column:is_show;type:tinyint(1);default:1" json:"is_show" form:"is_show"`           // 是否显示:1显示，0隐藏
	IsSearch     bool   `gorm:"column:is_search;type:tinyint(1);default:1" json:"is_search" form:"is_search"`     // 允许检索:1允许，0禁止
	IsDeleted    bool   `gorm:"column:is_deleted;type:tinyint(1)" json:"is_deleted" form:"is_deleted"`            // 删除标识
}

// ApiArticleModel 文章模型
type ApiCategoryOneModel struct {
	ApiCategoryFindModel
	Target         string `gorm:"column:target;type:varchar(16);default:_blank;comment:是否开启浏览器新窗口" json:"target" form:"target"`
	SeoTitle       string `gorm:"column:seo_title;type:varchar(128)" json:"seo_title" form:"seo_title"`                   // SEO标题
	SeoKeyword     string `gorm:"column:seo_keyword;type:varchar(128)" json:"seo_keyword" form:"seo_keyword"`             // SEO关健字
	SeoDescription string `gorm:"column:seo_description;type:varchar(128)" json:"seo_description" form:"seo_description"` // SEO描述
	Content        string `gorm:"column:content;type:text" json:"content" form:"content"`                                 // 内容介绍
	Template       string `gorm:"column:template;type:varchar(256);comment:模板路径" json:"template" form:"template"`
}

// ApiArticleListModel 文章查询模型
type ApiArticleListModel struct {
	ArtilceID         int64     `gorm:"column:article_id;type:bigint;" json:"article_id" form:"article_id"`                                // 主键
	SiteID            int64     `gorm:"column:site_id;type:bigint" json:"site_id" form:"site_id"`                                          // 所属站点
	ChannelID         int64     `gorm:"column:channel_id;type:bigint" json:"channel_id" form:"channel_id"`                                 // 所属频道
	CategoryID        int64     `gorm:"column:category_id;type:bigint" json:"category_id" form:"category_id"`                              // 所属类别
	Title             string    `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                                          // 文章标题
	SubTitle          string    `gorm:"column:sub_title;type:varchar(128)" json:"sub_title" form:"sub_title"`                              // 副标题
	IcoURL1           string    `gorm:"column:ico_url1;type:varchar(256)" json:"ico_url1" form:"ico_url1"`                                 // 图标地址
	IcoURL2           string    `gorm:"column:ico_url2;type:varchar(256)" json:"ico_url2" form:"ico_url2"`                                 // 图标地址
	CallIndex         string    `gorm:"column:call_index;type:varchar(64)" json:"call_index" form:"call_index"`                            // 调用别名
	Source            string    `gorm:"column:source;type:varchar(64)" json:"source" form:"source"`                                        // 文章来源
	Author            string    `gorm:"column:author;type:varchar(64)" json:"author" form:"author"`                                        // 作者
	LinkURL           string    `gorm:"column:link_url;type:varchar(256)" json:"link_url" form:"link_url"`                                 // 外部链接
	ImgURL1           string    `gorm:"column:img_url1;type:varchar(256)" json:"img_url1" form:"img_url1"`                                 // 图片地址
	ImgURL2           string    `gorm:"column:img_url2;type:varchar(256)" json:"img_url2" form:"img_url2"`                                 // 图片地址
	SeoTitle          string    `gorm:"column:seo_title;type:varchar(128)" json:"seo_title" form:"seo_title"`                              // SEO标题
	SeoKeyword        string    `gorm:"column:seo_keyword;type:varchar(128)" json:"seo_keyword" form:"seo_keyword"`                        // SEO关健字
	SeoDescription    string    `gorm:"column:seo_description;type:varchar(128)" json:"seo_description" form:"seo_description"`            // SEO描述
	Tags              string    `gorm:"column:tags;type:varchar(128)" json:"tags" form:"tags"`                                             // 标签
	Summary           string    `gorm:"column:summary;type:varchar(128)" json:"summary" form:"summary"`                                    // 摘要
	Click             int32     `gorm:"column:click;type:int" json:"click" form:"click"`                                                   // 点击数
	IsLock            bool      `gorm:"column:is_lock;type:tinyint(1)" json:"is_lock" form:"is_lock"`                                      // 是否锁定
	IsComment         bool      `gorm:"column:is_comment;type:tinyint(1)" json:"is_comment" form:"is_comment"`                             // 是否允许评论
	LikeCount         int32     `gorm:"column:like_count;type:int" json:"like_count" form:"like_count"`                                    // 点赞数
	IsTop             bool      `gorm:"column:is_top;type:tinyint(1)" json:"is_top" form:"is_top"`                                         // 是否置顶
	IsHot             bool      `gorm:"column:is_hot;type:tinyint(1)" json:"is_hot" form:"is_hot"`                                         // 是否热门
	IsSlide           bool      `gorm:"column:is_slide;type:tinyint(1)" json:"is_slide" form:"is_slide"`                                   // 是否幻灯片
	StaticURL         string    `gorm:"column:static_url;type:varchar(256)" json:"static_url" form:"static_url"`                           // 静态化地址
	PublishTime       time.Time `gorm:"column:publish_time;type:datetime" json:"publish_time" form:"publish_time"`                         // 发布时间
	CategoryTitle     string    `gorm:"column:category_title;type:varchar(128)" json:"category_title" form:"category_title"`               // category_类别标题
	CategoryCallIndex string    `gorm:"column:category_call_index;type:varchar(64)" json:"category_call_index" form:"category_call_index"` // category_调用别名
	CategoryLinkURL   string    `gorm:"column:category_link_url;type:varchar(256)" json:"category_link_url" form:"link_url"`               // category_外部链接
	ChannelName       string    `gorm:"column:channel_name;type:varchar(128)" json:"channel_name" form:"channel_name"`                     // channel_字段名
	ChannelTitle      string    `gorm:"column:channel_title;type:varchar(128)" json:"channel_title" form:"channel_title"`                  // channel_标题
}

// ApiCategoryNav 导航栏目
type ApiCategoryNav struct {
	Title     string `json:"title"`
	CallIndex string `json:"call_index"`
	LinkURL   string `json:"link_url"`
	NavType   string `json:"nav_type"`
}

// ApiArticleOneModel 文章模型
type ApiArticleOneModel struct {
	ArticleID         int64     `gorm:"column:article_id;type:bigint;" json:"article_id" form:"article_id"`                                  // 主键
	SiteID            int64     `gorm:"column:site_id;type:bigint" json:"site_id" form:"site_id"`                                            // 所属站点
	ChannelID         int64     `gorm:"column:channel_id;type:bigint" json:"channel_id" form:"channel_id"`                                   // 所属频道
	CategoryID        int64     `gorm:"column:category_id;type:bigint" json:"category_id" form:"category_id"`                                // 类别ID
	Title             string    `gorm:"column:title;type:varchar(256)" json:"title" form:"title"`                                            // 内容标题
	SubTitle          string    `gorm:"column:sub_title;type:varchar(128)" json:"sub_title" form:"sub_title"`                                // 副标题
	IcoUrl1           string    `gorm:"column:ico_url1;type:varchar(256)" json:"ico_url1" form:"ico_url1"`                                   // 图标地址
	IcoUrl2           string    `gorm:"column:ico_url2;type:varchar(256)" json:"ico_url2" form:"ico_url2"`                                   // 图标地址
	CallIndex         string    `gorm:"column:call_index;type:varchar(64)" json:"call_index" form:"call_index"`                              // 调用别名
	Source            string    `gorm:"column:source;type:varchar(64)" json:"source" form:"source"`                                          // 来源
	Author            string    `gorm:"column:author;type:varchar(64)" json:"author" form:"author"`                                          // 作者
	LinkURL           string    `gorm:"column:link_url;type:varchar(256)" json:"link_url" form:"link_url"`                                   // 外部链接
	ImgUrl1           string    `gorm:"column:img_url1;type:varchar(256)" json:"img_url1" form:"img_url1"`                                   // 图片地址
	ImgUrl2           string    `gorm:"column:img_url2;type:varchar(256)" json:"img_url2" form:"img_url2"`                                   // 图片地址
	SeoTitle          string    `gorm:"column:seo_title;type:varchar(128)" json:"seo_title" form:"seo_title"`                                // SEO标题
	SeoKeyword        string    `gorm:"column:seo_keyword;type:varchar(128)" json:"seo_keyword" form:"seo_keyword"`                          // SEO关健字
	SeoDescription    string    `gorm:"column:seo_description;type:varchar(128)" json:"seo_description" form:"seo_description"`              // SEO描述
	Tags              string    `gorm:"column:tags;type:varchar(128)" json:"tags" form:"tags"`                                               // TAG标签逗号分隔
	Summary           string    `gorm:"column:summary;type:varchar(128)" json:"summary" form:"summary"`                                      // 内容摘要
	Content           string    `gorm:"column:content;type:text" json:"content" form:"content"`                                              // 详细内容
	SortID            int32     `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                                               // 排序
	Click             int32     `gorm:"column:click;type:int" json:"click" form:"click"`                                                     // 浏览次数
	IsLock            int32     `gorm:"column:is_lock;type:tinyint" json:"is_lock" form:"is_lock"`                                           // 是否锁定（不允许编辑）
	IsComment         int32     `gorm:"column:is_comment;type:tinyint" json:"is_comment" form:"is_comment"`                                  // 是否允许评论:0禁止1允许
	CommentCount      int32     `gorm:"column:comment_count;type:int" json:"comment_count" form:"comment_count"`                             // 评论总数
	LikeCount         int32     `gorm:"column:like_count;type:int" json:"like_count" form:"like_count"`                                      // 点赞总数
	IsTop             int32     `gorm:"column:is_top;type:tinyint" json:"is_top" form:"is_top"`                                              // 是否置顶
	IsRed             int32     `gorm:"column:is_red;type:tinyint" json:"is_red" form:"is_red"`                                              // 是否推荐
	IsHot             int32     `gorm:"column:is_hot;type:tinyint" json:"is_hot" form:"is_hot"`                                              // 是否热门
	IsSlide           int32     `gorm:"column:is_slide;type:tinyint" json:"is_slide" form:"is_slide"`                                        // 是否幻灯片
	StaticURL         string    `gorm:"column:static_url;type:varchar(256)" json:"static_url" form:"static_url"`                             // 静态链接
	PublishTime       time.Time `gorm:"column:publish_time;type:datetime;default:CURRENT_TIMESTAMP" json:"publish_time" form:"publish_time"` // 发布时间
	CategoryTitle     string    `gorm:"column:category_title;type:varchar(128)" json:"category_title" form:"category_title"`                 // category_类别标题
	CategoryCallIndex string    `gorm:"column:category_call_index;type:varchar(64)" json:"category_call_index" form:"category_call_index"`   // category_调用别名
	CategoryLinkURL   string    `gorm:"column:category_link_url;type:varchar(256)" json:"category_link_url" form:"link_url"`                 // category_外部链接
	ChannelName       string    `gorm:"column:channel_name;type:varchar(128)" json:"channel_name" form:"channel_name"`                       // channel_字段名
	ChannelTitle      string    `gorm:"column:channel_title;type:varchar(128)" json:"channel_title" form:"channel_title"`                    // channel_标题
}
