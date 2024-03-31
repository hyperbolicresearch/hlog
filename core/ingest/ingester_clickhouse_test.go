package ingest

import (
	"reflect"
	"testing"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/logs"
)

// func TestGracefulStop(t *testing.T) {
// 	stop := make(chan struct{})
// 	ingester := NewClickHouseIngester(&config.DefaultConfig)
// 	go ingester.Start(stop)
// 	_ = ingester.Stop(stop)
// 	ingester.RLock()
// 	isRunning := ingester.IsRunning
// 	ingester.RUnlock()
// 	if isRunning != false {
// 		t.Errorf("Expected=%v, Got=%v", false, ingester.IsRunning)
// 	}
// }

func TestTransform(t *testing.T) {
	var tests = []struct {
		name  string
		input logs.Log
		want  map[string]interface{}
	}{
		{
			"Data: {foo: 1}",
			logs.Log{
				Channel:   "testnet",
				LogId:     "0000-0000-0000-0000-0000",
				SenderId:  "test-1234",
				Timestamp: 1709118220916,
				Level:     "debug",
				Message:   "lorem ipsum dolor",
				Data:      map[string]interface{}{"foo": 1, "bar": "helloworld"},
			},
			map[string]interface{}{
				// metadata
				"_channel":   "testnet",
				"_logid":     "0000-0000-0000-0000-0000",
				"_senderid":  "test-1234",
				"_timestamp": int64(1709118220916),
				"_level":     "debug",
				"_message":   "lorem ipsum dolor",
				"_data":      map[string]interface{}{"bar": "helloworld", "foo": 1},
				// fields
				"foo": 1,
				"bar": "helloworld",
				// field arrays
				"int.keys":      []string{"foo"},
				"int.values":    []int{1},
				"string.keys":   []string{"bar"},
				"string.values": []string{"helloworld"},
			},
		},
	}

	ingester := NewClickHouseIngester(&config.DefaultConfig)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ingester.Lock()
			ingester.Messages = &Messages{
				Data: []*logs.Log{&tt.input},
			}
			ingester.Unlock()
			_ = ingester.Transform()
			ingester.Messages.RLock()
			ans := ingester.Messages.TransformedData[0]
			ingester.Messages.RUnlock()
			eq := reflect.DeepEqual(tt.want, ans)
			if !eq {
				t.Errorf("Expected=%v, Got=%v", tt.want, ans)
			}
		})
	}
}

func TestExtractSchemas(t *testing.T) {
	ingester := NewClickHouseIngester(&config.DefaultConfig)
	err := ingester.ExtractSchemas()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
