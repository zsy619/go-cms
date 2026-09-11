package admin

import (
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
)

func timeNow() time.Time {
	return time.Now()
}

type LogController struct{ BaseController }

// LogPage 日志列表页面
// @router /admin/log/page [get]
func (ctrl *LogController) LogPage() {
	ctrl.display("admin/Admin/Log")
}

// LogData 日志列表数据
// @router /admin/log/data [get]
func (ctrl *LogController) LogData() {
	page, _ := ctrl.GetInt("page", 1)
	pageSize, _ := ctrl.GetInt("limit", 20)

	userID, _ := ctrl.GetInt64("user_id")
	params := &service.LogQueryParams{
		UserID:   userID,
		Method:   ctrl.GetString("method"),
		Path:     ctrl.GetString("path"),
		Page:     page,
		PageSize: pageSize,
	}

	list, total, err := service.NewAdminLogService().LogList(params)
	if err != nil {
		logs.Error("LogData error:", err)
		ctrl.JSONPage(lib.CodeError, err.Error(), nil, 0)
		return
	}

	ctrl.JSONPageSuccess(list, total)
}

// LogDelete 删除日志
// @router /admin/log/delete [get]
func (ctrl *LogController) LogDelete() {
	id, _ := ctrl.GetInt64("id")
	if id == 0 {
		ctrl.JSONError("id is required")
		return
	}

	if err := service.NewAdminLogService().LogDelete(id); err != nil {
		logs.Error("LogDelete error:", err)
		ctrl.JSONError(err.Error())
		return
	}

	ctrl.JSONSuccess("删除成功", nil)
}

// LogBatchDelete 批量删除
// @router /admin/log/batch/delete [get]
func (ctrl *LogController) LogBatchDelete() {
	idsStr := ctrl.GetString("ids")
	if idsStr == "" {
		ctrl.JSONError("ids is required")
		return
	}

	var ids []int64
	for _, s := range splitString(idsStr, ",") {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	if err := service.NewAdminLogService().LogBatchDelete(ids); err != nil {
		logs.Error("LogBatchDelete error:", err)
		ctrl.JSONError(err.Error())
		return
	}

	ctrl.JSONSuccess("批量删除成功", nil)
}

// LogClear 清空日志（保留最近 N 天）
// @router /admin/log/clear [get]
func (ctrl *LogController) LogClear() {
	days, _ := ctrl.GetInt("days", 30)
	if days <= 0 {
		days = 30
	}

	cutoff := timeNowAddDays(-days)
	if err := service.NewAdminLogService().LogClear(cutoff); err != nil {
		logs.Error("LogClear error:", err)
		ctrl.JSONError(err.Error())
		return
	}

	ctrl.JSONSuccess("清理成功", nil)
}

// splitString 字符串分割辅助
func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	result := []string{}
	current := ""
	for _, ch := range s {
		if string(ch) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// timeNowAddDays 计算 N 天前的时间
func timeNowAddDays(days int) time.Time {
	return timeNow().AddDate(0, 0, days)
}
