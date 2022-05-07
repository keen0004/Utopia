package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
)

const (
	LOGFLAGS = log.Ldate | log.Ltime
	PDEBUG   = "[DEBUG] "
	PINFO    = "[INFO] "
	PWARN    = "[WARN] "
	PERROR   = "[ERROR] "
)

var (
	logfile     io.Writer
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	defaultPath = "./utopia.log"
)

func init() {
	logfile = os.Stderr
	debugLogger = log.New(logfile, PDEBUG, LOGFLAGS)
	infoLogger = log.New(logfile, PINFO, LOGFLAGS)
	warnLogger = log.New(logfile, PWARN, LOGFLAGS)
	errorLogger = log.New(logfile, PERROR, LOGFLAGS)
}

func Debug(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	debugLogger.Printf(fmt.Sprintf("%s:%d %s", path.Base(file), line, format), v...)
}

func Info(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	infoLogger.Printf(fmt.Sprintf("%s:%d %s", path.Base(file), line, format), v...)
}

func Warn(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	warnLogger.Printf(fmt.Sprintf("%s:%d %s", path.Base(file), line, format), v...)
}

func Error(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	errorLogger.Printf(fmt.Sprintf("%s:%d %s", path.Base(file), line, format), v...)
}

func SetLogPath(path string) {
	if path == "" {
		path = defaultPath
	}

	var err error
	logfile, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	debugLogger.SetOutput(logfile)
	infoLogger.SetOutput(logfile)
	warnLogger.SetOutput(logfile)
	errorLogger.SetOutput(logfile)
}
