package logging

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

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
