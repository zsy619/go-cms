package vmodel

type ReportFormsModel struct {
	Dates []string               `json:"dates"`
	Items []ReportFormsItemModel `json:"items"`
}

type ReportFormsItemModel struct {
	Name string  `json:"name"`
	Data []int64 `json:"data"`
	Type string  `json:"type"`
}
