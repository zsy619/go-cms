package service

type Status int32

const (
	StatusNormal Status = iota // 草稿
	StatusSubmit               // 提交
	StatusPass                 // 审核通过
	StatusReject               // 审核未通过
	StatusBack                 // 驳回
)
