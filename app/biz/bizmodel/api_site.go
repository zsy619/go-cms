package bizmodel

type ApiChannelFindModel struct {
	ChannelID  int64  `gorm:"column:channel_id;type:bigint;primaryKey;" json:"channel_id" form:"channel_id"` // 主键
	ParentID   int64  `gorm:"column:parent_id;type:bigint" json:"parent_id" form:"parent_id"`                // 父级ID
	Title      string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                      // 标题
	Name       string `gorm:"column:name;type:varchar(128)" json:"name" form:"name"`                         // 频道名称
	Kind       int32  `gorm:"column:kind;type:tinyint" json:"kind" form:"kind"`                              // 频道类型
	ClassLayer int32  `gorm:"column:class_layer;type:tinyint" json:"class_layer" form:"class_layer"`         // 层级
	ImgUrl1    string `gorm:"column:img_url1;type:varchar(256)" json:"img_url1" form:"img_url1"`             // 图片地址
	ImgUrl2    string `gorm:"column:img_url2;type:varchar(256)" json:"img_url2" form:"img_url2"`             // 图片地址
	SortID     int32  `gorm:"column:sort_id;type:tinyint" json:"sort_id" form:"sort_id"`                     // 排序
	IsAlbum    int32  `gorm:"column:is_album;type:tinyint" json:"is_album" form:"is_album"`                  // 是否相册
	IsAttach   int32  `gorm:"column:is_attach;type:tinyint" json:"is_attach" form:"is_attach"`               // 是否附件
	IsSpec     int32  `gorm:"column:is_spec;type:tinyint" json:"is_spec" form:"is_spec"`                     // 是否规格
	Template   string `gorm:"column:template;type:varchar(256)" json:"template" form:"template"`             // 模板路径
}

type ApiNavFindModel struct {
	NavID    int64              `gorm:"column:nav_id;type:bigint;" json:"nav_id" form:"nav_id"`             // 主键
	Title    string             `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`           // 标题
	Name     string             `gorm:"column:name;type:varchar(128)" json:"name" form:"name"`              // 名称
	LinkURL  string             `gorm:"column:link_url;type:varchar(256);" json:"link_url" form:"link_url"` // 外部链接
	Target   string             `gorm:"column:target;type:varchar(16)" json:"target" form:"target"`         // 是否开启浏览器新窗口
	ImgUrl1  string             `gorm:"column:img_url1;type:varchar(256)" json:"img_url1" form:"img_url1"`  // 图片地址
	ImgUrl2  string             `gorm:"column:img_url2;type:varchar(256)" json:"img_url2" form:"img_url2"`  // 图片地址
	SortID   int32              `gorm:"column:sort_id;type:tinyint" json:"sort_id" form:"sort_id"`          // 排序
	Type     string             `gorm:"-" json:"type" form:"type"`                                          // 类型 channel:频道 category:分类
	Children []*ApiNavFindModel `gorm:"-" json:"children"`                                                  // 子级
}

type ApiSiteModel struct {
	SiteID          int64  `gorm:"column:site_id;type:bigint;comment:主键" json:"site_id" form:"site_id"`
	ParentID        int64  `gorm:"column:parent_id;type:bigint;not null;comment:父级ID" json:"parent_id" form:"parent_id"`
	Name            string `gorm:"column:name;type:varchar(128);comment:站点名称" json:"name" form:"name"`
	Flag            string `gorm:"column:flag;type:varchar(64);comment:站点标识" json:"flag" form:"flag"`
	Title           string `gorm:"column:title;type:varchar(128);comment:标题" json:"title" form:"title"`
	Template        string `gorm:"column:template;type:varchar(128);comment:模板名称" json:"template" form:"template"`
	IsDefault       bool   `gorm:"column:is_default;type:tinyint(1);comment:是否默认站" json:"is_default" form:"is_default"`
	IsMobile        bool   `gorm:"column:is_mobile;type:tinyint(1);comment:是否移动端" json:"is_mobile" form:"is_mobile"`
	Logo1           string `gorm:"column:logo1;type:varchar(512);comment:网站LOGO" json:"logo1" form:"logo1"`
	Logo2           string `gorm:"column:logo2;type:varchar(512);comment:网站LOGO" json:"logo2" form:"logo2"`
	Icon1           string `gorm:"column:icon1;type:varchar(512);comment:网站icon" json:"icon1" form:"icon1"`
	Icon2           string `gorm:"column:icon2;type:varchar(512);comment:网站icon" json:"icon2" form:"icon2"`
	Company         string `gorm:"column:company;type:varchar(512);comment:公司名称" json:"company" form:"company"`
	Address         string `gorm:"column:address;type:varchar(512);comment:通讯地址" json:"address" form:"address"`
	Telphone        string `gorm:"column:telphone;type:varchar(64);comment:联系电话" json:"telphone" form:"telphone"`
	Fax             string `gorm:"column:fax;type:varchar(64);comment:传真" json:"fax" form:"fax"`
	Email           string `gorm:"column:email;type:varchar(64);comment:邮箱" json:"email" form:"email"`
	Crod            string `gorm:"column:crod;type:varchar(64);comment:备案号" json:"crod" form:"crod"`
	Cache           int32  `gorm:"column:cache;type:int;comment:缓存时间" json:"cache" form:"cache"`
	MaxLength       int32  `gorm:"column:max_length;type:int;default:2048;comment:最大文件上传" json:"max_length" form:"max_length"`
	FileType        string `gorm:"column:file_type;type:varchar(64);default:png|gif|jpg|jpeg|zip|rar;comment:上传文件类型" json:"file_type" form:"file_type"`
	HomeTitle       string `gorm:"column:home_title;type:varchar(128);comment:首页标题" json:"home_title" form:"home_title"`
	Copyright       string `gorm:"column:copyright;type:varchar(128);comment:版权信息" json:"copyright" form:"copyright"`
	Statcode        string `gorm:"column:statcode;type:varchar(128);comment:统计代码" json:"statcode" form:"statcode"`
	Robots          string `gorm:"column:robots;type:varchar(128);comment:爬虫规则" json:"robots" form:"robots"`
	MetaKeyword     string `gorm:"column:meta_keyword;type:varchar(512);comment:META关键词" json:"meta_keyword" form:"meta_keyword"`
	MetaDescription string `gorm:"column:meta_description;type:varchar(128);comment:META描述" json:"meta_description" form:"meta_description"`
}
