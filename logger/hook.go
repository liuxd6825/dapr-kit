package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"time"
)

// Level 日志级别
type Level = logrus.Level

type logHook struct {
	logger *logrus.Logger
}

// OutputType 输出类型
type OutputType int

type consoleWrite struct {
}

const (
	OutputTypeConsole OutputType = iota
	OutputTypeFile
	OutputTypeAll
)

var fileLogger *logrus.Logger
var outputType OutputType

func newLogHook(logger *logrus.Logger) *logHook {
	return &logHook{
		logger: logger,
	}
}

func (h *logHook) Levels() []logrus.Level {
	return []Level{logrus.DebugLevel, logrus.InfoLevel, logrus.ErrorLevel, logrus.PanicLevel, logrus.WarnLevel, logrus.TraceLevel}
}

func (h *logHook) Fire(entry *logrus.Entry) error {
	if fileLogger == nil {
		return nil
	}
	switch entry.Level {
	case logrus.DebugLevel:
		fileLogger.WithFields(entry.Data).Debug(entry.Message)
	case logrus.InfoLevel:
		fileLogger.WithFields(entry.Data).Info(entry.Message)
	case logrus.ErrorLevel:
		fileLogger.WithFields(entry.Data).Error(entry.Message)
	case logrus.PanicLevel:
		fileLogger.WithFields(entry.Data).Panic(entry.Message)
	case logrus.WarnLevel:
		fileLogger.WithFields(entry.Data).Warn(entry.Message)
	case logrus.TraceLevel:
		fileLogger.WithFields(entry.Data).Trace(entry.Message)
	}
	return nil
}

func InitFileLogger(saveFile string, level string, saveDays int, rotationHour int, outType string) error {
	var err error
	outputType, err = ParseOutputType(outType)
	if err != nil {
		return err
	}

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	if fileLogger == nil && outputType != OutputTypeConsole {
		fileLogger = logrus.New()
		fileLogger.Log(lvl)
		fileLogger.SetOutput(newFileWriter(saveFile, saveDays, rotationHour))
	}
	return nil
}

func newFileWriter(saveFile string, saveDays int, rotationHour int) io.Writer {
	saveFile = strings.ReplaceAll(saveFile, "___", "")
	logFile := saveFile + ".%Y-%m-%d-%H.log"
	// 配置日志每隔 1 小时轮转一个新文件，保留最近 30 天的日志文件，多余的自动清理掉。
	writer, _ := rotatelogs.New(
		logFile,
		rotatelogs.WithLinkName(saveFile),
		rotatelogs.WithMaxAge(time.Duration(24*saveDays)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(rotationHour)*time.Hour),
	)
	return writer
}

func newConsoleWrite() io.Writer {
	return &consoleWrite{}
}

func (c *consoleWrite) Write(p []byte) (n int, err error) {
	if outputType != OutputTypeFile {
		return os.Stderr.Write(p)
	}
	return 0, err
}

func ParseOutputType(val string) (OutputType, error) {
	val = strings.ToLower(val)
	switch val {
	case "console":
		return OutputTypeConsole, nil
	case "file":
		return OutputTypeFile, nil
	case "all":
		return OutputTypeAll, nil
	default:
		return OutputTypeConsole, nil
	}
}
