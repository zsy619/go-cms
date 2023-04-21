package funcs

/**
 * @description:
 * @return {*}
 */
func IIF(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
