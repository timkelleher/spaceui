package logger

import (
	"fmt"
	"time"
)

var logger LogRepository

func init() {
	logger = NewLogRepository(10)
}

func Debug(msg string) {
	log("green", msg)
}

func Info(msg string) {
	log("blue", msg)
}

func Warn(msg string) {
	log("yellow", msg)
}

func Error(msg string) {
	log("red", msg)
}

func log(color, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logger.Push(fmt.Sprintf("[%s]%s[-] %s", color, timestamp, msg))
}

func Logs() string {
	return logger.Logs()
}
