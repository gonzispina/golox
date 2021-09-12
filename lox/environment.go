package lox

// NewEnvironment constructor
func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    map[string]interface{}{},
		enclosing: enclosing,
	}
}

// Environment representation
type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func (e *Environment) define(s string, value interface{}) {
	e.values[s] = value
}

func (e *Environment) get(s string) (interface{}, bool) {
	v, ok := e.values[s]
	if ok {
		return v, ok
	}

	if e.enclosing == nil {
		return nil, false
	}

	return e.enclosing.get(s)
}

func (e *Environment) assign(s string, value interface{}) bool {
	_, ok := e.values[s]
	if ok {
		e.values[s] = value
		return true
	}

	if e.enclosing == nil {
		return false
	}

	return e.enclosing.assign(s, value)
}
