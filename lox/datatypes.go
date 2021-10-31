package lox

import "reflect"

type dataType string

const (
	number   dataType = "number"
	str      dataType = "string"
	boolean  dataType = "boolean"
	class    dataType = "class"
	object   dataType = "object"
	function dataType = "function"
)

func getDataType(v interface{}) dataType {
	if v == nil {
		return object
	}
	if _, ok := v.(*Class); ok {
		return class
	}
	if _, ok := v.(Instance); ok {
		return object
	}
	if _, ok := v.(*Function); ok {
		return function
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

// Callable contract for all callables in lox
type Callable interface {
	Call(i *Interpreter, paren *Token, arguments []interface{}) (interface{}, error)
}

// NewBaseCallable constructor
func NewBaseCallable(parameters []*Token, closure *Environment) *BaseCallable {
	return &BaseCallable{
		parameters: parameters,
		closure:    closure,
	}
}

// BaseCallable to create compositions
type BaseCallable struct {
	parameters  []*Token
	environment *Environment
	closure     *Environment
}

// Call method
func (c *BaseCallable) Call(i *Interpreter, paren *Token, arguments []interface{}) (interface{}, error) {
	if len(c.parameters) != len(arguments) {
		return nil, WrongNumberOfArguments(paren, len(arguments), len(c.parameters))
	}

	c.environment = NewEnvironment(c.closure)
	for index, parameter := range c.parameters {
		c.environment.define(parameter.lexeme, arguments[index])
	}

	return nil, nil
}

func NewFunction(statement *FunctionStmt, closure *Environment) *Function {
	return &Function{
		BaseCallable: NewBaseCallable(statement.params, closure),
		statement:    statement,
	}
}

// Function representation
type Function struct {
	*BaseCallable
	statement *FunctionStmt
}

func (f *Function) String() string {
	return "function"
}

func (f *Function) Call(i *Interpreter, paren *Token, arguments []interface{}) (interface{}, error) {
	_, err := f.BaseCallable.Call(i, paren, arguments)
	if err != nil {
		return nil, err
	}

	prev := *i.environment
	i.environment = f.environment

	defer func() {
		i.environment = &prev
		*f.statement.rt = false
	}()

	for _, stmt := range f.statement.body.statements {
		s, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}

		if *f.statement.rt {
			return s, nil
		}
	}

	return nil, nil
}

// NewClass constructor
func NewClass(statement *ClassStmt, methods map[string]*Function, closure *Environment) *Class {
	return &Class{
		BaseCallable: NewBaseCallable(make([]*Token, 0), closure),
		statement:    statement,
		methods:      methods,
	}
}

// Class representation
type Class struct {
	*BaseCallable
	statement *ClassStmt
	methods   map[string]*Function
}

func (c *Class) String() string {
	return c.statement.name.lexeme
}

func (c *Class) Call(i *Interpreter, paren *Token, arguments []interface{}) (interface{}, error) {
	_, err := c.BaseCallable.Call(i, paren, arguments)
	if err != nil {
		return nil, err
	}
	return NewInstance(c), nil
}

func (c *Class) Get(method *Token) (Callable, bool) {
	v, ok := c.methods[method.lexeme]
	return v, ok
}

// NewInstance constructor
func NewInstance(class *Class) *Instance {
	return &Instance{class: class, properties: map[string]interface{}{}}
}

// Instance representation
type Instance struct {
	class      *Class
	properties map[string]interface{}
}

func (i *Instance) String() string {
	return i.class.statement.name.lexeme + " instance"
}

func (i *Instance) Get(property *Token) (interface{}, error) {
	v, ok := i.properties[property.lexeme]
	if ok {
		return v, nil
	}

	if m, ok := i.class.Get(property); ok {
		return m, nil
	}

	return nil, InvalidProperty(property)
}

func (i *Instance) Set(property *Token, value interface{}) {
	i.properties[property.lexeme] = value
}
