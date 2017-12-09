package util

import (
	"fmt"
)

// StringifyKeysMapValue recurses into v and changes all instances of
// map[interface{}]interface{} to map[string]interface{}. This is useful to
// work around the impedence mismatch between JSON and YAML unmarshaling that's
// described here:
//
// https://github.com/go-yaml/yaml/issues/139
func StringifyKeysMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return stringifyKeysInterfaceArray(v)
	case map[interface{}]interface{}:
		return stringifyKeysInterfaceMap(v)
	default:
		return v
	}
}

//
// helpers
//

func stringifyKeysInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = StringifyKeysMapValue(v)
	}
	return res
}

func stringifyKeysInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		stringK, ok := k.(string)
		if !ok {
			panic(fmt.Sprintf("Expected all keys in maps to be strings. Got: %v", k))
		}
		res[stringK] = StringifyKeysMapValue(v)
	}
	return res
}
