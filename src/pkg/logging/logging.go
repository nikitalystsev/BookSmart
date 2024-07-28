package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
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

	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

type Logger struct {
	*logrus.Entry
}

func NewLogger() (*Logger, error) {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		PrettyPrint: false,
	})

	err := os.MkdirAll("logs", 0755) // исправлено на 0755 для каталогов
	if err != nil {
		return nil, err
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	return &Logger{logrus.NewEntry(l)}, nil
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}
