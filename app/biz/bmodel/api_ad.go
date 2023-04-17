package bmodel

import "time"

type ApiAdListModel struct {
	AdID          int64     `gorm:"column:ad_id;type:bigint;" json:"ad_id" form:"ad_id"`                                           // 主键
	SiteID        int64     `gorm:"column:site_id;type:bigint" json:"site_id" form:"site_id"`                                      // 所属站点
	ChannelID     int64     `gorm:"column:channel_id;type:bigint" json:"channel_id" form:"channel_id"`                             // 所属频道
	CategoryID    int64     `gorm:"column:category_id;type:bigint" json:"category_id" form:"category_id"`                          // 类别ID
	CategoryTitle string    `gorm:"column:category_title;type:varchar(128)" json:"category_title" form:"category_title"`           // 列表标题
	Title         string    `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`                                      // 标题
	CallIndex     string    `gorm:"column:call_index;type:varchar(64)" json:"call_index" form:"call_index"`                        // 调用别名
	LinkURL       string    `gorm:"column:link_url;type:varchar(256)" json:"link_url" form:"link_url"`                             // 外部链接
	Target        string    `gorm:"column:target;type:varchar(16);default:_blank" json:"target" form:"target"`                     // 是否开启浏览器新窗口
	ImgUrl1       string    `gorm:"column:img_url1;type:varchar(256)" json:"img_url1" form:"img_url1"`                             // 图片地址
	ImgUrl2       string    `gorm:"column:img_url2;type:varchar(256)" json:"img_url2" form:"img_url2"`                             // 图片地址
	Remark        string    `gorm:"column:remark;type:varchar(256)" json:"remark" form:"remark"`                                   // 备注
	SortID        int32     `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                                         // 排序
	Click         int32     `gorm:"column:click;type:int" json:"click" form:"click"`                                               // 浏览次数
	Status        int32     `gorm:"column:status;type:tinyint" json:"status" form:"status"`                                        // 状态0草稿1提交2审核通过3审核未通过4驳回
	IsLock        int32     `gorm:"column:is_lock;type:tinyint" json:"is_lock" form:"is_lock"`                                     // 是否锁定（不允许编辑）
	IsTop         int32     `gorm:"column:is_top;type:tinyint" json:"is_top" form:"is_top"`                                        // 是否置顶
	IsRed         int32     `gorm:"column:is_red;type:tinyint" json:"is_red" form:"is_red"`                                        // 是否推荐
	IsHot         int32     `gorm:"column:is_hot;type:tinyint" json:"is_hot" form:"is_hot"`                                        // 是否热门
	IsSlide       int32     `gorm:"column:is_slide;type:tinyint" json:"is_slide" form:"is_slide"`                                  // 是否幻灯片
	BeginTime     time.Time `gorm:"column:begin_time;type:datetime;default:CURRENT_TIMESTAMP" json:"begin_time" form:"begin_time"` // 创建时间
	EndTime       time.Time `gorm:"column:end_time;type:datetime;default:CURRENT_TIMESTAMP" json:"end_time" form:"end_time"`       // 创建时间
}
