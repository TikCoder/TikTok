package log

import (
	"os"
	"path/filepath"
	"testing"
	"tiktok2023/config"
	"time"
)

// TestLogrusLogger 测试使用 logrus 的日志记录器。
func TestLogrusLogger(t *testing.T) {
	// 创建日志记录器实例
	logger := NewLogrusLogger()

	// 初始化日志包
	err := Init()
	if err != nil {
		t.Fatalf("Failed to initialize log package: %v", err)
	}

	// 测试日志记录功能
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// 在真实情况下，日志文件会在异步日志记录器中被写入，我们这里使用时间等待模拟异步写入的等待
	time.Sleep(500 * time.Millisecond)

	// TODO: 在这里添加检查日志文件是否正确写入的代码

	cf := config.DefaultLogConfig

	// 在真实情况下，我们会在日志文件轮转后检查旧的日志文件是否被删除，这里简化测试，只检查目录是否正确被创建
	_ = time.Now().Format("2006-01-02")
	oldDate := time.Now().AddDate(0, 0, -8).Format("2006-01-02")
	oldLogDir := filepath.Join(cf.LogPath, oldDate)
	if _, err := os.Stat(oldLogDir); !os.IsNotExist(err) {
		t.Errorf("Old log directory should have been deleted: %s", oldLogDir)
	}
}
