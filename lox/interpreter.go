package lox

import (
	"fmt"
)

// NewInterpreter constructor
func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

// Interpreter of the lox language
type Interpreter struct {
}

// Interpret the given expression
func (i *Interpreter) Interpret(e Expression) (string, error) {
	v, err := e.Accept(i)
	if err != nil {
		return "", err
	}
	return i.stringify(v), nil
}

func (i *Interpreter) stringify(v interface{}) string {
	dt := getDataType(v)
	if dt == object {
		return "nil"
	} else if dt == boolean {
		if v.(bool) {
			return "true"
		} else {
			return "false"
		}
	}
	return fmt.Sprintf("%s", v)
}

func (i *Interpreter) evaluate(e Expression) (interface{}, error) {
	return e.Accept(i)
}

func (i *Interpreter) visitGrouping(e *Grouping) (interface{}, error) {
	return i.evaluate(e.expression)
}

func (i *Interpreter) visitLiteral(e *Literal) (interface{}, error) {
	return e.value, nil
}

func (i *Interpreter) visitUnary(e *Unary) (interface{}, error) {
	right, err := i.evaluate(e.right)
	if err != nil {
		return nil, err
	}

	switch e.operator.tokenType {
	case MINUS:
		v, ok := right.(float64)
		if !ok {
			return nil, InvalidDataTypeError(e.operator, getDataType(right), number)
		}
		return -v, nil
	case BANG:
		return i.isTruthy(right), nil
	default:
		return nil, nil
	}
}

func (i *Interpreter) isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) visitBinary(e *Binary) (interface{}, error) {
	left, err := i.evaluate(e.left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(e.right)
	if err != nil {
		return nil, err
	}

	switch e.operator.tokenType {
	case GREATER:
		return greaterThan(left, right, e.operator)
	case GREATER_EQUAL:
		return greaterEqual(left, right, e.operator)
	case LESS:
		return lesserThan(left, right, e.operator)
	case LESS_EQUAL:
		return lesserEqual(left, right, e.operator)
	case EQUAL_EQUAL:
		return isEqual(left, right, e.operator)
	case BANG_EQUAL:
		return notEqual(left, right, e.operator)
	case MINUS:
		v1, v2, err := both2Float(left, right, e.operator)
		if err != nil {
			return nil, err
		}

		return v1 - v2, nil
	case STAR:
		v1, v2, err := both2Float(left, right, e.operator)
		if err != nil {
			return nil, err
		}

		return v1 * v2, nil
	case SLASH:
		v1, v2, err := both2Float(left, right, e.operator)
		if err != nil {
			return nil, err
		}

		if v2 == 0 {
			return nil, DivisionByZeroError(e.operator)
		}

		return v1 / v2, nil
	case PLUS:
		return addValues(left, right, e.operator)
	}

	return nil, nil
}
