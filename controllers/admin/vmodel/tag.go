package vmodel

type Tag_SaveSortIdModel struct {
	TagId  int64 `json:"tag_id"`
	SortId int   `json:"sort_id"`
}

type Tag_ChangeStatusModel struct {
	TagIds []int64 `json:"tag_ids"`
	Status int32   `json:"status"`
}
