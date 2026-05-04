package vmodel

type Ads_SaveSortIdModel struct {
	AdsId  int64 `json:"ads_id"`
	SortId int   `json:"sort_id"`
}

type Ads_ChangeStatusModel struct {
	AdsIds []int64 `json:"ads_ids"`
	Status int32   `json:"status"`
}
