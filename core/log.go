package core

import (
	"fs.video/log"
	"github.com/sirupsen/logrus"
	"strings"
)

//  Lm = LogModel
var (
	LmChainClient      = log.RegisterModule("vc-cli", logrus.InfoLevel) 
	LmChainType        = log.RegisterModule("vc-ty", logrus.InfoLevel)  // types
	LmChainKeeper      = log.RegisterModule("vc-kp", logrus.InfoLevel)  //keeper
	LmChainMsgServer   = log.RegisterModule("vc-ms", logrus.InfoLevel)  //msg server
	LmChainRest        = log.RegisterModule("vc-re", logrus.InfoLevel)  //msg rest
	LmChainMsgAnalysis = log.RegisterModule("vc-mas", logrus.InfoLevel) //msg msg analysis
	LmChainUtil        = log.RegisterModule("vc-ut", logrus.InfoLevel)  // util
)


func BuildLog(funcName string, modules ...log.LogModule) *logrus.Entry {
	moduleName := ""
	for _, v := range modules {
		if moduleName != "" {
			moduleName += "/"
		}
		moduleName += string(v)
	}
	logEntry := log.Log.WithField("module", strings.ToLower(moduleName))
	if funcName != "" {
		logEntry = logEntry.WithField("method", strings.ToLower(funcName))
	}
	return logEntry
}
