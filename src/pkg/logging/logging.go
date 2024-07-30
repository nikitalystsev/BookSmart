package logging

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

// InfluxDBWriter реализует интерфейс io.Writer для записи логов в InfluxDB
type InfluxDBWriter struct {
	writeAPI api.WriteAPI
}

func (w *InfluxDBWriter) Write(p []byte) (n int, err error) {
	var logEntry map[string]interface{}

	if err = json.Unmarshal(p, &logEntry); err != nil {
		return 0, fmt.Errorf("failed to unmarshal log entry: %v", err)
	}

	point := influxdb2.NewPointWithMeasurement("logs").SetTime(time.Now())

	for key, value := range logEntry {
		switch v := value.(type) {
		case float64:
			point.AddField(key, v)
		case string:
			point.AddField(key, v)
		case bool:
			point.AddField(key, v)
		case map[string]interface{}:
			for subKey, subValue := range v {
				point.AddField(fmt.Sprintf("%s_%s", key, subKey), subValue)
			}
		default:
			point.AddField(key, fmt.Sprintf("%v", value))
		}
	}

	w.writeAPI.WritePoint(point)
	w.writeAPI.Flush()

	return len(p), nil
}

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

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return nil, err
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}

	l.SetOutput(io.Discard)

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// InfluxDB settings
	influxDBUrl := os.Getenv("INFLUXDB_URL")
	influxDBToken := os.Getenv("INFLUXDB_INIT_CLIENT_TOKEN")
	influxDBOrg := os.Getenv("INFLUXDB_INIT_ORG")
	influxDBBucket := os.Getenv("INFLUXDB_INIT_BUCKET")

	client := influxdb2.NewClient(influxDBUrl, influxDBToken)
	writeAPI := client.WriteAPI(influxDBOrg, influxDBBucket)

	influxWriter := &InfluxDBWriter{writeAPI: writeAPI}

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout, influxWriter},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	return &Logger{logrus.NewEntry(l)}, nil
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}
