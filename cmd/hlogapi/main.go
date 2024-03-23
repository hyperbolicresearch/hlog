package main

import (
	"context"
	"fmt"

	"github.com/hyperbolicresearch/hlog/internal/clickhouseservice"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func main() {
	addrs := []string{"127.0.0.1:9000"}
	chConn, err := clickhouseservice.Conn(addrs)
	if err != nil {
		panic(err)
	}

	count(chConn)
	// describe(chConn)
}

func count(chConn driver.Conn) {
	query := "SELECT COUNT(*) FROM default.default"
	row := chConn.QueryRow(context.Background(), query)
	var count uint64
	if err := row.Scan(&count); err != nil {
		panic(err)
	}
	fmt.Println(count)
}

func describe(chConn driver.Conn) {
	query := "SELECT name FROM system.databases"
	rows, err := chConn.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			panic(err)
		}
		fmt.Println(name)
	}
}
