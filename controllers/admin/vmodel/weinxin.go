package vmodel

type Account_SaveSortIdModel struct {
	AccountId int64 `json:"account_id"`
	SortId    int   `json:"sort_id"`
}

type Account_ChangeStatusModel struct {
	AccountIds []int64 `json:"account_ids"`
	Status     int32   `json:"status"`
}

type Menu_SaveSortIdModel struct {
	MenuId int64 `json:"menu_id"`
	SortId int   `json:"sort_id"`
}

type Menu_ChangeStatusModel struct {
	MenuIds []int64 `json:"menu_ids"`
	Status  int32   `json:"status"`
}

type Rule_SaveSortIdModel struct {
	RuleId int64 `json:"rule_id"`
	SortId int32 `json:"sort_id"`
}
