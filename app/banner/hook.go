package banner

import (
	"fmt"
	"runtime"

	beegoWeb "github.com/beego/beego/v2/server/web"
)

// BeegoHook 不再依赖 Beego 内部 Hooks，完全手动控制
type BeegoHook struct {
	bannerRenderer *Renderer
	logger         *StartupLogger
}

func NewBeegoHook(bannerCfg *Config) *BeegoHook {
	if bannerCfg == nil {
		bannerCfg = DefaultConfig()
	}
	return &BeegoHook{
		bannerRenderer: New(bannerCfg),
		logger:         NewLogger(bannerCfg.AppName),
	}
}

// BeforeRun 应在 beego.Run() 之前调用
// 打印 Banner + 初始化日志
func (h *BeegoHook) BeforeRun() {
	// 1. 渲染 Banner
	h.bannerRenderer.Render()

	// 2. Spring Boot 风格启动日志
	h.logger.Info("Starting %s using Go %s (%s/%s)",
		h.logger.AppName,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
	h.logger.Info("PID: %d", h.logger.PID)

	// 读取 Beego 运行模式
	runMode := beegoWeb.BConfig.RunMode
	if runMode == "" {
		runMode = "dev"
	}
	h.logger.Info("Beego RunMode: %s", runMode)
	h.logger.Info("No active profile set, falling back to default profile: default")
	h.logger.Info("Initializing Beego HttpServer...")
}

// AfterRun 应在 beego.Run() 成功启动后调用
// 打印端口监听信息 + 启动汇总
func (h *BeegoHook) AfterRun() {
	port := beegoWeb.BConfig.Listen.HTTPPort
	if port == 0 {
		port = 8080
	}
	h.logger.Info("Beego HttpServer started on http port(s): %d", port)
	h.logger.Info("Mapped URL path [/**] onto Beego Router")
	h.logger.Summary(port, runtime.Version())
	fmt.Printf("[YY-Cms] 服务运行中，按 Ctrl+C 退出...\n\n")
}
