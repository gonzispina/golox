package lox

import "reflect"

type dataType string

const (
	number  dataType = "number"
	str     dataType = "string"
	boolean dataType = "boolean"
	object  dataType = "object"
)

func getDataType(v interface{}) dataType {
	if v == nil {
		return object
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		return str
	case reflect.Bool:
		return boolean
	default:
		return number
	}
}
