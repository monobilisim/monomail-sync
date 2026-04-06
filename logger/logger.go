package logger

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = logrus.New()
var log = Log

func SetupLogger() {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		Log.Warnf("Failed to create log directory: %v", err)
	}

	logFile := filepath.Join(logDir, "monomail-sync.log")

	multiWriter := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})

	Log.SetOutput(multiWriter)

	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			slash := strings.LastIndex(f.File, "/")
			filename := f.File[slash+1:]
			return "", "[" + filename + ":" + strconv.Itoa(f.Line) + "]"
		},
	})

	Log.SetReportCaller(true)
	Log.SetLevel(logrus.InfoLevel)
}
