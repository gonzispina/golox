package lox

import "reflect"

type dataType string

const (
	number   dataType = "number"
	str      dataType = "string"
	boolean  dataType = "boolean"
	object   dataType = "object"
	callable dataType = "callable"
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
	case reflect.Func:
		return callable
	default:
		return number
	}
}

// Callable contract for all callables in lox
type Callable interface {
	Call(i *Interpreter, paren *Token, arguments []interface{}) (interface{}, error)
}

// NewBaseCallable constructor
func NewBaseCallable(parameters []string) *BaseCallable {
	return &BaseCallable{parameters: parameters}
}

// BaseCallable to create compositions
type BaseCallable struct {
	parameters []string
}

// Call method
func (c *BaseCallable) Call(i *Interpreter, paren *Token, arguments []interface{}) (interface{}, error) {
	if len(c.parameters) != len(arguments) {
		return nil, WrongNumberOfArguments(paren, len(arguments), len(c.parameters))
	}

	return nil, nil
}
