package lox

import "time"

// NewClockFunction constructor
func NewClockFunction() *Clock {
	return &Clock{}
}

// Clock native function
type Clock struct{}

func (c *Clock) Call(_ *Interpreter, _ *Token, _ []interface{}) (interface{}, error) {
	return time.Now().UnixNano() / int64(time.Microsecond), nil
}

func (c *Clock) String() string {
	return "<native fn>"
}
