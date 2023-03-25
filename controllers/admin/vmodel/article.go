package vmodel

type Article_SaveSortIdModel struct {
	ArticleId int64 `json:"article_id"`
	SortId    int   `json:"sort_id"`
}

type Article_ChangeStatusModel struct {
	ArticleIds []int64 `json:"article_ids"`
	Status     int32   `json:"status"`
}
