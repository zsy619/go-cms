package funcs

import "time"

/**
 * @description: 时间戳转时间
 * @param {time.Time} input 时间戳
 * @return {*}
 */
func Time2Str(input time.Time) string {
	return input.Format("2006-01-02 15:04:05")
}
