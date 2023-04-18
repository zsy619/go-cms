package bizmodel

import (
	"time"

	"haedu.gov.cn/cms/app/dal/model"
)

type Weixin_ContentSubscribeOrDefaultModel struct {
	AccountID   int64                         `json:"account_id" form:"account_id"`     // 归属公众号
	RequestType int32                         `json:"request_type" form:"request_type"` // 请求类型(0默认回复1文字2图片3语音4链接5地理位置6关注7取消关注8扫描带参数二维码事件9上报地理位置事件10自定义菜单事件）
	TextReply   *model.WeixinRequestContent   `json:"text_reply" form:"text_reply"`     // 文本回复
	ImageReply  []*model.WeixinRequestContent `json:"image_reply" form:"image_reply"`   // 图文回复
	SoundReply  *model.WeixinRequestContent   `json:"sound_reply" form:"sound_reply"`   // 语音回复
}

type Weixin_RuleModel struct {
	RuleID       int64     `gorm:"column:rule_id;type:bigint" json:"rule_id" form:"rule_id"`                                  // 主键
	AccountID    int64     `gorm:"column:account_id;type:bigint;not null" json:"account_id" form:"account_id"`                // 归属公众号
	Name         string    `gorm:"column:name;type:varchar(256)" json:"name" form:"name"`                                     // 规则名称
	Keywords     string    `gorm:"column:keywords;type:varchar(2048)" json:"keywords" form:"keywords"`                        // 请求关键词,逗号分隔
	RequestType  int32     `gorm:"column:request_type;type:int" json:"request_type" form:"request_type"`                      // 请求类型(0默认回复1文字2图片3语音4链接5地理位置6关注7取消关注8扫描带参数二维码事件9上报地理位置事件10自定义菜单事件）
	ResponseType int32     `gorm:"column:response_type;type:int" json:"response_type" form:"response_type"`                   // 回复类型(1文本2图文3语音4视频5第三方接口)
	IsLikeQuery  int32     `gorm:"column:is_like_query;type:tinyint" json:"is_like_query" form:"is_like_query"`               // 是否模糊查询
	IsDefault    int32     `gorm:"column:is_default;type:tinyint" json:"is_default" form:"is_default"`                        // 是否默认回复
	SortID       int32     `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                                     // 排序
	CreateTime   time.Time `gorm:"column:create_time;type:int unsigned;autoCreateTime" json:"create_time" form:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;type:int unsigned;autoUpdateTime" json:"update_time" form:"update_time"` // 修改时间
	ContentID    int64     `gorm:"column:content_id;type:bigint" json:"content_id" form:"content_id"`                         // 主键
	Title        string    `gorm:"column:title;type:varchar(512)" json:"title" form:"title"`                                  // 回复标题
	Content      string    `gorm:"column:content;type:text" json:"content" form:"content"`                                    // 回复内容
	LinkURL      string    `gorm:"column:link_url;type:varchar(512)" json:"link_url" form:"link_url"`                         // 详情链接地址
	ImgURL       string    `gorm:"column:img_url;type:varchar(512)" json:"img_url" form:"img_url"`                            // 图片地址
	MediaURL     string    `gorm:"column:media_url;type:varchar(512)" json:"media_url" form:"media_url"`                      // 语音或视频地址
	MediaHdURL   string    `gorm:"column:media_hd_url;type:varchar(512)" json:"media_hd_url" form:"media_hd_url"`             // 高清语音或者视频地址
}

