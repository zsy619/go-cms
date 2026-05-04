package global

import (
	"haedu.gov.cn/tools/xstring"
	"strings"
)

func IsSuper(roleType string) bool {
	return roleType == SuperFlag
}

// ReverseLowerString 字符串反转并转换为小写字母
func ReverseLowerString(s string) string {
	str, _ := xstring.Reverse(s)
	return strings.ToLower(str)
}
