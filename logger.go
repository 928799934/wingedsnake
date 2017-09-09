package wingedsnake

import (
	"log"
)

var logger *log.Logger

// SetLogger 设置logger
func SetLogger(log *log.Logger) {
	logger = log
}

func logf(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf(format, args...)
		return
	}
	log.Printf(format, args...)
}
