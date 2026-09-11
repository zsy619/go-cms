package controllers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"

	"haedu.gov.cn/cms/app/cms/service"
)

// MiddlewareLog 操作日志中间件
// 记录所有后台管理操作的请求日志到 cms_admin_log 表
// 不记录以下情况:
//   - 静态资源
//   - 心跳检查
//   - 验证码
//   - 登录相关（保护密码安全）
var MiddlewareLog = func(ctx *context.Context) {
	// 跳过白名单
	if isLogSkip(ctx) {
		return
	}

	// 跳过 GET 请求（只记录写操作）
	if !shouldLogRequest(ctx) {
		return
	}

	// 记录开始时间
	startTime := time.Now()

	// 在请求处理完成后记录日志
	// Beego 中间件通过 defer 实现
	defer func() {
		// 跳过登录/登出请求中的敏感信息
		path := ctx.Request.RequestURI
		method := ctx.Request.Method

		// 获取用户信息
		var userID int64
		var userName string
		if v, ok := ctx.Input.Session("adminId").(int64); ok {
			userID = v
		}
		if v, ok := ctx.Input.Session("adminName").(string); ok {
			userName = v
		}

		// 获取 IP
		ip := ctx.Request.RemoteAddr
		if idx := strings.LastIndex(ip, ":"); idx > 0 {
			ip = ip[:idx]
		}
		if xff := ctx.Request.Header.Get("X-Forwarded-For"); xff != "" {
			ip = xff
		}
		if xri := ctx.Request.Header.Get("X-Real-IP"); xri != "" {
			ip = xri
		}

		// 获取响应状态码
		statusCode := ctx.ResponseWriter.Status
		if statusCode == 0 {
			statusCode = 200
		}

		// 构建查询参数（过滤敏感字段）
		query := buildLogQuery(ctx)

		// 写入日志
		logMdl := &service.AdminLog{
			UserID:     userID,
			UserName:   userName,
			Method:     method,
			Path:       path,
			Query:      query,
			StatusCode: intToStr(statusCode),
			IP:         ip,
			TenantID:   0,
		}

		// 异步写入，不阻塞响应
		service.NewAdminLogService().LogRecordAsync(logMdl)

		// 同时输出到控制台
		cost := time.Since(startTime)
		logs.Info("[%s] %s %s %d %v", userName, method, path, statusCode, cost)
	}()
}

// isLogSkip 判断是否跳过日志记录
func isLogSkip(ctx *context.Context) bool {
	path := ctx.Request.RequestURI

	// 跳过静态资源
	if strings.HasPrefix(path, "/static/") {
		return true
	}
	if strings.HasPrefix(path, "/uploads/") {
		return true
	}
	if strings.HasPrefix(path, "/downloads/") {
		return true
	}

	// 跳过验证码
	if strings.Contains(path, "/captcha") {
		return true
	}

	// 跳过帮助文档
	if strings.HasPrefix(path, "/help/") {
		return true
	}

	// 跳过健康检查
	if path == "/health" || path == "/ping" {
		return true
	}

	// 跳过 SSE
	if strings.HasPrefix(path, "/sse") {
		return true
	}

	return false
}

// shouldLogRequest 判断是否记录此请求
func shouldLogRequest(ctx *context.Context) bool {
	method := strings.ToUpper(ctx.Request.Method)
	path := ctx.Request.RequestURI

	// 跳过登录接口（保护密码）
	if strings.Contains(path, "/admin/login") || strings.Contains(path, "/admin/logout") {
		return false
	}

	// 只记录后台操作
	if !strings.HasPrefix(path, "/admin/") && !strings.HasPrefix(path, "/cms/admin/") {
		return false
	}

	// 记录所有写操作
	if method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {
		return true
	}

	// 记录重要的 GET 操作
	if method == "GET" {
		if strings.Contains(path, "delete") || strings.Contains(path, "sort") ||
			strings.Contains(path, "status") || strings.Contains(path, "save") ||
			strings.Contains(path, "init") || strings.Contains(path, "reset") {
			return true
		}
	}

	return false
}

// buildLogQuery 构建日志查询参数
func buildLogQuery(ctx *context.Context) string {
	method := strings.ToUpper(ctx.Request.Method)

	// POST 请求从 Body 读取
	if method == "POST" || method == "PUT" || method == "PATCH" {
		if ctx.Request.Body != nil {
			body := make(map[string]string)
			ctx.Request.ParseForm()
			for k, v := range ctx.Request.PostForm {
				val := strings.Join(v, ",")
				// 过滤敏感字段
				if isSensitiveField(k) {
					val = "***"
				}
				body[k] = val
			}
			if len(body) > 0 {
				jsonBytes, err := json.Marshal(body)
				if err == nil {
					return string(jsonBytes)
				}
			}
		}
	}

	// GET 请求从 URL 参数读取
	query := ctx.Request.URL.RawQuery
	return query
}

// isSensitiveField 判断是否为敏感字段
func isSensitiveField(field string) bool {
	field = strings.ToLower(field)
	sensitive := []string{"password", "passwd", "pwd", "secret", "token", "key", "credential"}
	for _, s := range sensitive {
		if strings.Contains(field, s) {
			return true
		}
	}
	return false
}

// intToStr 整数转字符串
func intToStr(i int) string {
	if i == 0 {
		return "0"
	}
	negative := false
	if i < 0 {
		negative = true
		i = -i
	}
	digits := "0123456789"
	result := ""
	for i > 0 {
		result = string(digits[i%10]) + result
		i /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}

// LogAllMiddleware 记录所有请求的日志（用于审计）
var LogAllMiddleware = func(ctx *context.Context) {
	if isLogSkip(ctx) {
		return
	}

	// 仅记录 /admin/* 和 /cms/admin/*
	path := ctx.Request.RequestURI
	if !strings.HasPrefix(path, "/admin/") && !strings.HasPrefix(path, "/cms/admin/") {
		return
	}

	defer func() {
		var userID int64
		var userName string
		if v, ok := ctx.Input.Session("adminId").(int64); ok {
			userID = v
		}
		if v, ok := ctx.Input.Session("adminName").(string); ok {
			userName = v
		}

		ip := ctx.Request.RemoteAddr
		if idx := strings.LastIndex(ip, ":"); idx > 0 {
			ip = ip[:idx]
		}
		if xff := ctx.Request.Header.Get("X-Forwarded-For"); xff != "" {
			ip = xff
		}

		statusCode := ctx.ResponseWriter.Status
		if statusCode == 0 {
			statusCode = 200
		}

		query := buildLogQuery(ctx)

		logMdl := &service.AdminLog{
			UserID:     userID,
			UserName:   userName,
			Method:     ctx.Request.Method,
			Path:       path,
			Query:      query,
			StatusCode: intToStr(statusCode),
			IP:         ip,
		}

		service.NewAdminLogService().LogRecordAsync(logMdl)
	}()
}
