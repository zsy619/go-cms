package vmodel

type Site_SaveSortIdModel struct {
	SiteId int64 `json:"site_id"`
	SortId int32 `json:"sort_id"`
}

type Channel_SaveSortIdModel struct {
	ChannelID int64 `json:"channel_id"`
	SortId    int32 `json:"sort_id"`
}
