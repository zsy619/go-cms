// Package banner 提供 Spring Boot 风格的启动 Banner 渲染能力：
// ANSI 彩色输出、自定义 banner.txt 模板、${...} 占位符替换，
// 以及输出目标不支持颜色时自动降级为纯文本的防御性处理。
package banner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

// ============================================================
// ANSI 颜色常量（供本包及外部直接引用）
// ============================================================
const (
	Reset       = "\x1b[0m"
	Bold        = "\x1b[1m"
	Dim         = "\x1b[2m"
	Italic      = "\x1b[3m"
	Green       = "\x1b[32m"
	BrightGreen = "\x1b[92m"
	Cyan        = "\x1b[36m"
	BrightCyan  = "\x1b[96m"
	Yellow      = "\x1b[33m"
	Red         = "\x1b[31m"
	White       = "\x1b[37m"
	Blue        = "\x1b[34m"
	Magenta     = "\x1b[35m"
)

// styleTokens 是 Banner 模板（内置或 banner.txt）中可用的样式占位符，
// 映射到对应的 ANSI 码。当输出目标不支持颜色（ModeLog / 非终端 /
// NO_COLOR / TERM=dumb）时，这些占位符会被替换为空串，避免日志与
// 管道中出现裸 ANSI 序列。新增样式只需在此登记，渲染逻辑无需改动。
var styleTokens = map[string]string{
	"${reset}":       Reset,
	"${bold}":        Bold,
	"${dim}":         Dim,
	"${italic}":      Italic,
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

// dataTokens 是 Banner 模板中可用的数据占位符，取值来自 Config。
var dataTokens = []string{
	"${version}",
	"${appName}",
	"${serverAddr}",
}

// bannerTxtFile 兼容 Spring Boot 约定：自动读取运行目录下的 banner.txt。
const bannerTxtFile = "banner.txt"

// utf8BOM 是 UTF-8 字节序标记，Windows 编辑器保存的文本常带此前缀。
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// ============================================================
// 默认 YY-Cms Banner（Spring Boot 风格）
// 模板占位符：数据 ${version}/${appName}/${serverAddr}；
// 样式 ${reset}/${bold}/${dim}/${italic} 及 ${green} 等颜色占位符。
// 自定义 banner.txt 中的未知占位符会被原样保留，便于排查拼写错误。
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
${bold}${green}  :: ${appName} ::${reset}${dim}${white}  (v${version})${reset}
${dim}${white}  Beego Framework | ${serverAddr}${reset}
`

// ============================================================
// Banner 模式
// ============================================================
// Mode 决定 Banner 的渲染与着色策略。
type Mode string

const (
	// ModeOff 完全关闭 Banner 输出。
	ModeOff Mode = "off"
	// ModeConsole 面向终端输出：仅当输出目标是终端且环境允许时着色，
	// 否则（如服务日志重定向）自动降级为纯文本。
	ModeConsole Mode = "console"
	// ModeLog 面向日志输出：始终不携带 ANSI 颜色，保证日志干净可读。
	ModeLog Mode = "log"
)

// ============================================================
// Config Banner 配置
// ============================================================
// Config 描述 Banner 的渲染配置。
type Config struct {
	// Mode 显示模式：ModeConsole（默认）/ ModeLog / ModeOff；
	// 空值或未识别值按 ModeConsole 处理。
	Mode Mode
	// FilePath 可选：自定义 Banner 文件路径，支持 ${...} 占位符；
	// 读取失败或内容为空时自动回退到 ./banner.txt 与内置 Banner。
	FilePath string
	// Version 版本号，对应模板中的 ${version}。
	Version string
	// AppName 应用名称，对应模板中的 ${appName}，亦用于启动日志。
	AppName string
	// ServerAddr 服务监听地址，对应模板中的 ${serverAddr}。
	ServerAddr string
}

// DefaultConfig 返回默认配置。
func DefaultConfig() *Config {
	return &Config{
		Mode:       ModeConsole,
		Version:    "1.0.0",
		AppName:    "YY-Cms Application",
		ServerAddr: "localhost:8080",
	}
}

// ============================================================
// Renderer Banner 渲染器
// ============================================================
// Renderer 负责 Banner 模板的解析与渲染。
// 模板内容按需加载并缓存（sync.Once），多次渲染不会重复读文件；
// 版本号等数据占位符每次渲染时实时取自 Config，保证取值始终最新。
type Renderer struct {
	config   *Config
	once     sync.Once
	template string // 解析后的模板内容，首次渲染时填充
}

// New 创建 Banner 渲染器；cfg 为 nil 时使用默认配置。
func New(cfg *Config) *Renderer {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &Renderer{config: cfg}
}

// Render 渲染 Banner 并输出到标准输出。
// 返回的 error 仅表示写出失败；模板加载失败会在内部回退到默认 Banner 并告警。
func (r *Renderer) Render() error {
	return r.RenderTo(os.Stdout)
}

// RenderTo 渲染 Banner 并写入 w，便于测试或接入日志输出。
// 着色与否由 Mode 与 w 的类型共同决定：非 *os.File 一律视为不支持颜色。
func (r *Renderer) RenderTo(w io.Writer) error {
	if w == nil {
		return fmt.Errorf("banner: 输出 writer 不能为 nil")
	}
	if r.config.Mode == ModeOff {
		return nil
	}
	out := r.substitute(r.resolve(), r.useColor(w))
	if out == "" {
		return nil
	}
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	_, err := io.WriteString(w, out)
	return err
}

// useColor 判断本次渲染是否携带 ANSI 颜色：
// ModeLog 永不彩色；其余模式在环境允许且输出目标是终端时才彩色。
func (r *Renderer) useColor(w io.Writer) bool {
	if r.config.Mode == ModeLog {
		return false
	}
	if !colorAllowedByEnv() {
		return false
	}
	return isTerminalWriter(w)
}

// isTerminalWriter 判断 writer 是否直连终端（字符设备）。
// 被重定向到文件/管道/服务日志时返回 false，从而自动关闭颜色。
func isTerminalWriter(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// colorAllowedByEnv 依据环境变量判断是否允许彩色输出，
// 遵循 no-color.org 约定：NO_COLOR 存在（无论取值）与 TERM=dumb 均视为禁用。
func colorAllowedByEnv() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return true
}

// resolve 返回 Banner 模板内容，只解析一次并缓存。优先级：
// 1) 显式配置的 FilePath（不可读/为空时告警并继续回退）；
// 2) 运行目录下的 banner.txt（兼容 Spring Boot 习惯）；
// 3) 内置 defaultBanner。
func (r *Renderer) resolve() string {
	r.once.Do(func() {
		r.template = defaultBanner

		if fp := r.config.FilePath; fp != "" {
			data, err := os.ReadFile(fp)
			switch {
			case err != nil:
				r.warnf("无法读取自定义 Banner 文件 %q: %v，已回退到默认 Banner", fp, err)
			case len(bytes.TrimSpace(data)) == 0:
				r.warnf("自定义 Banner 文件 %q 内容为空，已回退到默认 Banner", fp)
			default:
				r.template = normalizeTemplate(data)
				return
			}
		}

		if data, err := os.ReadFile(bannerTxtFile); err == nil && len(bytes.TrimSpace(data)) > 0 {
			r.template = normalizeTemplate(data)
		}
	})
	return r.template
}

// normalizeTemplate 清理文件模板：去除 UTF-8 BOM，CRLF 统一为 LF，
// 避免终端/日志出现 ^M 或整行回退错乱。
func normalizeTemplate(data []byte) string {
	s := string(bytes.TrimPrefix(data, utf8BOM))
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// warnf 将加载告警输出到 stderr，不污染 stdout 的业务输出。
func (r *Renderer) warnf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[banner] 警告: %s\n", fmt.Sprintf(format, args...))
}

// substitute 对模板执行占位符替换：数据占位符替换为当前配置值，
// 样式占位符在 color=false 时替换为空串。基于 strings.Replacer 的
// 单次扫描替换，配置值中的 "${...}" 不会被二次替换；未知占位符原样保留。
func (r *Renderer) substitute(tmpl string, color bool) string {
	styleKeys := make([]string, 0, len(styleTokens))
	for k := range styleTokens {
		styleKeys = append(styleKeys, k)
	}
	sort.Strings(styleKeys)

	pairs := make([]string, 0, (len(styleTokens)+len(dataTokens))*2)
	for _, k := range styleKeys {
		v := styleTokens[k]
		if !color {
			v = ""
		}
		pairs = append(pairs, k, v)
	}
	for _, k := range dataTokens {
		pairs = append(pairs, k, r.dataValue(k))
	}
	return strings.NewReplacer(pairs...).Replace(tmpl)
}

// dataValue 返回数据占位符对应的配置值。
func (r *Renderer) dataValue(key string) string {
	switch key {
	case "${version}":
		return r.config.Version
	case "${appName}":
		return r.config.AppName
	case "${serverAddr}":
		return r.config.ServerAddr
	}
	return ""
}
