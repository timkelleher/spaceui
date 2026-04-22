package logger

import (
	"fmt"
	"os"
	"time"
)

var logger LogRepository
var fileName string

func init() {
	logger = NewLogRepository(10)

	fileName = fmt.Sprintf("logs/%s", time.Now().Format("2006-01-02 15:04:05.txt"))
	f, err := os.Create(fileName)
	if err != nil {
		Error(err.Error())
	}
	f.Close()
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

	writeToFile(msg)
}

func Logs() string {
	return logger.Logs()
}

func writeToFile(msg string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		Error(err.Error())
		return
	}
	f.WriteString(msg + "\n")
	f.Close()
}
