package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"tiktok/config"
	"time"
)

// FileDataHook 是一个 logrus 钩子，用于将日志条目写入文件。
type FileDataHook struct {
	file          *os.File     // 当前日志文件指针。
	logPath       string       // 日志文件目录路径。
	fileDate      string       // 当前日志文件的日期（YYYY-MM-DD）。
	FileName      string       // 用于日志文件名的应用程序名称。
	asyncWriter   *AsyncWriter // 用于处理异步日志记录的异步写入器。
	MaxLogAgeDays int          // 添加 MaxLogAgeDays 字段，用于表示最大日志保留天数。
}

// Config 表示日志初始化的配置。
type Config struct {
	LogPath       string `json:"logPath"`       // 日志文件目录路径。
	FileName      string `json:"FileName"`      // 用于日志文件名的应用程序名称。
	LogLevel      string `json:"logLevel"`      // 日志级别，例如 "info"、"warn"、"error" 等。
	MaxLogAgeDays int    `json:"maxLogAgeDays"` // 最大日志保留天数。
	AsyncLogging  bool   `json:"asyncLogging"`  // 是否启用异步日志记录。
}

// Levels 返回 FileDataHook 支持的所有日志级别。
func (hook *FileDataHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 将日志条目写入对应的日志文件，基于日期进行文件切换。
func (hook *FileDataHook) Fire(entry *logrus.Entry) error {
	timer := entry.Time.Format("2006-01-02")
	line, err := entry.String()
	if err != nil {
		return err
	}

	// 如果日期没有变化，将日志条目写入当前日志文件。
	if hook.fileDate == timer {
		if hook.asyncWriter != nil {
			hook.asyncWriter.Write(line)
		} else {
			_, _ = hook.file.WriteString(line)
		}
		return nil
	}

	hook.file.Close()

	// 执行日志文件轮转，切换到新的日志文件。
	hook.rotateLog(timer)

	// 创建新的日志文件以存储新的日志。
	fileName := filepath.Join(hook.logPath, timer, fmt.Sprintf("%s.log", hook.FileName))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Error(err)
		return err
	}

	hook.file = file
	hook.fileDate = timer

	// 如果使用异步写入器，设置新的文件以进行异步日志记录。
	if hook.asyncWriter != nil {
		hook.asyncWriter.SetFile(file)
	} else {
		_, _ = hook.file.WriteString(line)
	}

	return nil
}

// rotateLog 处理日志文件轮转逻辑。
func (hook *FileDataHook) rotateLog(newDate string) {
	if hook.asyncWriter != nil {
		hook.asyncWriter.Stop()
	}

	// 将当前目录重命名为新日期以进行日志文件轮转。
	currentPath := filepath.Join(hook.logPath, hook.fileDate)
	newPath := filepath.Join(hook.logPath, newDate)
	_ = os.Rename(currentPath, newPath)

	// 根据 maxLogAgeDays 删除旧的日志目录。
	hook.deleteOldLogs()

	if hook.asyncWriter != nil {
		hook.asyncWriter.Start()
	}
}

// deleteOldLogs 根据 maxLogAgeDays 删除旧的日志目录。
func (hook *FileDataHook) deleteOldLogs() {
	if hook.fileDate == "" {
		return
	}

	_, err := time.Parse("2006-01-02", hook.fileDate)
	if err != nil {
		return
	}

	files, err := filepath.Glob(filepath.Join(hook.logPath, "*/*.log"))
	if err != nil {
		return
	}

	maxLogAge := time.Duration(hook.MaxLogAgeDays) * 24 * time.Hour
	for _, file := range files {
		fileDate := filepath.Base(filepath.Dir(file))
		logFileDate, err := time.Parse("2006-01-02", fileDate)
		if err != nil {
			continue
		}

		if time.Since(logFileDate) > maxLogAge {
			err = os.Remove(file)
		}
	}
}

// Formatter 是一个自定义的日志格式化器，将附加字段包含在日志输出中。
type Formatter struct{}

var HttpStatusCode int

