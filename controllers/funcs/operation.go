package funcs

/**
 * @description: 三元运算
 * @param {bool} condition 条件
 * @param {*} trueVal
 * @param {interface{}} falseVal
 * @return {*}
 */
func IIF(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
