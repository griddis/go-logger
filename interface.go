package logging

import (
	"io"
	"net/http"

	"google.golang.org/grpc"
)

type Logger interface {
	With(keyvals ...interface{}) Logger
	Log(info string, keyvals ...interface{}) error
	Info(info string, keyvals ...interface{})
	Fatal(info string, keyvals ...interface{})
	Error(info string, keyvals ...interface{})
	Warn(info string, keyvals ...interface{})
	Debug(info string, keyvals ...interface{})

	// Logger options
	Logger() Logger
	SetWriter(writer io.Writer) Logger
	SetDefaultFieldName(def string) Logger

	//plugins
	ChiRequestLogger() func(next http.Handler) http.Handler
	LoggerMiddleware() grpc.UnaryServerInterceptor
	LoggerClientMiddleware() grpc.UnaryClientInterceptor
}
