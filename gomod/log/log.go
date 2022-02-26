package log

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

var Log *logrus.Logger = logrus.New()

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func InitLogger(logPath string, isSaveLog bool, logLevel logrus.Level) {
	Log = logrus.New()
	// 
	if isSaveLog {
		baseLogPath := path.Join(logPath)
		writer, err := rotatelogs.New(
			baseLogPath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(baseLogPath),      // 
			rotatelogs.WithMaxAge(time.Hour*24*7),     // 
			rotatelogs.WithRotationTime(time.Hour*24), // 
		)
		if err != nil {
			panic(err.Error())
		}
		lfHook := lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: writer, // 
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		}, nil)
		Log.AddHook(lfHook)
	}
	Log.SetLevel(logLevel)

}

type nilWriter struct {
}

func (nw *nilWriter) Write(data []byte) (n int, err error) {
	return 0, nil
}