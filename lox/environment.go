package lox

// NewEnvironment constructor
func NewEnvironment() *Environment {
	return &Environment{values: map[string]interface{}{}}
}

// Environment representation
type Environment struct {
	values map[string]interface{}
}

func (e *Environment) define(s string, value interface{}) {
	e.values[s] = value
}

func (e *Environment) get(s string) (interface{}, bool) {
	v, ok := e.values[s]
	return v, ok
}

func (e *Environment) assign(s string, value interface{}) bool {
	_, ok := e.values[s]
	if !ok {
		return false
	}

	e.values[s] = value
	return true
}
