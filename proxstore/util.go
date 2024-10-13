package proxstore

import (
	"github.com/brianvoe/gofakeit/v7"
	"reflect"
)

// isNil returns true if the value is nil
func isNil[T any](t T) bool {
	switch val := any(t).(type) {
	case interface{}:
		return val == nil
	case nil:
		return true
	default:
		// panic(fmt.Errorf("not supported type: %s", val))
		if v := reflect.ValueOf(t); (v.Kind() == reflect.Ptr ||
			v.Kind() == reflect.Interface ||
			v.Kind() == reflect.Slice ||
			v.Kind() == reflect.Map ||
			v.Kind() == reflect.Chan ||
			v.Kind() == reflect.Func) && v.IsNil() {
			return true
		}
		return true
	}
}

func RandomString(length int) string {
	return gofakeit.Password(true, true, true, false, false, length)
}
