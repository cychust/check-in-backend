package controller

import (
	"github.com/sirupsen/logrus"
	"phs-mp-develop/src/util/log"
)

var indexLogger = log.GetLogger()

func WriteLog(funcName, filename, errMsg string, err error) {
	writeLog(funcName, filename, errMsg, err)
}

func writeLog(filename, funcName, errMsg string, err error) {
	indexLogger.WithFields(logrus.Fields{
		"package":  "controller",
		"file":     filename,
		"function": funcName,
		"err":      err,
	}).Warn(errMsg)
}

