package vmodel

type Notice_SaveSortIdModel struct {
	NoticeId int64 `json:"notice_id"`
	SortId   int   `json:"sort_id"`
}

type Notice_ChangeStatusModel struct {
	NoticeIds []int64 `json:"notice_ids"`
	Status    int32   `json:"status"`
}
