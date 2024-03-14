package ingest

import (
	"context"
	"testing"

	"github.com/ClickHouse/clickhouse-go/v2"
	clickhouse_connector "github.com/hyperbolicresearch/hlog/internal/clickhouse"
)

func TestSink(t *testing.T) {
	schema := `
		CREATE TABLE IF NOT EXISTS test_sink (
			foo String,
			bar Int8,
		)
	`
	data := []map[string]interface{}{
		{"foo": "lorem", "bar": int8(1)},
		{"foo": "ipsum", "bar": int8(2)},
		{"foo": "dolor", "bar": int8(3)},
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
	defaultConn, err := clickhouse_connector.Conn(addr)
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
		ch := make(chan struct{})
		count, err := batcher.Sink(test.input, ch)
		if err != nil {
			t.Errorf("error sinking the data: %v", err)
		}
		if count != test.expect {
			t.Errorf("Expected=%v, Got=%v", test.expect, count)
			t.Fail()
		}
	})
}
