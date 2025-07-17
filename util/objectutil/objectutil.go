package objectutil

import (
	"fmt"
	"reflect"
)

func ReturnFirstNotNil[T any](obj ...T) T {
	for _, v := range obj {
		if !reflect.ValueOf(v).IsNil() {
			return v
		}
	}

	var zero T
	return zero
}

func IsEnumEquals(e1, e2 any) bool {
	return reflect.DeepEqual(fmt.Sprintf("%v", e1), fmt.Sprintf("%v", e2))
}
