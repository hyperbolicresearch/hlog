package ingest

import (
	"context"
	"fmt"
	"strings"

	clickhouse_connector "github.com/hyperbolicresearch/hlog/internal/clickhouse"
)

// GetStorableData gets a list of transformed messages and return only a list
// of further transformed messages, with only the fields that will be stored
// on ClickHouse.
func GetStorableData(raw []map[string]interface{}) []map[string]interface{} {
	storableData := make([]map[string]interface{}, 0, len(raw))
	for _, item := range storableData {
		kv := map[string]interface{}{}
		for k, v := range item {
			if strings.HasPrefix(k, "_") || strings.Contains(k, ".") {
				kv[k] = v
			}
		}
		storableData = append(storableData, kv)
	}
	return storableData
}

// GetDataByChannel transforms a slice of messages and return a
// map where they are grouped by <_channel>.
func GetDataByChannel(data []map[string]interface{}) map[string][]map[string]interface{} {
	dataByChannel := make(map[string][]map[string]interface{})
	for _, item := range data {
		// we are grouping by channel (channel <=> table)
		key := item["_channel"].(string)
		if _, exists := dataByChannel[key]; !exists {
			dataByChannel[key] = []map[string]interface{}{item}
		} else {
			dataByChannel[key] = append(dataByChannel[key], item)
		}
	}
	return dataByChannel
}

// GenerateSQLAndApply generates the SQL query for either creating or altering the
// Clickhouse schema for a given table and makes the given changes to the database.
func GenerateSQLAndApply(schema map[string]string, table string, isAlter bool) error {
	var _sql string
	switch isAlter {
	case true:
		_sql += fmt.Sprintf("ALTER TABLE %s ADD COLUMN (\n", table)
	case false:
		_sql += fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", table)
	}

	for key, value := range schema {
		newLine := fmt.Sprintf("  `%s` %s,\n", key, value)
		// we will sort by logid, so it should not be nullable. indeed,
		// all log is required by design to have a logid.
		if key == "_logid" {
			newLine = fmt.Sprintf("  `%s` String,\n", key)
		}
		_sql += newLine
	}
	_sql += ")"
	_sql += "\nENGINE = MergeTree"
	_sql += "\nPRIMARY KEY (_logid)"
	_sql += "\nORDER BY _logid"
	// _sql += "\nSET allow_nullable_key = true"

	addrs := []string{"127.0.0.1:9000"}
	chConn, err := clickhouse_connector.Conn(addrs)
	if err != nil {
		return err
	}
	err = chConn.Exec(context.Background(), _sql)
	if err != nil {
		return err
	}

	return nil
}
