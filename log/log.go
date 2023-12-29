package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// DefaultLogLevel 默认日志级别.
	DefaultLogLevel = logrus.TraceLevel
)

// Options 定义日志配置项.
type Options struct {
	LogLevel string `yaml:"log_level"`
}

/* available level: TRACE DEBUG INFO WARN ERROR FATAL PANIC */
func (o *Options) getLogLevel() logrus.Level {
	if o.LogLevel == "" {
		return DefaultLogLevel
	}
	level, err := logrus.ParseLevel(o.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("unknown log level[%s]", o.LogLevel))
	}
	return level
}

func Config(opts Options) func() {
	return LoggerConfig(logrus.StandardLogger(), opts)
}

func LoggerConfig(logger *logrus.Logger, opts Options) func() {
	logger.SetReportCaller(true)
	logger.SetLevel(opts.getLogLevel())
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})

	logrus.SetOutput(os.Stdout)
	return func() {}
}
