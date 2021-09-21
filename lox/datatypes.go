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
func NewBaseCallable(parameters []*Token) *BaseCallable {
	return &BaseCallable{parameters: parameters}
}

// BaseCallable to create compositions
type BaseCallable struct {
	parameters  []*Token
	environment *Environment
}

// Call method
func (c *BaseCallable) Call(i *Interpreter, paren *Token, arguments []interface{}) (interface{}, error) {
	if len(c.parameters) != len(arguments) {
		return nil, WrongNumberOfArguments(paren, len(arguments), len(c.parameters))
	}

	c.environment = NewEnvironment(i.environment)
	for index, parameter := range c.parameters {
		c.environment.define(parameter.lexeme, arguments[index])
	}

	return nil, nil
}

func NewFunction(statement *FunctionStmt) *Function {
	return &Function{
		BaseCallable: NewBaseCallable(statement.params),
		statement:    statement,
	}
}

// Function representation
type Function struct {
	*BaseCallable
	statement *FunctionStmt
}

func (f *Function) Call(i *Interpreter, paren *Token, arguments []interface{}) (interface{}, error) {
	_, _ = f.BaseCallable.Call(i, paren, arguments)
	prev := *i.environment
	i.environment = f.environment

	defer func() {
		i.environment = &prev
		f.statement.rt.value = false
	}()

	for _, stmt := range f.statement.body.statements {
		_, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}

		if f.statement.rt.value {
			return nil, nil
		}
	}

	return nil, nil
}
