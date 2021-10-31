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
		parameters:  parameters,
		environment: NewEnvironment(closure),
	}
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

func (f *Function) Bind(this *Instance) *Function {
	method := &Function{
		BaseCallable: f.BaseCallable,
		statement:    f.statement,
	}
	method.environment.define("this", this)
	return method
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
func NewClass(statement *ClassStmt, methods map[string]*Function) *Class {
	return &Class{
		statement: statement,
		methods:   methods,
	}
}

// Class representation
type Class struct {
	statement *ClassStmt
	methods   map[string]*Function
}

func (c *Class) String() string {
	return c.statement.name.lexeme
}

func (c *Class) Call(i *Interpreter, paren *Token, arguments []interface{}) (interface{}, error) {
	return NewInstance(c, i, paren, arguments)
}

// NewInstance constructor
func NewInstance(class *Class, i *Interpreter, paren *Token, arguments []interface{}) (*Instance, error) {
	instance := &Instance{
		class:      class,
		properties: map[string]interface{}{},
		methods:    map[string]*Function{},
	}

	for name, method := range class.methods {
		instance.methods[name] = method.Bind(instance)
	}

	if init, ok := instance.methods["init"]; ok {
		_, err := init.Call(i, paren, arguments)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

// Instance representation
type Instance struct {
	class      *Class
	properties map[string]interface{}
	methods    map[string]*Function
}

func (i *Instance) String() string {
	return i.class.statement.name.lexeme + " instance"
}

func (i *Instance) Get(property *Token) (interface{}, error) {
	v, ok := i.properties[property.lexeme]
	if ok {
		return v, nil
	}

	if m, ok := i.methods[property.lexeme]; ok {
		return m, nil
	}

	return nil, InvalidProperty(property)
}

func (i *Instance) Set(property *Token, value interface{}) {
	i.properties[property.lexeme] = value
}
