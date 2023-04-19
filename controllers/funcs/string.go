package funcs

import "strings"

/**
 * @description: 字符串拼接
 * @param {...string} input 输入字符串
 * @return {*}
 */
func ConcatStr(input ...string) string {
	var output string
	for _, v := range input {
		output += v
	}
	return output
}

/**
 * @description: 字符串截取
 * @param {string} s 输入字符串
 * @param {*} start 开始位置
 * @param {int} length 截取长度
 * @return {*}
 */
func SubStr(s string, start, length int) string {
	bt := []rune(s)
	var ss string

	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
		ss = string(bt[start:end])
	} else {
		end = start + length
		ss = "..."
		bt1 := []rune(ss)
		ss = string(append(bt[start:end], bt1[0:3]...))
	}
	return ss
}

/**
 * @description: 字符串比较
 * @param {string} x
 * @param {string} y
 * @return {*}
 */
func StrCheck(x, y string) bool {
	return strings.Compare(x, y) == 0
}
