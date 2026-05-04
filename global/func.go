package global

import (
	"strings"

	"github.com/zsy619/tools/xstring"
)

func IsSuper(roleType string) bool {
	return roleType == SuperFlag
}

// ReverseLowerString 字符串反转并转换为小写字母
func ReverseLowerString(s string) string {
	str, _ := xstring.Reverse(s)
	return strings.ToLower(str)
}
