package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
)

type SomeStruct struct {
	RandomList []string `faker:"slice_len=20"`
}

func TestNewLogger(t *testing.T) {
	config := Config{
		Level:  "debug",
		Format: "2006-01-02T15:04:05.999999999Z07:00",
	}
	logger := NewLogger(&config)
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name string
		args args
		want Logger
	}{
		{
			"new",
			args{
				cfg: &config,
			},
			logger,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimBody(t *testing.T) {
	var body = []byte("{ \n    \"notification\": { \n        \"applicationID\":\"58593a7e8004bb1d152ac476\",\n        \"token\":\"1b3bd17ecc2ec7c07ced453943c2ee4273a44db52547079443403d85a414eb46\",\n        \"systemNotification\":\"iOS\",\n        \"message\":\"Test message для Сергея\",\n        \"type\":\"CheckCode\",\n        \"badge\":5,\n        \"data\": {\n            \"data\":\"34587\",\n            \"type\": \"\"\n        }\n    }\n}")

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, body); err != nil {
		fmt.Println(err)
	}
	t.Logf("body() = %v", buffer)
}

/*func TestNewLogger_1(t *testing.T) {
	config := &Config{
		Level:  "debug",
		Format: "2006-01-02T15:04:05.999999999Z07:00",
	}
	logger := NewLogger(config)
	type args struct {
		level      string
		timeFormat string
	}
	tests := []struct {
		name string
		args args
		want Logger
	}{
		{
			"new",
			args{
				"debug",
				"2006-01-02T15:04:05.999999999Z07:00",
			},
			logger,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(tt.args.level, tt.args.timeFormat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}*/

func Test_logger_Info(t *testing.T) {
	type fields struct {
		next Logger
	}
	type args struct {
		info string
		val1 string
		val2 string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"success",
			fields{
				NewLogger(&Config{}).Logger(),
			},
			args{
				"many single success",
				"123",
				"321",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.next
			s.Info(tt.args.info, tt.args.val1, tt.args.val2)
		})
	}
}

func Test_logger_InfoSlice(t *testing.T) {
	_ = faker.SetRandomMapAndSliceSize(20)
	arg := SomeStruct{}
	_ = faker.FakeData(&arg)
	type fields struct {
		next Logger
	}
	type args struct {
		info    string
		keyvals []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"success",
			fields{
				NewLogger(&Config{}).Logger(),
			},
			args{
				"slice success",
				arg.RandomList,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.next
			s.Info(tt.args.info, tt.args.keyvals)
		})
	}
}

func BenchmarkLoggerInfo(b *testing.B) {

	_ = faker.SetRandomMapAndSliceSize(20)
	arg := SomeStruct{}
	_ = faker.FakeData(&arg)
	config := Config{
		Level:  "debug",
		Format: "2006-01-02T15:04:05.999999999Z07:00",
	}
	logger := NewLogger(&config)
	logger = logger.SetWriter(ioutil.Discard)
	fmt.Printf("dump %v\n", arg.RandomList)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1; j++ {
			logger.Info("test", arg.RandomList)
		}
	}
}

/*func Test_logger_Fatal(t *testing.T) {
	type fields struct {
		next Logger
	}
	type args struct {
		info    string
		keyvals []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"success",
			fields{
				NewLogger(&Config{}).Logger(),
			},
			args{
				"many single fatal",
				[]string{"123", "321"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.next
			s.Fatal(tt.args.info, tt.args.keyvals)
		})
	}
}*/

func Test_convert(t *testing.T) {
	type args struct {
		val interface{}
	}
	timeNow := int64(1617000000)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"int",
			args{
				val: 21,
			},
			"21",
		},
		{
			"int32",
			args{
				val: int32(21),
			},
			"21",
		},
		{
			"int64",
			args{
				val: int64(21),
			},
			"21",
		},
		{
			"float32",
			args{
				val: float32(21.2),
			},
			"21.2",
		},
		{
			"float64",
			args{
				val: float64(21.2),
			},
			"21.2",
		},
		{
			"string",
			args{
				val: "test",
			},
			"test",
		},
		{
			"[]string",
			args{
				val: []string{"test", "test2"},
			},
			"test, test2",
		},
		{
			"time",
			args{
				val: time.Unix(timeNow, 0),
			},
			time.Unix(timeNow, 0).String(),
		},
		{
			"error",
			args{
				val: errors.New("test error"),
			},
			"test error",
		},
		{
			"struct",
			args{
				val: SomeStruct{
					RandomList: []string{"test1", "test2"},
				},
			},
			"{[test1 test2]}",
		},
		{
			"misk",
			args{},
			"undefined",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := _convert(tt.args.val); got != tt.want {
				t.Errorf("_convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithAddsFields(t *testing.T) {
	var buf bytes.Buffer

	cfg := &Config{Level: "debug"}
	log := NewLogger(cfg).With("request_id", "abc-123").SetWriter(&buf)
	log.Info("test message")

	ctx := context.Background()
	ctx = WithContext(ctx, log)
	logger2 := FromContext(ctx)
	logger2.Info("logger2")

	output := buf.String()
	t.Logf("Output: %s", output)
	if !strings.Contains(output, "abc-123") {
		t.Errorf("Field was not logged: %s", output)
	}
}
