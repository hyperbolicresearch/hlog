package ingest

import (
	"reflect"
	"testing"

	"github.com/hyperbolicresearch/hlog/internal/core"
)

func TestGracefulStop(t *testing.T) {
	ingester := NewIngesterWorker()
	go ingester.Start()
	_ = ingester.Stop()
	ingester.RLock()
	isRunning := ingester.IsRunning
	ingester.RUnlock()
	if isRunning != false {
		t.Errorf("Expected=%v, Got=%v", false, ingester.IsRunning)
	}
}

func TestTransform(t *testing.T) {
	var tests = []struct {
		name  string
		input core.Log
		want  map[string]interface{}
	}{
		{
			"Data: {foo: 1}",
			core.Log{
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

	ingester := NewIngesterWorker()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ingester.Lock()
			ingester.Messages = &Messages{
				Data: []*core.Log{&tt.input},
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
