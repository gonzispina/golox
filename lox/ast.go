package lox

// Expression representation
type Expression interface {
	Accept(v ExpressionVisitor) (interface{}, error)
}

// ExpressionVisitor defines the visit method of every Expression
type ExpressionVisitor interface {
	visitBinary(e *Binary) (interface{}, error)
	visitGrouping(e *Grouping) (interface{}, error)
	visitLiteral(e *Literal) (interface{}, error)
	visitUnary(e *Unary) (interface{}, error)
}

// NewLiteral Expression constructor
func NewLiteral(value interface{}) *Literal {
	return &Literal{
		value: value,
	}
}

// Literal Expression implementation
type Literal struct {
	value interface{}
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Literal) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitLiteral(e)
}

// NewUnary Expression constructor
func NewUnary(operator *Token, right Expression) *Unary {
	return &Unary{
		operator: operator,
		right: right,
	}
}

// Unary Expression implementation
type Unary struct {
	operator *Token
	right Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Unary) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitUnary(e)
}

// NewBinary Expression constructor
func NewBinary(left Expression, operator *Token, right Expression) *Binary {
	return &Binary{
		left: left,
		operator: operator,
		right: right,
	}
}

// Binary Expression implementation
type Binary struct {
	left Expression
	operator *Token
	right Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Binary) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitBinary(e)
}

// NewGrouping Expression constructor
func NewGrouping(expression Expression) *Grouping {
	return &Grouping{
		expression: expression,
	}
}

// Grouping Expression implementation
type Grouping struct {
	expression Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Grouping) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitGrouping(e)
}