// Format 将日志条目格式化为字节切片以进行日志记录。
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = &bytes.Buffer{}
	}

	// 基于 HTTP 状态码定义自定义日志级别。
	var level logrus.Level
	if HttpStatusCode >= 500 {
		level = logrus.ErrorLevel
	} else if HttpStatusCode >= 400 {
		level = logrus.WarnLevel
	} else {
		level = logrus.InfoLevel
	}

	HttpStatusCode = 0

	// 格式化日志条目，包括时间戳和日志消息。
	_, _ = fmt.Fprintf(buf, "[%s] %s\n", level.String(), entry.Message)

	return buf.Bytes(), nil
}

// AsyncWriter 表示一个异步日志写入器。
type AsyncWriter struct {
	ch   chan string // 用于传入日志行的通道。
	file *os.File    // 当前日志文件指针。
}

// NewAsyncWriter 创建一个新的异步日志写入器实例。
func NewAsyncWriter(file *os.File) *AsyncWriter {
	writer := &AsyncWriter{
		ch:   make(chan string, 1000),
		file: file,
	}
	writer.Start()
	return writer
}

// Start 启动异步日志写入器。
func (w *AsyncWriter) Start() {
	go func() {
		for line := range w.ch {
			_, _ = w.file.WriteString(line)
		}
	}()
}

// Write 将日志行写入通道以进行异步处理。
func (w *AsyncWriter) Write(line string) {
	w.ch <- line
}

// SetFile 更改当前用于异步日志记录的日志文件。
func (w *AsyncWriter) SetFile(file *os.File) {
	w.Stop()
	w.file = file
	w.Start()
}

// Stop 停止异步日志写入器。
func (w *AsyncWriter) Stop() {
	close(w.ch)
}

// InitFile 使用提供的配置文件初始化日志包。
// configFile 是 JSON 格式的配置文件路径，配置文件包含日志初始化所需的参数。
// 配置文件的格式应该符合 Config 结构体的字段定义。
// 返回错误如果读取配置文件或初始化日志时发生问题。
func Init() error {
	logrus.SetFormatter(&Formatter{})

	// TODO:从配置文件加载配置

	cf := config.DefaultLogConfig

	// 创建当日的日志目录
	fileDate := time.Now().Format("2006-01-02")
	err := os.MkdirAll(filepath.Join(cf.LogPath, fileDate), os.ModePerm)
	if err != nil {
		return err
	}

	// 打开当前日志文件
	fileName := filepath.Join(cf.LogPath, fileDate, fmt.Sprintf("%s.log", cf.FileName))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	var fileHook *FileDataHook
	if cf.AsyncLogging {
		fileHook = &FileDataHook{
			logPath:     cf.LogPath,
			fileDate:    fileDate,
			FileName:    cf.FileName,
			asyncWriter: NewAsyncWriter(file),
		}
	} else {
		fileHook = &FileDataHook{
			file:     file,
			logPath:  cf.LogPath,
			fileDate: fileDate,
			FileName: cf.FileName,
		}
	}

	logrus.AddHook(fileHook)
	logLevel, err := logrus.ParseLevel(cf.LogLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)
	return nil
}

// Logger 是一个用于记录日志的接口。
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

// LogrusLogger 是一个使用 logrus 的日志记录器。
type LogrusLogger struct {
	logger *logrus.Logger
}

// NewLogrusLogger 创建一个新的 LogrusLogger 实例。
func NewLogrusLogger() *LogrusLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	return &LogrusLogger{
		logger: logger,
	}
}

// Debug 记录调试级别的日志。
func (l *LogrusLogger) Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Info 记录信息级别的日志。
func (l *LogrusLogger) Info(args ...interface{}) {
	logrus.Info(args...)
}

// Warn 记录警告级别的日志。
func (l *LogrusLogger) Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Error 记录错误级别的日志。
func (l *LogrusLogger) Error(args ...interface{}) {
	logrus.Error(args...)
}

// Fatal 记录严重错误级别的日志，并导致程序退出。
func (l *LogrusLogger) Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}
