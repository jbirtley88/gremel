package data

import "reflect"

func GetType(v any) reflect.Kind {
	rv := reflect.ValueOf(v)
	return rv.Kind()
}
