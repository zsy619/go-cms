package biz

import (
	"haedu.gov.cn/tools/xcache"
)

// ApiCache 缓存
var ApiCache *xcache.ExpiredMap

// init 初始化
func init() {
	ApiCache = xcache.NewExpiredMap()
}

// CleanCahe 清除缓存
func CleanCahe() {
	ApiCache = xcache.NewExpiredMap()
}
