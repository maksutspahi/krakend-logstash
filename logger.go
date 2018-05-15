//Package gologging provides a logger implementation based on the github.com/op/go-logging pkg
package logstsash

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/devopsfaith/krakend-gologging"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
)

var (
	ErrNothingToLog = errors.New("nothing to log")
	hostname        = "localhost"
)

func init() {
	name, err := os.Hostname()
	if err != nil {
		hostname = name
	}
}

// NewLogger returns a krakend logger wrapping a gologging logger
func NewLogger(cfg config.ExtraConfig, ws ...io.Writer) (*Logger, error) {
	serviceName := "KRAKEND"
	gologging.LoggingPattern = "%{message}"
	if tmp, ok := cfg[gologging.Namespace]; ok {
		if section, ok := tmp.(map[string]interface{}); ok {
			if tmp, ok = section["prefix"]; ok {
				if v, ok := tmp.(string); ok {
					serviceName = v
				}
				delete(section, "prefix")
			}
		}
	}

	loggr, err := gologging.NewLogger(cfg, ws...)
	if err != nil {
		return nil, err
	}

	return &Logger{loggr, serviceName}, nil
}

// Logger is a wrapper over a github.com/devopsfaith/krakend/logging logger
type Logger struct {
	logger      logging.Logger
	serviceName string
}

func (l Logger) format(logLevel LogLevel, v ...interface{}) ([]byte, error) {
	if len(v) == 0 {
		return []byte{}, ErrNothingToLog
	}
	msg, ok := v[0].(string)
	if !ok {
		return []byte{}, ErrNothingToLog
	}
	record, ok := v[1].(map[string]interface{})
	if !ok {
		record = map[string]interface{}{}
	}
	record["@version"] = 1
	record["@timestamp"] = time.Now().Format("2006-01-02T15:04:05.000000-07:00")
	record["module"] = l.serviceName
	record["host"] = hostname
	record["message"] = msg
	record["level"] = logLevel

	return json.Marshal(record)
}

// Debug implements the logger interface
func (l Logger) Debug(v ...interface{}) {
	data, err := l.format(LEVEL_DEBUG, v...)
	if err != nil {
		return
	}
	l.logger.Debug(string(data))
}

// Info implements the logger interface
func (l Logger) Info(v ...interface{}) {
	data, err := l.format(LEVEL_INFO, v...)
	if err != nil {
		return
	}
	l.logger.Info(string(data))
}

// Warning implements the logger interface
func (l Logger) Warning(v ...interface{}) {
	data, err := l.format(LEVEL_WARNING, v...)
	if err != nil {
		return
	}
	l.logger.Warning(string(data))
}

// Error implements the logger interface
func (l Logger) Error(v ...interface{}) {
	data, err := l.format(LEVEL_ERROR, v...)
	if err != nil {
		return
	}
	l.logger.Error(string(data))
}

// Critical implements the logger interface
func (l Logger) Critical(v ...interface{}) {
	data, err := l.format(LEVEL_CRITICAL, v...)
	if err != nil {
		return
	}
	l.logger.Critical(string(data))
}

// Fatal implements the logger interface
func (l Logger) Fatal(v ...interface{}) {
	data, err := l.format(LEVEL_CRITICAL, v...)
	if err != nil {
		return
	}
	l.logger.Fatal(string(data))
}

type LogLevel string

const (
	LEVEL_DEBUG    = "DEBUG"
	LEVEL_INFO     = "INFO"
	LEVEL_WARNING  = "WARNING"
	LEVEL_ERROR    = "ERROR"
	LEVEL_CRITICAL = "CRITICAL"
)