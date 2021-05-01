package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

type ILogger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
}

type logger struct {
	logger   *log.Logger
	writer   io.Writer
	logLevel int
}

var Log ILogger = New(os.Stdout, info)

func New(writer io.Writer, logLevel int) ILogger {
	return &logger{
		logger:   log.New(writer, "", log.LstdFlags),
		writer:   writer,
		logLevel: logLevel,
	}
}

const (
	critical = iota
	err
	warn
	info
	debug
)

func (l *logger) isEnabledLevel(level int) bool {
	return level <= l.logLevel
}

func getTrace() string {
	if pt, file, line, ok := runtime.Caller(4); ok {
		funcName := runtime.FuncForPC(pt).Name()
		return fmt.Sprintf("%20s:%d | %20s | ", file, line, funcName)
	}
	return ""
}

func (l *logger) print(loglevel int, v ...interface{}) {
	if l.isEnabledLevel(loglevel) {
		now := time.Now().Format("2006/01/02 15:04:05")
		// trace := getTrace()
		format := "%v | message:【%v】\n"
		v = append([]interface{}{now}, v...)
		fmt.Fprintf(l.writer, format, v...)
	}
}

func (l *logger) Debug(v ...interface{})    { l.print(debug, v...) }
func (l *logger) Info(v ...interface{})     { l.print(info, v...) }
func (l *logger) Warning(v ...interface{})  { l.print(warn, v...) }
func (l *logger) Error(v ...interface{})    { l.print(err, v...) }
func (l *logger) Critical(v ...interface{}) { l.print(critical, v...) }