type Weixin_PictureModel struct {
	RuleID       int64                         `gorm:"column:rule_id;type:bigint" json:"rule_id" form:"rule_id"`                                  // 主键
	AccountID    int64                         `gorm:"column:account_id;type:bigint;not null" json:"account_id" form:"account_id"`                // 归属公众号
	Name         string                        `gorm:"column:name;type:varchar(256)" json:"name" form:"name"`                                     // 规则名称
	Keywords     string                        `gorm:"column:keywords;type:varchar(2048)" json:"keywords" form:"keywords"`                        // 请求关键词,逗号分隔
	RequestType  int32                         `gorm:"column:request_type;type:int" json:"request_type" form:"request_type"`                      // 请求类型(0默认回复1文字2图片3语音4链接5地理位置6关注7取消关注8扫描带参数二维码事件9上报地理位置事件10自定义菜单事件）
	ResponseType int32                         `gorm:"column:response_type;type:int" json:"response_type" form:"response_type"`                   // 回复类型(1文本2图文3语音4视频5第三方接口)
	IsLikeQuery  int32                         `gorm:"column:is_like_query;type:tinyint" json:"is_like_query" form:"is_like_query"`               // 是否模糊查询
	IsDefault    int32                         `gorm:"column:is_default;type:tinyint" json:"is_default" form:"is_default"`                        // 是否默认回复
	SortID       int32                         `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                                     // 排序
	CreateTime   time.Time                     `gorm:"column:create_time;type:int unsigned;autoCreateTime" json:"create_time" form:"create_time"` // 创建时间
	UpdateTime   time.Time                     `gorm:"column:update_time;type:int unsigned;autoUpdateTime" json:"update_time" form:"update_time"` // 修改时间
	ImageReply   []*model.WeixinRequestContent `json:"image_reply" form:"image_reply"`
}

type Weixin_RuleCountModel struct {
	RuleID       int64     `gorm:"column:rule_id;type:bigint" json:"rule_id" form:"rule_id"`                                  // 主键
	AccountID    int64     `gorm:"column:account_id;type:bigint;not null" json:"account_id" form:"account_id"`                // 归属公众号
	Name         string    `gorm:"column:name;type:varchar(256)" json:"name" form:"name"`                                     // 规则名称
	Keywords     string    `gorm:"column:keywords;type:varchar(2048)" json:"keywords" form:"keywords"`                        // 请求关键词,逗号分隔
	RequestType  int32     `gorm:"column:request_type;type:int" json:"request_type" form:"request_type"`                      // 请求类型(0默认回复1文字2图片3语音4链接5地理位置6关注7取消关注8扫描带参数二维码事件9上报地理位置事件10自定义菜单事件）
	ResponseType int32     `gorm:"column:response_type;type:int" json:"response_type" form:"response_type"`                   // 回复类型(1文本2图文3语音4视频5第三方接口)
	IsLikeQuery  int32     `gorm:"column:is_like_query;type:tinyint" json:"is_like_query" form:"is_like_query"`               // 是否模糊查询
	IsDefault    int32     `gorm:"column:is_default;type:tinyint" json:"is_default" form:"is_default"`                        // 是否默认回复
	SortID       int32     `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                                     // 排序
	CreateTime   time.Time `gorm:"column:create_time;type:int unsigned;autoCreateTime" json:"create_time" form:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;type:int unsigned;autoUpdateTime" json:"update_time" form:"update_time"` // 修改时间
	Count        int32     `gorm:"column:count;type:int" json:"count" form:"count"`                                           // 汇总数量
}

type Weixin_ContentModel struct {
	ContentID      int64     `gorm:"column:content_id;type:bigint" json:"content_id" form:"content_id"`                         // 主键
	AccountID      int64     `gorm:"column:account_id;type:bigint;not null" json:"account_id" form:"account_id"`                // 归属公众号
	OpenID         string    `gorm:"column:open_id;type:varchar(128)" json:"open_id" form:"open_id"`                            // 用户微信ID
	RequestType    string    `gorm:"column:request_type;type:varchar(32)" json:"request_type" form:"request_type"`              // 数据类型 文本消息：text 图片消息:image 地理位置消息:location 链接消息:link 事件:event
	RequestContent string    `gorm:"column:request_content;type:varchar(2048)" json:"request_content" form:"request_content"`   // 数据内容
	ResponseType   string    `gorm:"column:response_type;type:varchar(32)" json:"response_type" form:"response_type"`           // 回复的类型 文本消息：text 图片消息:image 地理位置消息:location 链接消息:link
	ReponseContent string    `gorm:"column:reponse_content;type:varchar(2048)" json:"reponse_content" form:"reponse_content"`   // 系统回复的内容
	MeidaHdURL     string    `gorm:"column:meida_hd_url;type:varchar(512)" json:"meida_hd_url" form:"meida_hd_url"`             // 高清语音或者视频地址
	CreateTime     time.Time `gorm:"column:create_time;type:int unsigned;autoCreateTime" json:"create_time" form:"create_time"` // 创建时间
	XMLContent     string    `gorm:"column:xml_content;type:varchar(2048)" json:"xml_content" form:"xml_content"`               // xml原始内容
	AddTime        time.Time `gorm:"column:add_time;type:datetime;default:CURRENT_TIMESTAMP" json:"add_time" form:"add_time"`   // 录入系统的时间
	AccountName    string    `gorm:"column:account_name;type:varchar(256)" json:"account_name" form:"account_name"`             // 公众号名称
}
