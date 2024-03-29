package sink

import (
	"sort"
)

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

type Values string

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
