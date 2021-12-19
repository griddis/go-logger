package logging

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	kitloglevel "github.com/go-kit/kit/log/level"
	rz "gitlab.com/bloom42/libs/rz-go"
)

type loggerKey struct{}

const missingKey string = "undefined"

var (
	ErrLoggerLevel = errors.New("can`t find level in (emerg, alert, crit, err, warn, notice, info, debug)")
	Log            *logger
	// onceInit guarantee initialize logger only once
	onceInit sync.Once
)

type Config struct {
	Level string
	Type  string
	Time  struct {
		Enabled bool
		Format  string
	}
	Caller struct {
		Enabled bool
	}
	DefaultFieldName string
	Format           string
	Writer           *os.File
}

type logger struct {
	next rz.Logger
}

func NewLogger(cfg *Config) Logger {
	/*lvl, err := getLevel(cfg.Level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger init: %s", err)
		os.Exit(1)
	}
	format := "plain"
	var klog kitlog.Logger

	if format == "json" {
		klog = kitlog.NewJSONLogger(kitlog.NewSyncWriter(cfg.Writer))
	} else {
		klog = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(cfg.Writer))
	}
	klog = kitloglevel.NewFilter(klog, lvl)
	if cfg.Time.Enabled {
		klog = kitlog.With(klog, "ts", kitlog.DefaultTimestampUTC)
	}

	onceInit.Do(func() {
		Log = &Logger{klog}
	})*/

	// new logger

	//hostname, _ := os.Hostname()
	logLevel, _ := rz.ParseLevel(cfg.Level)
	/*if err != nil {

	}*/

	// update global logger's context fields
	log := rz.New()
	log = log.With(
		/*rz.Fields(
			rz.String("hostname", hostname),
		),*/
		rz.Level(logLevel),
		//rz.Formatter(rz.FormatterCLI()),
	)
	//log.Info("info from logger", rz.String("hello", "world"))

	onceInit.Do(func() {
		Log = &logger{log}
	})

	return Log
}

func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey{}).(logger); ok {
		return &logger
	}
	return Log
}

func (s *logger) Logger() Logger {
	return s
}

func (s *logger) SetDefaultFieldName(def string) Logger {
	return &logger{s.next.With(rz.MessageFieldName(def))}
}

func (s *logger) SetWriter(writer io.Writer) Logger {
	return &logger{s.next.With(rz.Writer(writer))}
}

func (s *logger) With(keyvals ...interface{}) Logger {
	return &logger{s.next.With(rz.Fields(s._parse(keyvals...)...))}
}

func (s *logger) Log(info string, keyvals ...interface{}) error {
	s.next.Info(info, s._parse(keyvals...)...)
	return nil
}

func (s *logger) Info(info string, keyvals ...interface{}) {
	s.next.Info(info, s._parse(keyvals...)...)
}

func (s *logger) _parse(keyvals ...interface{}) []rz.Field {

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, missingKey)
	}
	var l []rz.Field = make([]rz.Field, len(keyvals)%2)
	for i := 0; i < len(keyvals); i += 2 {
		switch keyvals[i].(type) {
		case string, int, int64:
			l = append(l, rz.String(s._convert(keyvals[i]), s._convert(keyvals[i+1])))
			break
		case error:
			l = append(l, rz.Error("error", keyvals[i].(error)))
			//i -= 1
			break
		//case map[string][]string:

		case []string:
			l = append(l, s._parseSliceString(keyvals[i].([]string))...)
			break
		}
	}

	return l
}

/*func (s *logger) _parseMapSliceString(keyvals map[string][]string) (l []rz.Field) {
	for i, k := range keyvals {
		l = append(l, rz.String(keyvals[i], keyvals[i+1]))
	}
	return l
}*/

func (s *logger) _parseSliceString(keyvals []string) (l []rz.Field) {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, missingKey)
	}
	for i := 0; i < len(keyvals); i += 2 {
		l = append(l, rz.String(keyvals[i], keyvals[i+1]))
	}
	return l
}

func (s *logger) _convert(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	case []string:
		slice := val.([]string)
		return strings.Join(slice[:], ",")
	case time.Duration:
		return val.(time.Duration).String()
	case int:
		return strconv.Itoa(val.(int))
	case int64:
		i, _ := val.(int64)
		return strconv.FormatInt(i, 10)
	case error:
		return val.(error).Error()
	}
	return missingKey
}

func (s *logger) Fatal(info string, keyvals ...interface{}) {
	s.next.Fatal(info, s._parse(keyvals...)...)
}

func (s *logger) Error(info string, keyvals ...interface{}) {
	s.next.Error(info, s._parse(keyvals...)...)
}

func (s *logger) Warn(info string, keyvals ...interface{}) {
	s.next.Warn(info, s._parse(keyvals...)...)
}

func (s *logger) Debug(info string, keyvals ...interface{}) {
	s.next.Debug(info, s._parse(keyvals...)...)
}

/*func (s *loggerKey) With(keyvals ...interface{}) *Logger {
	return &Logger{kitlog.With(s, keyvals...)}
}

func (s *loggerKey) Fatal(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Error(s).Log(keyvals...)
}*/

/*
func (s *Logger) Error(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Error(s).Log(keyvals...)
}

func (s *Logger) Warn(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Warn(s).Log(keyvals...)
}

func (s *Logger) Info(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Info(s).Log(keyvals...)
}

func (s *Logger) Print(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Info(s).Log(keyvals...)
}

func (s *Logger) Debug(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Debug(s).Log(keyvals...)
}

// default logger
func With(keyvals ...interface{}) *Logger {
	return &Logger{kitlog.With(Log, keyvals...)}
}

func Fatal(keyvals ...interface{}) {
	kitloglevel.Error(Log).Log(keyvals...)
}

func Error(keyvals ...interface{}) {
	kitloglevel.Error(Log).Log(keyvals...)
}

func Warn(keyvals ...interface{}) {
	kitloglevel.Warn(Log).Log(keyvals...)
}

func Info(keyvals ...interface{}) {
	kitloglevel.Info(Log).Log(keyvals...)
}

func Print(keyvals ...interface{}) {
	kitloglevel.Info(Log).Log(keyvals...)
}

func Debug(keyvals ...interface{}) {
	kitloglevel.Debug(Log).Log(keyvals...)
}*/

func getLevel(lvl string) (kitloglevel.Option, error) {
	switch lvl {
	case "emerg":
		return kitloglevel.AllowError(), nil
	case "alert":
		return kitloglevel.AllowError(), nil
	case "crit":
		return kitloglevel.AllowError(), nil
	case "err":
		return kitloglevel.AllowError(), nil
	case "warning":
		return kitloglevel.AllowWarn(), nil
	case "notice":
		return kitloglevel.AllowInfo(), nil
	case "info":
		return kitloglevel.AllowInfo(), nil
	case "debug":
		return kitloglevel.AllowDebug(), nil
	}
	return nil, fmt.Errorf("level %s is incorrect. Level can be (emerg, alert, crit, err, warn, notice, info, debug)", lvl)
}

func caller() string {
	_, file, no, ok := runtime.Caller(2)
	if ok {
		/*short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short*/
		return fmt.Sprintf("%v:%v ", file, no)
	}
	return "???:0 "
}

func CustomFormatter() rz.LogFormatter {
	return func(ev *rz.Event) ([]byte, error) {
		var ret = new(bytes.Buffer)
		return ret.Bytes(), nil
	}
}
