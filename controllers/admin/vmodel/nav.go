package vmodel

type Nav_SaveSortIdModel struct {
	NavId  int64 `json:"nav_id"`
	SortId int   `json:"sort_id"`
}

type Nav_ChangeStatusModel struct {
	NavIds []int64 `json:"nav_ids"`
	Status int32   `json:"status"`
}
