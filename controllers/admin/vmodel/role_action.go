package vmodel

type RoleAction struct {
	IsSuccess   bool `json:"isSuccess" form:"isSuccess"`
	IsSuper     bool `json:"isSuper" form:"isSuper"`
	IsHasAdd    bool `json:"isHasAdd" form:"isHasAdd"`
	IsHasAudit  bool `json:"isHasAudit" form:"isHasAudit"`
	IsHasEdit   bool `json:"isHasEdit" form:"isHasEdit"`
	IsHasDelete bool `json:"isHasDelete" form:"isHasDelete"`
	IsHasView   bool `json:"isHasView" form:"isHasView"`
	IsHasAttach bool `json:"isHasAttach" form:"isHasAttach"`
	IsHasAlbum  bool `json:"isHasAlbum" form:"isHasAlbum"`
}
