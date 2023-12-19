package model

type TreeNode struct {
	Id       int64       `json:"id"`
	Name     string      `json:"name"`
	Open     bool        `json:"open"`
	Checked  bool        `json:"checked"`
	Selected bool        `json:"selected"`
	Spread   bool        `json:"spread"`
	Icon     string      `json:"icon"`
	Children []*TreeNode `json:"children"`
}
