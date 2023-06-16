package global

import "strings"

func IsSuper(roleType string) bool {
	return roleType == SuperFlag
}

// ReverseLowerString 字符串反转并转换为小写字母
func ReverseLowerString(s string) string {
	b := []byte(s)
	n := len(b)
	for i := 0; i < n/2; i++ {
		b[i], b[n-i-1] = b[n-i-1], b[i]
	}
	return strings.ToLower(string(b))
}
