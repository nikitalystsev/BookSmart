package logging

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"runtime"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		_, err = w.Write([]byte(line))
		if err != nil {
			return err
		}
	}

	return nil
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func setLogLevel(l *logrus.Logger) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	logLevel := os.Getenv("LOG_LEVEL")

	switch logLevel {
	case "panic":
		l.SetLevel(logrus.PanicLevel)
	case "fatal":
		l.SetLevel(logrus.FatalLevel)
	case "error":
		l.SetLevel(logrus.ErrorLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "trace":
		l.SetLevel(logrus.TraceLevel)
	default:
		l.Warn("Invalid LOG_LEVEL, defaulting to Info level")
		l.SetLevel(logrus.InfoLevel)
	}
}

func NewLogger() (*logrus.Entry, error) {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)

			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		PrettyPrint: false,
	})

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return nil, err
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}

	l.SetOutput(io.Discard)

	writer := []io.Writer{allFile}

	l.AddHook(&writerHook{
		Writer:    writer,
		LogLevels: logrus.AllLevels,
	})

	setLogLevel(l)

	return logrus.NewEntry(l), nil
}

func GetLoggerForTests() (logger *logrus.Entry) {
	l := logrus.New()
	l.SetOutput(io.Discard)
	logger = logrus.NewEntry(l)

	return
}
