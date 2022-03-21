package log

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	config "github.com/toolboxconfig"
	logger "github.com/toolboxlogger"
)

var (
	globalLogger *logger.Logger
	setupOnce    sync.Once
)

func isDebug() bool {
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	return debug
}

func Global() *logger.Logger {
	if globalLogger == nil {
		setupOnce.Do(func() {
			globalLogger = logger.NewCommonLogger("LoggerTest", os.Getenv("VERSION"), os.Getenv("ENV"), config.GetHostName(), isDebug())
		})
	}
	return globalLogger
}

func Entry() *logrus.Entry {
	return Global().Entry()
}

func WithField(key string, value interface{}) *logrus.Entry {
	return Global().WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Global().WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return Global().WithError(err)
}
