package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger 自定义日志记录器
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

var AppLogger *Logger

// InitLogger 初始化日志记录器
func InitLogger() {
	// 创建logs目录
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err = os.Mkdir(logsDir, 0755)
		if err != nil {
			log.Fatal("Failed to create logs directory:", err)
		}
	}

	// 创建日志文件
	timestamp := time.Now().Format("20060102")
	logFile := filepath.Join(logsDir, "app_"+timestamp+".log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// 设置日志格式
	AppLogger = &Logger{
		infoLogger:  log.New(file, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(file, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(file, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// 同时输出到控制台
	AppLogger.infoLogger.SetOutput(os.Stdout)
	AppLogger.errorLogger.SetOutput(os.Stdout)
	AppLogger.debugLogger.SetOutput(os.Stdout)

	AppLogger.Info("Logger initialized successfully")
}

// Info 记录信息日志
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error 记录错误日志
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Debug 记录调试日志
func (l *Logger) Debug(format string, v ...interface{}) {
	l.debugLogger.Printf(format, v...)
}

// Info 全局Info日志函数
func Info(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Info(format, v...)
	}
}

// Error 全局Error日志函数
func Error(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Error(format, v...)
	}
}

// Debug 全局Debug日志函数
func Debug(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Debug(format, v...)
	}
}