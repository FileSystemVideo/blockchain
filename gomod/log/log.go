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
	// 如果需要保存日志
	if isSaveLog {
		baseLogPath := path.Join(logPath)
		writer, err := rotatelogs.New(
			baseLogPath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(time.Hour*24*7),     // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
		)
		if err != nil {
			panic(err.Error())
		}
		lfHook := lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
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