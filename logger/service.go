package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	Reset   = "\033[0m"
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[95m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

type serviceLogger struct {
	serviceName     string
	logFile         *os.File
	logger          *log.Logger
	logPath         string
	isTest          bool
	started         time.Time
	subProcessName  string
	subProcessStart *time.Time
	actionName      string
	actionStart     *time.Time
	tabs            int
}

type ServiceLoggerConfig struct {
	ServiceName string
	Path        string
	Testing     *testing.T
	Prefix      string
}

type ServiceLogger interface {
	StartSubProcess(subprocess string)
	EndSubProcess()
	StartAction(action string)
	EndAction()
	Tabs() string
	//
	Info(msg string)
	Infof(format string, v ...interface{})
	Warn(msg string)
	Warnf(format string, v ...interface{})
	Error(msg string)
	Errorf(format string, v ...interface{})
	Debug(msg string)
	Debugf(format string, v ...interface{})
	//
	GetLogPath() string
	GetServiceName() string
	Close() error
}

func NewServiceLogger(config ServiceLoggerConfig) (*serviceLogger, error) {
	sanitizedName := SanitizeFilename(config.ServiceName)
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s.ans", timestamp)

	logDir := filepath.Join(config.Path, sanitizedName)

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", logDir, err)
	}

	logPath := filepath.Join(logDir, filename)

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %s: %w", logPath, err)
	}

	var writer io.Writer = logFile

	if config.Testing != nil && testing.Verbose() {
		writer = io.MultiWriter(logFile, os.Stdout)
	}

	logger := log.New(writer, config.Prefix, log.Ldate|log.Ltime)
	logger.SetFlags(0)

	slog := &serviceLogger{
		serviceName: config.ServiceName,
		logFile:     logFile,
		logger:      logger,
		logPath:     logPath,
		isTest:      config.Testing != nil,
		started:     time.Now(),
	}

	slog.logger.Printf("%s[SERVICE] %s%s started", Magenta, strings.ToUpper(config.ServiceName), Reset)

	return slog, nil
}

func (sl *serviceLogger) Close() error {
	if sl.logFile == nil {
		return nil
	}

	endtime := time.Since(sl.started)
	hours := int(endtime.Hours())
	minutes := int(endtime.Minutes()) % 60
	seconds := int(endtime.Seconds()) % 60
	milliseconds := int(endtime.Milliseconds()) % 1000

	sl.logger.Printf("%s[SERVICE] %s%s run time: %02dh:%02dm:%02ds.%03dms", Magenta, strings.ToUpper(sl.serviceName), Reset, hours, minutes, seconds, milliseconds)

	return sl.logFile.Close()
}

func (sl *serviceLogger) StartSubProcess(subprocess string) {
	start := time.Now()
	sl.subProcessName = subprocess
	sl.subProcessStart = &start
	sl.tabs += 1
	sl.logger.Printf("\n%s%s[SUBPROCESS] %s%s started", sl.Tabs(), Cyan, strings.ToUpper(subprocess), Reset)
}

func (sl *serviceLogger) EndSubProcess() {
	if sl.subProcessStart == nil {
		return
	}
	endtime := time.Since(*sl.subProcessStart)
	minutes := int(endtime.Minutes()) % 60
	seconds := int(endtime.Seconds()) % 60
	milliseconds := int(endtime.Milliseconds()) % 1000
	sl.logger.Printf("%s%s[SUBPROCESS] %s%s run time: %02dm:%02ds.%03dms\n", sl.Tabs(), Cyan, strings.ToUpper(sl.subProcessName), Reset, minutes, seconds, milliseconds)
	sl.subProcessStart = nil
	sl.tabs -= 1
}

func (sl *serviceLogger) StartAction(action string) {
	start := time.Now()
	sl.actionName = action
	sl.actionStart = &start
	sl.tabs += 1
	sl.logger.Printf("%s%s[ACTION] %s%s started", sl.Tabs(), Green, strings.ToUpper(action), Reset)
}

func (sl *serviceLogger) EndAction() {
	if sl.actionStart == nil {
		return
	}
	endtime := time.Since(*sl.actionStart)
	minutes := int(endtime.Minutes()) % 60
	seconds := int(endtime.Seconds()) % 60
	milliseconds := int(endtime.Milliseconds()) % 1000
	sl.logger.Printf("%s%s[ACTION] %s%s run time: %02dm:%02ds.%03dms", sl.Tabs(), Green, strings.ToUpper(sl.actionName), Reset, minutes, seconds, milliseconds)
	sl.actionStart = nil
	sl.tabs -= 1
}

func (sl *serviceLogger) Tabs() string {
	return strings.Repeat("\t", sl.tabs)
}

func (sl *serviceLogger) Info(msg string) {
	sl.logger.Printf(sl.Tabs()+"\t"+"%s%s", Reset, msg)
}

func (sl *serviceLogger) Infof(format string, v ...interface{}) {
	sl.logger.Printf(sl.Tabs()+"\t"+Reset+format, v...)
}

func (sl *serviceLogger) Warn(msg string) {
	sl.logger.Printf("%s%s[WARN]%s %s", sl.Tabs()+"\t", Yellow, Reset, msg)
}

func (sl *serviceLogger) Warnf(format string, v ...interface{}) {
	sl.logger.Printf(sl.Tabs()+"\t"+Yellow+"[WARN]"+Reset+" "+format, v...)
}

func (sl *serviceLogger) Error(msg string) {
	sl.logger.Printf("%s%s[WARN]%s %s", sl.Tabs()+"\t", Red, Reset, msg)
}

func (sl *serviceLogger) Errorf(format string, v ...interface{}) {
	sl.logger.Printf(sl.Tabs()+"\t"+Red+"[ERROR]"+Reset+" "+format, v...)
}

func (sl *serviceLogger) Debug(msg string) {
	sl.logger.Printf(sl.Tabs()+"[DEBUG] %s", msg)
}

func (sl *serviceLogger) Debugf(format string, v ...interface{}) {
	sl.logger.Printf(sl.Tabs()+"[DEBUG] "+format, v...)
}

func (sl *serviceLogger) GetLogPath() string {
	return sl.logPath
}

func (sl *serviceLogger) GetServiceName() string {
	return sl.serviceName
}

func SanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, "*", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "\"", "_")
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	name = strings.ReplaceAll(name, "|", "_")

	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}

	return strings.Trim(name, "_")
}
