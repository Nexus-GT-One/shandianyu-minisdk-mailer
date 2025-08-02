package maputil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
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
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回map所有值
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0)
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Equals 判断两个map是否相等
func Equals[K string, V any](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}

	// 比较键的数量是否一致
	keys1 := make([]string, 0, len(m1))
	keys2 := make([]string, 0, len(m2))

	for k := range m1 {
		keys1 = append(keys1, fmt.Sprintf("%v", k))
	}
	for k := range m2 {
		keys2 = append(keys2, fmt.Sprintf("%v", k))
	}

	sort.Strings(keys1)
	sort.Strings(keys2)

	if !reflect.DeepEqual(keys1, keys2) {
		return false
	}

	// 比较键的值是否一致
	for _, k := range keys1 {
		m1IsMap := strings.HasPrefix(reflect.TypeOf(m1[K(k)]).String(), "map[")
		m2IsMap := strings.HasPrefix(reflect.TypeOf(m2[K(k)]).String(), "map[")
		if m2IsMap && m2IsMap {
			mm1 := make(map[K]V)
			b1, _ := json.Marshal(m1[K(k)])
			json.Unmarshal(b1, &mm1)
			mm2 := make(map[K]V)
			b2, _ := json.Marshal(m2[K(k)])
			json.Unmarshal(b2, &mm2)
			if !Equals(mm1, mm2) {
				return false
			}
		}

		if (!m1IsMap || !m2IsMap) && fmt.Sprintf("%v", m1[K(k)]) != fmt.Sprintf("%v", m2[K(k)]) {
			return false
		}
	}

	return true
}
