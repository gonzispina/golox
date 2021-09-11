package lox

import (
	"fmt"
)

// NewASTPrinter constructor
func NewASTPrinter() *ASTPrinter {
	return &ASTPrinter{}
}

// ASTPrinter visitor. Traverses the whole tree and creates a string representation of the tree
type ASTPrinter struct {
}

func (p *ASTPrinter) visitAssign(e *Assign) (interface{}, error) {
	panic("implement me")
}

func (p *ASTPrinter) visitVariable(e *Variable) (interface{}, error) {
	panic("implement me")
}

// Print the given expression
func (p *ASTPrinter) Print(e Expression) (string, error) {
	v, err := e.Accept(p)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (p *ASTPrinter) visitBinary(e *Binary) (interface{}, error) {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p *ASTPrinter) visitGrouping(e *Grouping) (interface{}, error) {
	return p.parenthesize("group", e.expression)
}

func (p *ASTPrinter) visitLiteral(e *Literal) (interface{}, error) {
	if e.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", e.value), nil
}

func (p *ASTPrinter) visitUnary(e *Unary) (interface{}, error) {
	return p.parenthesize(e.operator.lexeme, e.right)
}

func (p *ASTPrinter) parenthesize(name string, expressions ...Expression) (interface{}, error) {
	line := fmt.Sprintf("(%s", name)
	for _, e := range expressions {
		line += " "
		v, err := e.Accept(p)
		if err != nil {
			return nil, err
		}
		line += v.(string)
	}
	line += ")"
	return line, nil
}
