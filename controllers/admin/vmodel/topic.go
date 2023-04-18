package vmodel

type Topic_SaveSortIdModel struct {
	TopicId int64 `json:"topic_id"`
	SortId  int   `json:"sort_id"`
}

type Topic_ChangeStatusModel struct {
	TopicIds []int64 `json:"topic_ids"`
	Status   int32   `json:"status"`
}
