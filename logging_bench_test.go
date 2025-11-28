package logging

import (
	"errors"
	"io"
	"testing"
)

func BenchmarkLogger_Info_Simple(b *testing.B) {
	config := Config{
		Level: "debug",
	}
	logger := NewLogger(&config)
	logger = logger.SetWriter(io.Discard)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("test message")
	}
}

func BenchmarkLogger_Info_WithFields(b *testing.B) {
	config := Config{
		Level: "debug",
	}
	logger := NewLogger(&config)
	logger = logger.SetWriter(io.Discard)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "key1", "value1", "key2", "value2")
	}
}

func BenchmarkLogger_Info_WithInts(b *testing.B) {
	config := Config{
		Level: "debug",
	}
	logger := NewLogger(&config)
	logger = logger.SetWriter(io.Discard)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "count", 123, "status", 200)
	}
}

func BenchmarkLogger_Info_WithError(b *testing.B) {
	config := Config{
		Level: "debug",
	}
	logger := NewLogger(&config)
	logger = logger.SetWriter(io.Discard)
	err := errors.New("test error")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", err)
	}
}

func BenchmarkLogger_Info_Slice(b *testing.B) {
	config := Config{
		Level: "debug",
	}
	logger := NewLogger(&config)
	logger = logger.SetWriter(io.Discard)
	slice := []string{"val1", "val2", "val3", "val4"}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", slice)
	}
}

func BenchmarkLogger_Info_SliceValue(b *testing.B) {
	config := Config{
		Level: "debug",
	}
	logger := NewLogger(&config)
	logger = logger.SetWriter(io.Discard)
	slice := []string{"val1", "val2", "val3", "val4"}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "tags", slice)
	}
}
