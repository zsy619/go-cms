package funcs

import (
	"math"
	"strconv"
	"time"
)

/**
 * @description: 时间轴转时间字符串
 * @param {int} timeUnix 时间戳
 * @return {*}
 */
func UnixTimeFormat(timeUnix int) string {
	// 转化所需模板
	timeLayout := "2006-01-02 15:04:05"
	return time.Unix(int64(timeUnix), 0).Format(timeLayout)
}

/**
 * @description:格式化文件大小单位
 * @param {string} size 文件大小
 * @param {string} delimiter 分隔符
 * @return {*}
 */
func SizeFormat(size, delimiter string) string {
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return ""
	}
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	for i = 0; sizeInt >= 1024 && i < 5; i++ {
		sizeInt /= 1024
	}
	return strconv.FormatFloat(math.Round(float64(sizeInt)), 'f', -1, 64) + delimiter + units[i]
}
