package helpers

import (
	"html"
	"reflect"
)

func EscapeString(myStruct interface{}) {
	value := reflect.ValueOf(myStruct).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if field.Type().Kind() != reflect.String {
			continue
		}

		str := field.String()
		field.SetString(html.EscapeString(str))
	}
}
