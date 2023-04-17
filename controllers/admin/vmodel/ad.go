package vmodel

type Ad_SaveSortIdModel struct {
	AdId   int64 `json:"ad_id"`
	SortId int   `json:"sort_id"`
}

type Ad_ChangeStatusModel struct {
	AdIds  []int64 `json:"ad_ids"`
	Status int32   `json:"status"`
}
