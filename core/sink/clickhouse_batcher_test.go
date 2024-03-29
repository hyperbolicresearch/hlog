package sink

import (
	"context"
	"testing"

	"github.com/ClickHouse/clickhouse-go/v2"
	clickhouseservice "github.com/hyperbolicresearch/hlog/storage/clickhouse"
)

func TestSink(t *testing.T) {
	data := []map[string]interface{}{
		{"foo": "lorem", "bar": int8(1), "_channel": "test_sink"},
		{"foo": "ipsum", "bar": int8(2), "_channel": "test_sink"},
		{"foo": "dolor", "bar": int8(3), "_channel": "test_sink"},
	}
	test := struct {
		name   string
		input  []map[string]interface{}
		expect int
	}{
		"Insertion of maps",
		data,
		len(data),
	}

	// clickhouse: create test_sink db
	addr := []string{"localhost:9000"}
	defaultConn, err := clickhouseservice.Conn(addr)
	if err != nil {
		t.Errorf("Error while connecting to the default db: %v", err)
	}
	if err := defaultConn.Ping(context.Background()); err != nil {
		t.Errorf("Connected, but connot ping the defaul db: %v", err)
	}
	createTestDbQuery := "CREATE DATABASE IF NOT EXISTS test_sink"
	err = defaultConn.Exec(context.Background(), createTestDbQuery)
	if err != nil {
		t.Errorf("Error while creating test db: %v", err)
	}
	// clickhouse: connect to test_sink db
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "test_sink",
			Username: "default",
		},
	})
	if err != nil {
		t.Error(err)
	}
	if err := conn.Ping(context.Background()); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			t.Errorf("Exception [%d] %s \n%s\n",
				exception.Code,
				exception.Message,
				exception.StackTrace)
		}
	}
	err = conn.Exec(context.Background(), "DROP TABLE IF EXISTS test_sink")
	if err != nil {
		t.Error(err)
	}
	schema := `
		CREATE TABLE IF NOT EXISTS test_sink (
			_channel String,
			bar Int8 PRIMARY KEY,
			foo String,
		)
		ENGINE = MergeTree
		ORDER BY bar
	`
	err = conn.Exec(context.Background(), schema)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		conn.Exec(context.Background(), "DROP TABLE IF EXISTS test_table")
	}()
	// batcher: initialializing
	batcher := BatcherWorker{
		Conn: conn,
	}
	t.Run(test.name, func(t *testing.T) {
		count, err := batcher.Sink(test.input)
		if err != nil {
			t.Errorf("error sinking the data: %v", err)
		}
		if count != test.expect {
			t.Errorf("Expected=%v, Got=%v", test.expect, count)
			t.Fail()
		}
	})
}
