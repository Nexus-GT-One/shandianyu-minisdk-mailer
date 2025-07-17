package maputil

import (
	"fmt"
	"reflect"
)

// GetValueFromMap 从映射中获取值，如果键不存在则返回默认值
func GetValueFromMap[K string, V any](m map[K]V, key K, defaultValue V) V {
	if m == nil {
		return defaultValue
	}
	if val, ok := m[key]; ok {
		if reflect.DeepEqual("<nil>", fmt.Sprintf("%v", val)) || reflect.ValueOf(val).IsZero() {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

// ContainsKey 判断映射中是否包含指定键
func ContainsKey[K string, V any](m map[K]V, key K) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

// Keys 返回map所有键
func Keys[K string, V any](m map[K]V) []K {
	keys := make([]K, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
