package domain

// SiteLanguage 站点语言关联视图
// 用于查询站点及其关联的语言信息
// 对应 SQL:
// SELECT a.site_id, b.language_id, a.language_code, b.name, b.description, b.icon
// FROM cms_site a
// LEFT JOIN cms_language b ON a.language_code = b.code
// WHERE b.language_id IS NOT NULL
// ORDER BY b.is_default DESC
type SiteLanguage struct {
	SiteID       int64  `json:"site_id"`       // 站点ID
	LanguageID   int64  `json:"language_id"`   // 语言ID
	LanguageCode string `json:"language_code"` // 语言代码
	Name         string `json:"name"`          // 语言名称
	Description  string `json:"description"`   // 语言描述
	Icon         string `json:"icon"`          // 语言图标
}

// SiteLanguageList 站点语言列表响应
type SiteLanguageList struct {
	List  []*SiteLanguage `json:"list"`
	Total int64           `json:"total"`
}
