package vmodel

type Register_ChangeReadModel struct {
	RegisterIds []int64 `json:"register_ids"`
	IsRead      int32   `json:"is_read"`
}
