package getfilter

import "strings"

func GenerateWhere(filter FilterItem) []interface{} {
	keySearch := strings.ToLower(filter.Value)
	var keyList []interface{}

	if filter.Operator == "like" {
		keyList = make([]interface{}, 1)
		keyList[0] = "%" + keySearch + "%"
	} else if filter.Operator == "in" {
		keys := strings.Split(keySearch, ",")
		keyList = make([]interface{}, len(keys))
		for i, key := range keys {
			keyList[i] = key
		}
	} else {
		keyList = make([]interface{}, 1)
		keyList[0] = keySearch
	}
	return keyList
}
