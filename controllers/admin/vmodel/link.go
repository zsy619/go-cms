package vmodel

type Link_SaveSortIdModel struct {
	LinkId int64 `json:"link_id"`
	SortId int   `json:"sort_id"`
}

type Category_SaveSortIdModel struct {
	CategoryId int64 `json:"category_id"`
	SortId     int   `json:"sort_id"`
}
