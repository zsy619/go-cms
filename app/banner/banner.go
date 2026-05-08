package banner

import (
	"fmt"
	"os"
	"strings"
)

// ============================================================
// ANSI 颜色常量
// ============================================================
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	Dim         = "\033[2m"
	Italic      = "\033[3m"
	Green       = "\033[32m"
	BrightGreen = "\033[92m"
	Cyan        = "\033[36m"
	BrightCyan  = "\033[96m"
	Yellow      = "\033[33m"
	Red         = "\033[31m"
	White       = "\033[37m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
)

// ============================================================
// 默认 YY-Cms Banner（Spring Boot 风格）
// ============================================================
const defaultBanner = `
${bold}${brightGreen}
██╗   ██╗██╗   ██╗        ██████╗███╗   ███╗███████╗
╚██╗ ██╔╝╚██╗ ██╔╝       ██╔════╝████╗ ████║██╔════╝
 ╚████╔╝  ╚████╔╝        ██║     ██╔████╔██║███████╗
  ╚██╔╝    ╚██╔╝         ██║     ██║╚██╔╝██║╚════██║
   ██║      ██║          ╚██████╗██║ ╚═╝ ██║███████║
   ╚═╝      ╚═╝           ╚═════╝╚═╝     ╚═╝╚══════╝
${reset}
${bold}${green}  :: YY-Cms ::${reset}${dim}${white} (v${version})${reset}
${dim}${white}  Beego Framework | http://yycms.local${reset}
`

// ============================================================
// Banner 模式
// ============================================================
type Mode string

const (
	ModeOff     Mode = "off"
	ModeConsole Mode = "console"
	ModeLog     Mode = "log"
)

// ============================================================
// BannerConfig 配置
// ============================================================
type Config struct {
	Mode       Mode   // 显示模式
	FilePath   string // 自定义 banner 文件路径
	Version    string // 版本号
	AppName    string // 应用名称
	ServerAddr string // 服务地址
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Mode:       ModeConsole,
		FilePath:   "",
		Version:    "1.0.0",
		AppName:    "YY-Cms Application",
		ServerAddr: "localhost:8080",
	}
}

// ============================================================
// Renderer Banner 渲染器
// ============================================================
type Renderer struct {
	config *Config
}

// New 创建 Banner 渲染器
func New(cfg *Config) *Renderer {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &Renderer{config: cfg}
}

// Render 渲染并输出 Banner
func (r *Renderer) Render() {
	output := r.renderToString()
	if output != "" {
		fmt.Print(output)
	}
}

// renderToString 将 Banner 渲染为字符串
func (r *Renderer) renderToString() string {
	if r.config.Mode == ModeOff {
		return ""
	}

	raw := r.loadBannerContent()

	// 构建变量映射
	vars := map[string]string{
		"${version}":     r.config.Version,
		"${appName}":     r.config.AppName,
		"${serverAddr}":  r.config.ServerAddr,
		"${bold}":        Bold,
		"${dim}":         Dim,
		"${italic}":      Italic,
		"${reset}":       Reset,
		"${green}":       Green,
		"${brightGreen}": BrightGreen,
		"${cyan}":        Cyan,
		"${brightCyan}":  BrightCyan,
		"${yellow}":      Yellow,
		"${red}":         Red,
		"${white}":       White,
		"${blue}":        Blue,
		"${magenta}":     Magenta,
	}

	// LOG 模式下清除颜色
	if r.config.Mode == ModeLog {
		for k := range vars {
			if strings.HasPrefix(k, "${") && k != "${version}" &&
				k != "${appName}" && k != "${serverAddr}" {
				vars[k] = ""
			}
		}
	}

	// 执行替换
	replacer := buildReplacer(vars)
	return replacer.Replace(raw)
}

// loadBannerContent 加载 Banner 内容
func (r *Renderer) loadBannerContent() string {
	// 1. 优先从自定义文件加载
	if r.config.FilePath != "" {
		data, err := os.ReadFile(r.config.FilePath)
		if err == nil && len(data) > 0 {
			return string(data)
		}
		fmt.Printf("%s⚠ 无法加载自定义 Banner 文件: %s，使用默认 Banner%s\n",
			Yellow, r.config.FilePath, Reset)
	}

	// 2. 尝试加载项目根目录下的 banner.txt（兼容 Spring Boot 习惯）
	if data, err := os.ReadFile("banner.txt"); err == nil && len(data) > 0 {
		return string(data)
	}

	// 3. 使用内置默认 Banner
	return defaultBanner
}

// buildReplacer 构建字符串替换器
func buildReplacer(vars map[string]string) *strings.Replacer {
	pairs := make([]string, 0, len(vars)*2)
	for k, v := range vars {
		pairs = append(pairs, k, v)
	}
	return strings.NewReplacer(pairs...)
}
