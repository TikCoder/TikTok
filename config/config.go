package config

/*
log 配置
*/
type LogCfg struct {
	// 日志文件目录路径
	LogPath string
	// 用于日志文件名的应用程序名称
	FileName string
	// 日志级别，例如 "info"、"warn"、"error" 等
	LogLevel string
	// 最大日志保留天数
	MaxLogAgeDays int
	// 是否启用异步日志记录
	AsyncLogging bool
}

var DefaultLogConfig = &LogCfg{
	LogPath:       "./logs",
	FileName:      "TikTok",
	LogLevel:      "debug",
	MaxLogAgeDays: 7,
	AsyncLogging:  false,
}
