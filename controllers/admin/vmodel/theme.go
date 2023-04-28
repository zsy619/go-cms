package vmodel

type Theme_SaveSortIdModel struct {
	ThemeId int64 `json:"theme_id"`
	SortId  int   `json:"sort_id"`
}

type ThemeConfig struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Mail        string `json:"mail"`
	Version     string `json:"version"`
}
