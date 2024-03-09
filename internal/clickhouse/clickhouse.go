package clickhouse_connector

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func Conn(addrs []string) (driver.Conn, error) {
	ctx := context.Background()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: addrs,
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
		},
		MaxOpenConns: 5,
	})

	if err != nil {
		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", 
				exception.Code, 
				exception.Message, 
				exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}
