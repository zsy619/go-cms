package banner

import (
	"fmt"
	"os"
	"time"
)

// ============================================================
// StartupLogger 启动日志（Spring Boot 风格）
// ============================================================
type StartupLogger struct {
	AppName   string
	PID       int
	StartTime time.Time
}

// NewLogger 创建启动日志记录器
func NewLogger(appName string) *StartupLogger {
	return &StartupLogger{
		AppName:   appName,
		PID:       os.Getpid(),
		StartTime: time.Now(),
	}
}

// timestamp 生成日志时间戳
func (l *StartupLogger) timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// Info 输出 INFO 级别日志
func (l *StartupLogger) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s%s  %s INFO %s %d %s---%s [%s%10s%s] %-40s : %s\n",
		Dim, l.timestamp(), Reset,
		BrightGreen, Reset,
		l.PID,
		Dim, Reset,
		Green, "main", Reset,
		l.AppName,
		msg,
	)
}

// Warn 输出 WARN 级别日志
func (l *StartupLogger) Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s%s  %s WARN %s %d %s---%s [%s%10s%s] %-40s : %s\n",
		Dim, l.timestamp(), Reset,
		Yellow, Reset,
		l.PID,
		Dim, Reset,
		Yellow, "main", Reset,
		l.AppName,
		msg,
	)
}

// Error 输出 ERROR 级别日志
func (l *StartupLogger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s%s  %sERROR%s %d %s---%s [%s%10s%s] %-40s : %s\n",
		Dim, l.timestamp(), Reset,
		Red, Reset,
		l.PID,
		Dim, Reset,
		Red, "main", Reset,
		l.AppName,
		msg,
	)
}

// Summary 输出启动成功汇总
func (l *StartupLogger) Summary(port int, goVersion string) {
	elapsed := time.Since(l.StartTime)

	fmt.Printf("\n%s%s═══════════════════════════════════════════════════════════%s\n",
		Bold, Green, Reset)
	fmt.Printf("%s  YY-Cms 启动成功！%s\n", Bold+Green, Reset)
	fmt.Printf("%s  框架:     %sBeego%s\n", Dim, White+Bold, Reset)
	fmt.Printf("%s  Go 版本:  %s%s%s\n", Dim, White+Bold, goVersion, Reset)
	fmt.Printf("%s  进程 PID: %s%d%s\n", Dim, White+Bold, l.PID, Reset)
	fmt.Printf("%s  服务端口: %s%d%s\n", Dim, White+Bold, port, Reset)
	fmt.Printf("%s  启动耗时: %s%.3f 秒%s\n", Dim, White+Bold, elapsed.Seconds(), Reset)
	fmt.Printf("%s%s═══════════════════════════════════════════════════════════%s\n\n",
		Bold, Green, Reset)
}
