package lox

import (
	"fmt"
)

// NewInterpreter constructor
func NewInterpreter() *Interpreter {
	e := NewEnvironment()
	return &Interpreter{environment: e}
}

// Interpreter of the lox language
type Interpreter struct {
	environment *Environment
}

// Interpret the given expression
func (i *Interpreter) Interpret(s []Stmt) error {
	for _, stmt := range s {
		_, err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
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
	return fmt.Sprintf("%v", v)
}

func (i *Interpreter) execute(s Stmt) (interface{}, error) {
	return s.Accept(i)
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
		return isTruthy(right), nil
	default:
		return nil, nil
	}
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

func (i *Interpreter) visitVariable(e *Variable) (interface{}, error) {
	v, ok := i.environment.get(e.token.lexeme)
	if !ok {
		return nil, UndefinedVariable(e.token.lexeme, e.token)
	}
	return v, nil
}

func (i *Interpreter) visitPrintStmt(s *PrintStmt) (interface{}, error) {
	value, err := i.evaluate(s.expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(i.stringify(value))
	return nil, nil
}

func (i *Interpreter) visitExpressionStmt(s *ExpressionStmt) (interface{}, error) {
	v, err := i.evaluate(s.expression)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (i *Interpreter) visitVarStmt(e *VarStmt) (interface{}, error) {
	value, err := i.evaluate(e.initializer)
	if err != nil {
		return nil, err
	}

	i.environment.define(e.token.lexeme, value)
	return nil, nil
}

// assignment → IDENTIFIER "=" assignment | equality ;
func (i *Interpreter) visitAssign(e *Assign) (interface{}, error) {
	value, err := i.evaluate(e.value)
	if err != nil {
		return nil, err
	}

	ok := i.environment.assign(e.name.lexeme, value)
	if !ok {
		return nil, UndefinedVariable(e.name.lexeme, e.name)
	}

	return nil, nil
}
