package service

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/db"
)

// AdminLogService 管理日志服务
type AdminLogService struct{}

// NewAdminLogService 创建日志服务实例
func NewAdminLogService() *AdminLogService {
	return &AdminLogService{}
}

// AdminLog 日志记录参数
type AdminLog struct {
	UserID     int64
	UserName   string
	Method     string
	Path       string
	Query      string
	StatusCode string
	IP         string
	TenantID   int64
}

// LogRecord 写入一条日志
func (svc *AdminLogService) LogRecord(log *AdminLog) error {
	if log == nil {
		return nil
	}

	// 自动截断超长字段
	path := log.Path
	if len(path) > 128 {
		path = path[:128]
	}

	query := log.Query
	if len(query) > 128 {
		query = query[:128]
	}

	userName := log.UserName
	if len(userName) > 128 {
		userName = userName[:128]
	}

	ip := log.IP
	if len(ip) > 64 {
		ip = ip[:64]
	}

	statusCode := log.StatusCode
	if len(statusCode) > 64 {
		statusCode = statusCode[:64]
	}

	method := log.Method
	if len(method) > 32 {
		method = method[:32]
	}

	logMdl := &domain.CmsAdminLog{
		UserID:     log.UserID,
		UserName:   userName,
		Method:     method,
		Path:       path,
		Query:      query,
		StatusCode: statusCode,
		IP:         ip,
		CreateTime: time.Now(),
		TenantID:   log.TenantID,
		Deleted:    false,
	}

	return db.CmsDatabase.Create(logMdl).Error
}

// LogRecordAsync 异步写入日志（不阻塞请求）
func (svc *AdminLogService) LogRecordAsync(log *AdminLog) {
	go func(l *AdminLog) {
		defer func() {
			if r := recover(); r != nil {
				// 忽略错误，不影响业务
			}
		}()
		_ = svc.LogRecord(l)
	}(log)
}

// LogQueryParams 日志查询参数
type LogQueryParams struct {
	UserID    int64
	Method    string
	Path      string
	StartTime time.Time
	EndTime   time.Time
	Page      int
	PageSize  int
}

// LogList 查询日志列表
func (svc *AdminLogService) LogList(params *LogQueryParams) ([]*domain.CmsAdminLog, int64, error) {
	if params == nil {
		params = &LogQueryParams{
			Page:     1,
			PageSize: 20,
		}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	query := db.CmsDatabase.Model(&domain.CmsAdminLog{}).Where("deleted = ?", false)

	if params.UserID > 0 {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.Method != "" {
		query = query.Where("method = ?", strings.ToUpper(params.Method))
	}
	if params.Path != "" {
		query = query.Where("path LIKE ?", "%"+params.Path+"%")
	}
	if !params.StartTime.IsZero() {
		query = query.Where("create_time >= ?", params.StartTime)
	}
	if !params.EndTime.IsZero() {
		query = query.Where("create_time <= ?", params.EndTime)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*domain.CmsAdminLog
	err := query.Order("create_time DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&list).Error

	return list, total, err
}

// LogDelete 删除日志（逻辑删除）
func (svc *AdminLogService) LogDelete(id int64) error {
	return db.CmsDatabase.Model(&domain.CmsAdminLog{}).
		Where("log_id = ?", id).
		Update("deleted", true).Error
}

// LogBatchDelete 批量删除日志
func (svc *AdminLogService) LogBatchDelete(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return db.CmsDatabase.Model(&domain.CmsAdminLog{}).
		Where("log_id IN ?", ids).
		Update("deleted", true).Error
}

// LogClear 清空日志（软删除全部）
func (svc *AdminLogService) LogClear(beforeTime time.Time) error {
	return db.CmsDatabase.Model(&domain.CmsAdminLog{}).
		Where("create_time < ?", beforeTime).
		Update("deleted", true).Error
}

// 计算 MD5（用于生成日志唯一标识）
func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
