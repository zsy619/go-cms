package bizmodel

type ApiChannelModel struct {
	ChannelID  int64  `gorm:"column:channel_id;type:bigint;primaryKey;" json:"channel_id" form:"channel_id"` // 主键
	ParentID   int64  `gorm:"column:parent_id;type:bigint" json:"parent_id" form:"parent_id"`                // 父级ID
	Title      string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                      // 标题
	Name       string `gorm:"column:name;type:varchar(128)" json:"name" form:"name"`                         // 频道名称
	Kind       int32  `gorm:"column:kind;type:tinyint" json:"kind" form:"kind"`                              // 频道类型
	ClassLayer int32  `gorm:"column:class_layer;type:tinyint" json:"class_layer" form:"class_layer"`         // 层级
	LinkURL    string `gorm:"column:link_url;type:varchar(256)" json:"link_url" form:"link_url"`             // 外部链接
	ImgUrl1    string `gorm:"column:img_url1;type:varchar(256)" json:"img_url1" form:"img_url1"`             // 图片地址
	ImgUrl2    string `gorm:"column:img_url2;type:varchar(256)" json:"img_url2" form:"img_url2"`             // 图片地址
	SortID     int32  `gorm:"column:sort_id;type:tinyint" json:"sort_id" form:"sort_id"`                     // 排序
	IsAlbum    int32  `gorm:"column:is_album;type:tinyint" json:"is_album" form:"is_album"`                  // 是否相册
	IsAttach   int32  `gorm:"column:is_attach;type:tinyint" json:"is_attach" form:"is_attach"`               // 是否附件
	IsSpec     int32  `gorm:"column:is_spec;type:tinyint" json:"is_spec" form:"is_spec"`                     // 是否规格
	TmplChnl   string `gorm:"column:tmpl_chnl;type:varchar(256)" json:"tmpl_chnl" form:"tmpl_chnl"`          // 频道模板路径
	TmplCat    string `gorm:"column:tmpl_cat;type:varchar(256)" json:"tmpl_cat" form:"tmpl_cat"`             // 栏目模板路径
	TmplLst    string `gorm:"column:tmpl_lst;type:varchar(256)" json:"tmpl_lst" form:"tmpl_lst"`             // 列表模板路径
	TmplDtl    string `gorm:"column:tmpl_dtl;type:varchar(256)" json:"tmpl_dtl" form:"tmpl_dtl"`             // 明细模板路径
	SiteFlag   string `gorm:"column:site_flag;type:varchar(64)" json:"site_flag" form:"site_flag"`           // 站点标识
}

type ApiNavModel struct {
	NavID    int64          `gorm:"column:nav_id;type:bigint;" json:"nav_id" form:"nav_id"`             // 主键
	Title    string         `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`           // 标题
	Name     string         `gorm:"column:name;type:varchar(128)" json:"name" form:"name"`              // 名称
	LinkURL  string         `gorm:"column:link_url;type:varchar(256);" json:"link_url" form:"link_url"` // 外部链接
	Target   string         `gorm:"column:target;type:varchar(16)" json:"target" form:"target"`         // 是否开启浏览器新窗口
	ImgUrl1  string         `gorm:"column:img_url1;type:varchar(256)" json:"img_url1" form:"img_url1"`  // 图片地址
	ImgUrl2  string         `gorm:"column:img_url2;type:varchar(256)" json:"img_url2" form:"img_url2"`  // 图片地址
	SortID   int32          `gorm:"column:sort_id;type:tinyint" json:"sort_id" form:"sort_id"`          // 排序
	Type     string         `gorm:"-" json:"type" form:"type"`                                          // 类型 channel:频道 category:分类
	Children []*ApiNavModel `gorm:"-" json:"children"`                                                  // 子级
}

type ApiSiteModel struct {
	SiteID          int64  `gorm:"column:site_id;type:bigint;" json:"site_id" form:"site_id"`                                             // 站点ID
	ParentID        int64  `gorm:"column:parent_id;type:bigint;not null;" json:"parent_id" form:"parent_id"`                              // 父级ID
	Name            string `gorm:"column:name;type:varchar(128);" json:"name" form:"name"`                                                // 站点名称
	Flag            string `gorm:"column:flag;type:varchar(64);" json:"flag" form:"flag"`                                                 // 站点标识
	Title           string `gorm:"column:title;type:varchar(128);" json:"title" form:"title"`                                             // 标题
	Template        string `gorm:"column:template;type:varchar(128);" json:"template" form:"template"`                                    // 模板名称
	IsDefault       bool   `gorm:"column:is_default;type:tinyint(1);" json:"is_default" form:"is_default"`                                // 是否默认站
	IsMobile        bool   `gorm:"column:is_mobile;type:tinyint(1);" json:"is_mobile" form:"is_mobile"`                                   // 是否移动端
	Logo1           string `gorm:"column:logo1;type:varchar(512);" json:"logo1" form:"logo1"`                                             // 网站LOGO
	Logo2           string `gorm:"column:logo2;type:varchar(512);" json:"logo2" form:"logo2"`                                             // 网站LOGO
	Icon1           string `gorm:"column:icon1;type:varchar(512);" json:"icon1" form:"icon1"`                                             // 网站icon
	Icon2           string `gorm:"column:icon2;type:varchar(512);" json:"icon2" form:"icon2"`                                             // 网站icon
	Company         string `gorm:"column:company;type:varchar(512);" json:"company" form:"company"`                                       // 公司名称
	Address         string `gorm:"column:address;type:varchar(512);" json:"address" form:"address"`                                       // 通讯地址
	Telphone        string `gorm:"column:telphone;type:varchar(64);" json:"telphone" form:"telphone"`                                     // 联系电话
	Fax             string `gorm:"column:fax;type:varchar(64);" json:"fax" form:"fax"`                                                    // 传真
	Email           string `gorm:"column:email;type:varchar(64);" json:"email" form:"email"`                                              // 邮箱
	Crod            string `gorm:"column:crod;type:varchar(64);" json:"crod" form:"crod"`                                                 // 备案号
	Cache           int32  `gorm:"column:cache;type:int;" json:"cache" form:"cache"`                                                      // 缓存时间
	MaxLength       int32  `gorm:"column:max_length;type:int;default:2048;" json:"max_length" form:"max_length"`                          // 最大文件上传
	FileType        string `gorm:"column:file_type;type:varchar(64);default:png|gif|jpg|jpeg|zip|rar;" json:"file_type" form:"file_type"` // 上传文件类型
	HomeTitle       string `gorm:"column:home_title;type:varchar(128);" json:"home_title" form:"home_title"`                              // 首页标题
	Copyright       string `gorm:"column:copyright;type:varchar(128);" json:"copyright" form:"copyright"`                                 // 版权信息
	Statcode        string `gorm:"column:statcode;type:varchar(128);" json:"statcode" form:"statcode"`                                    // 统计代码
	Robots          string `gorm:"column:robots;type:varchar(128);" json:"robots" form:"robots"`                                          // 爬虫规则
	MetaKeyword     string `gorm:"column:meta_keyword;type:varchar(512);" json:"meta_keyword" form:"meta_keyword"`                        // META关键词
	MetaDescription string `gorm:"column:meta_description;type:varchar(128);" json:"meta_description" form:"meta_description"`            // META描述
}
