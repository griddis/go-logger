package logging

import (
	"context"
	"fmt"
	"strings"
	"time"

	rz "github.com/griddis/go-logger/libs/rz-go"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// HeadersToLog list headers to log in request info.
var HeadersToLog = []string{
	"content-type",
}

// logger *rz.Logger
func (s *logger) LoggerMiddleware() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		if s.next.GetLevel() == rz.DebugLevel {
			if p, ok := req.(proto.Message); ok {
				mo := protojson.MarshalOptions{
					EmitUnpopulated: true,
					UseProtoNames:   true,
				}
				if b, err := mo.Marshal(p); err == nil {
					s.Debug("incoming gRPC request payload", "method", info.FullMethod, "payload", string(b))
				} else {
					s.Debug("incoming gRPC request payload", "method", info.FullMethod, "payload_type", fmt.Sprintf("%T", req))
				}
			} else {
				s.Debug("incoming gRPC request payload", "method", info.FullMethod, "payload_type", fmt.Sprintf("%T", req))
			}
		}

		// Вызываем основной обработчик
		resp, err := handler(ctx, req)
		if err != nil {
			// Логируем результат
			s.Error("Processed gRPC request",
				"method", info.FullMethod,
				"duration", time.Since(start),
				"code", status.Code(err),
				"error", err,
				s.headersFields(ctx))
		} else {
			s.Info("Processed gRPC request",
				"method", info.FullMethod,
				"duration", time.Since(start),
				"code", status.Code(err),
				s.headersFields(ctx))
			//rz.Error("error", err))
		}
		return resp, err
	}
}

func (s *logger) headersFields(ctx context.Context) interface{} {
	var fields []string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if s.next.GetLevel() == rz.DebugLevel {
			for name, v := range md {
				fields = append(fields, name, strings.Join(v, ","))
			}
		} else {
			for _, h := range HeadersToLog {
				if v, ok := md[h]; ok && len(v) > 0 {
					fields = append(fields, h, strings.Join(v, ","))
				}
			}
		}
	}
	fields = append(fields, tracingFields(ctx)...)
	return fields
}

func tracingFields(ctx context.Context) []string {
	// set tracing info

	if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
		//	bag := baggage.FromContext(ctx)
		return tracingFieldsOtel(span) //, bag)
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		return openTracingFields(span)
	}
	return nil
}

func tracingFieldsOtel(span trace.Span) []string { //bag baggage.Baggage
	if span != nil && span.SpanContext().HasSpanID() {
		return []string{
			"trace_id", span.SpanContext().TraceID().String(),
			"span_id", span.SpanContext().SpanID().String(),
		}
	}
	return nil
}

func openTracingFields(span opentracing.Span) []string {
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		return []string{
			//eagleEyeLogFieldsOtel(bag),
			"trace_id", sc.TraceID().String(),
			"span_id", sc.SpanID().String(),
		}
	}
	return nil
}
