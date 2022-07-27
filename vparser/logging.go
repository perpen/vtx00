package vparser

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// InitLogging configures logging to the log file path, or stderr if empty string.
func InitLogging(path string, level log.Level) {
	if path == "" {
		log.SetOutput(os.Stderr)
	} else {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			log.Info("Failed to log to file, using default stderr")
		}
		log.SetOutput(file)
	}

	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	//log.Infoln("logging level:", level)
}
