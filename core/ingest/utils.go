package ingest

import (
	"context"
	"fmt"
	"sort"
	"strings"

	clickhouseservice "github.com/hyperbolicresearch/hlog/storage/clickhouse"
)

// GetStorableData gets a list of transformed messages and return only a list
// of further transformed messages, with only the fields that will be stored
// on ClickHouse.
func GetStorableData(raw []map[string]interface{}) []map[string]interface{} {
	storableData := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
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

// SortMap takes a map and returns a sorted version of it.
func SortMap(m map[string]interface{}) (map[string]interface{}, []string, []interface{}, error) {
	sortedMap := make(map[string]interface{})
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortedValues := []interface{}{}
	for _, item := range keys {
		sortedMap[item] = m[item]
		sortedValues = append(sortedValues, m[item])

	}
	return sortedMap, keys, sortedValues, nil
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
func GenerateSQLAndApply(schema map[string]interface{}, table string, isAlter bool) error {
	var _sql string
	switch isAlter {
	case true:
		_sql += fmt.Sprintf("ALTER TABLE %s ADD COLUMN (\n", table)
	case false:
		_sql += fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", table)
	}

	_, keys, _, err := SortMap(schema)
	if err != nil {
		return err
	}
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		value := schema[key]
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

	// TODO make configurable
	addrs := []string{"127.0.0.1:9000"}
	chConn, err := clickhouseservice.Conn(addrs)
	if err != nil {
		return err
	}
	err = chConn.Exec(context.Background(), _sql)
	if err != nil {
		return err
	}
	return nil
}
