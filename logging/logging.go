package logging

import (
	"freeDNS/config"
	"os"

	"github.com/imafaz/logger"
)

func Init() {
	debug := os.Getenv("freeDNS_DEBUG")
	if debug == "true" {
		config.Debug = true
	} else {
		config.Debug = false
	}
	logger.SetLogFile(config.GetLogPath())
	logger.SetOutput(logger.CONSOLE_AND_FILE)
}
func Debug(messages ...string) {
	if config.Debug {
		logger.Debug(messages...)
	}
}
func Debugf(format string, args ...interface{}) {
	if config.Debug {
		logger.Debugf(format, args...)
	}
}
