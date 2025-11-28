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

	kitloglevel "github.com/go-kit/log/level"
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
	if logger, ok := ctx.Value(loggerKey{}).(Logger); ok {
		return logger
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

func (s *logger) _parse(keyvals ...any) []rz.Field {
	n := len(keyvals)
	if n == 0 {
		return nil
	}
	// Pre-allocate slice with correct capacity
	l := make([]rz.Field, 0, (n/2)+1)

	for i := 0; i < n; i += 2 {
		var val any = missingKey
		if i+1 < n {
			val = keyvals[i+1]
		}

		switch k := keyvals[i].(type) {
		case error:
			l = append(l, rz.Error("error", k))
		case []string:
			l = append(l, s._parseSliceString(k)...)
		default:
			key := _convert(k)
			switch v := val.(type) {
			case string:
				l = append(l, rz.String(key, v))
			case int:
				l = append(l, rz.Int(key, v))
			case int8:
				l = append(l, rz.Int8(key, v))
			case int16:
				l = append(l, rz.Int16(key, v))
			case int32:
				l = append(l, rz.Int32(key, v))
			case int64:
				l = append(l, rz.Int64(key, v))
			case uint:
				l = append(l, rz.Uint(key, v))
			case uint8:
				l = append(l, rz.Uint8(key, v))
			case uint16:
				l = append(l, rz.Uint16(key, v))
			case uint32:
				l = append(l, rz.Uint32(key, v))
			case uint64:
				l = append(l, rz.Uint64(key, v))
			case float32:
				l = append(l, rz.Float32(key, v))
			case float64:
				l = append(l, rz.Float64(key, v))
			case bool:
				l = append(l, rz.Bool(key, v))
			case error:
				l = append(l, rz.Error(key, v))
			case []string:
				l = append(l, rz.Strings(key, v))
			case time.Duration:
				l = append(l, rz.Duration(key, v))
			case time.Time:
				l = append(l, rz.Time(key, v))
			default:
				l = append(l, rz.String(key, _convert(v)))
			}
		}
	}

	return l
}

func (s *logger) _parseSliceString(keyvals []string) (l []rz.Field) {
	n := len(keyvals)
	l = make([]rz.Field, 0, (n/2)+1)
	for i := 0; i < n; i += 2 {
		val := missingKey
		if i+1 < n {
			val = keyvals[i+1]
		}
		l = append(l, rz.String(keyvals[i], val))
	}
	return l
}

func _convert(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v[:], ", ")
	case time.Duration:
		return v.String()
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	// case float32:
	// 	return strconv.FormatFloat(float64(v), 'f', 3, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case error:
		return v.Error()
	default:
		if v != nil {
			return fmt.Sprint(v)
		} else {
			return missingKey
		}
	}
	// return missingKey
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
