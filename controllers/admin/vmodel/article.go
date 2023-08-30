package vmodel

type Article_SaveSortIdModel struct {
	ArticleId int64 `json:"article_id"`
	SortId    int   `json:"sort_id"`
}

type Article_ChangeStatusModel struct {
	ArticleIds []int64 `json:"article_ids"`
	Status     int32   `json:"status"`
}

type Article_AlbumSaveShowModel struct {
	AlbumIds []int64 `json:"album_ids"`
	Show     int32   `json:"show"`
}

type Article_AlbumSaveBatchdModel struct {
	AlbumID int64  `gorm:"column:album_id;type:bigint;" json:"album_id" form:"album_id"`      // 主键
	Title   string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`          // 标题
	LinkURL string `gorm:"column:link_url;type:varchar(256)" json:"link_url" form:"link_url"` // 外部链接
	Click   int32  `gorm:"column:click;type:int" json:"click" form:"click"`                   // 点击次数
	SortID  int32  `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`             // 排序
	IsShow  int32  `gorm:"column:is_show;type:tinyint" json:"is_show" form:"is_show"`         // 是否显示：1显示2隐藏
	Remark  string `gorm:"column:remark;type:varchar(256)" json:"remark" form:"remark"`       // 图片描述
}

type Article_AttachSaveShowModel struct {
	AttachIds []int64 `json:"attach_ids"`
	Show      int32   `json:"show"`
}

type Article_AttachSaveBatchdModel struct {
	AttachID int64  `gorm:"column:attach_id;type:bigint;" json:"attach_id" form:"attach_id"` // 主键
	Title    string `gorm:"column:title;type:varchar(128)" json:"title" form:"title"`        // 标题
	Point    int32  `gorm:"column:point;type:int" json:"point" form:"point"`                 // 下载所需积分
	Click    int32  `gorm:"column:click;type:int" json:"click" form:"click"`                 // 下载次数
	SortID   int32  `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`           // 排序
	Remark   string `gorm:"column:remark;type:varchar(256)" json:"remark" form:"remark"`     // 图片描述
}

type Comment_ChangeStatusModel struct {
	CommentIds []int64 `json:"comment_ids"`
	Status     int32   `json:"status"`
}

type Property_ChangeStatusModel struct {
	PropertyIds []int64 `json:"property_ids"`
	Status      int32   `json:"status"`
}

type Property_SaveSortIdModel struct {
	PropertyId int64 `json:"property_id"`
	SortId     int   `json:"sort_id"`
}
