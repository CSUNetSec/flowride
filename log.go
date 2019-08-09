package flowride

import (
	log "github.com/sirupsen/logrus"
)

func CheckLogFatal(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": true,
			"Error": err,
		}).Fatal("terminating")
	}
}

func CheckLogWarn(err error, msg string) {
	if err != nil {
		log.WithFields(log.Fields{
			"Fatal": false,
			"Error": err,
		}).Warn(msg)
	}
}

func LogInfo(msg string) {
	log.WithFields(log.Fields{
		"Fatal": false,
		"Error": nil,
	}).Info(msg)
}

func LogFatal(msg string) {
	log.WithFields(log.Fields{
		"Fatal": true,
		"Error": nil,
	}).Fatal(msg)
}
