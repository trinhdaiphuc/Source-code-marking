package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type AppLog struct {
	LogObject *logrus.Logger
}

func caller() (timeString, packageName, funcName, filename, line string) {
	pc, file, l, _ := runtime.Caller(3)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)

	funcName = parts[pl-1]
	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}
	filename = filepath.Base(file)
	line = fmt.Sprint(l)
	timeString = time.Now().Format(time.RFC3339)
	return
}

func NewAppLog(env string, level string, logOutput string) *AppLog {
	appLog := AppLog{
		LogObject: logrus.New(),
	}
	// set log format
	if env == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	// set log level
	switch level {
	case "TRACE":
		appLog.LogObject.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		appLog.LogObject.SetLevel(logrus.DebugLevel)
	case "INFO":
		appLog.LogObject.SetLevel(logrus.InfoLevel)
	case "WARNING":
		appLog.LogObject.SetLevel(logrus.WarnLevel)
	case "ERROR":
		appLog.LogObject.SetLevel(logrus.ErrorLevel)
	case "CRITICAL":
		appLog.LogObject.SetLevel(logrus.PanicLevel)
	case "FATAL":
		appLog.LogObject.SetLevel(logrus.FatalLevel)
	default:
		appLog.LogObject.SetLevel(logrus.WarnLevel)
	}

	// set log output
	if len(logOutput) > 0 {
		accessLogFileHandler, err := os.OpenFile(logOutput, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		appLog.LogObject.Out = accessLogFileHandler
	}

	return &appLog
}

func (appLog *AppLog) Info(args ...interface{}) {
	_, packageName, funcName, filename, line := caller()
	entry := appLog.LogObject.WithFields(logrus.Fields{
		"PACKAGE": packageName,
		"FILE":    filename,
		"LINE":    line,
		"FUNC":    funcName,
	})
	entry.Info(args...)
}

func (appLog *AppLog) Debug(args ...interface{}) {
	_, packageName, funcName, filename, line := caller()
	entry := appLog.LogObject.WithFields(logrus.Fields{
		"PACKAGE": packageName,
		"FILE":    filename,
		"LINE":    line,
		"FUNC":    funcName,
	})
	entry.Debug(args...)
}

func (appLog *AppLog) Warn(args ...interface{}) {
	_, packageName, funcName, filename, line := caller()
	entry := appLog.LogObject.WithFields(logrus.Fields{
		"PACKAGE": packageName,
		"FILE":    filename,
		"LINE":    line,
		"FUNC":    funcName,
	})
	entry.Warn(args...)
}

func (appLog *AppLog) Error(args ...interface{}) {
	_, packageName, funcName, filename, line := caller()
	entry := appLog.LogObject.WithFields(logrus.Fields{
		"PACKAGE": packageName,
		"FILE":    filename,
		"LINE":    line,
		"FUNC":    funcName,
	})
	entry.Error(args...)
}

func (appLog *AppLog) Fatal(args ...interface{}) {
	_, packageName, funcName, filename, line := caller()
	entry := appLog.LogObject.WithFields(logrus.Fields{
		"PACKAGE": packageName,
		"FILE":    filename,
		"LINE":    line,
		"FUNC":    funcName,
	})
	entry.Fatal(args...)
}

func (appLog *AppLog) Panic(args ...interface{}) {
	_, packageName, funcName, filename, line := caller()
	entry := appLog.LogObject.WithFields(logrus.Fields{
		"PACKAGE": packageName,
		"FILE":    filename,
		"LINE":    line,
		"FUNC":    funcName,
	})
	entry.Panic(args...)
}
