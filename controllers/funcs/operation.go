package funcs

/**
 * @description: 三元运算
 * @param {bool} condition 条件
 * @param {*} trueVal
 * @param {any} falseVal
 * @return {*}
 */
func IIF(condition bool, trueVal, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}
