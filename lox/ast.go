package lox

// Expression representation
type Expression interface {
	Accept(v ExpressionVisitor) (interface{}, error)
}

// ExpressionVisitor defines the visit method of every Expression
type ExpressionVisitor interface {
	visitAssign(e *Assign) (interface{}, error)
	visitBinary(e *Binary) (interface{}, error)
	visitGrouping(e *Grouping) (interface{}, error)
	visitLiteral(e *Literal) (interface{}, error)
	visitUnary(e *Unary) (interface{}, error)
	visitVariable(e *Variable) (interface{}, error)
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

// NewVariable Expression constructor
func NewVariable(token *Token) *Variable {
	return &Variable{
		token: token,
	}
}

// Variable Expression implementation
type Variable struct {
	token *Token
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Variable) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitVariable(e)
}

// NewAssign Expression constructor
func NewAssign(name *Token, value Expression) *Assign {
	return &Assign{
		name: name,
		value: value,
	}
}

// Assign Expression implementation
type Assign struct {
	name *Token
	value Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Assign) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitAssign(e)
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

// Stmt representation
type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}

// StmtVisitor defines the visit method of every Stmt
type StmtVisitor interface {
	visitExpressionStmt(e *ExpressionStmt) (interface{}, error)
	visitPrintStmt(e *PrintStmt) (interface{}, error)
	visitVarStmt(e *VarStmt) (interface{}, error)
	visitBlockStmt(e *BlockStmt) (interface{}, error)
}

// NewVarStmt Stmt constructor
func NewVarStmt(token *Token, initializer Expression) *VarStmt {
	return &VarStmt{
		token: token,
		initializer: initializer,
	}
}

// VarStmt Stmt implementation
type VarStmt struct {
	token *Token
	initializer Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *VarStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.visitVarStmt(e)
}

// NewBlockStmt Stmt constructor
func NewBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{
		statements: statements,
	}
}

// BlockStmt Stmt implementation
type BlockStmt struct {
	statements []Stmt
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *BlockStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.visitBlockStmt(e)
}

// NewExpressionStmt Stmt constructor
func NewExpressionStmt(expression Expression) *ExpressionStmt {
	return &ExpressionStmt{
		expression: expression,
	}
}

// ExpressionStmt Stmt implementation
type ExpressionStmt struct {
	expression Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *ExpressionStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.visitExpressionStmt(e)
}

// NewPrintStmt Stmt constructor
func NewPrintStmt(expression Expression) *PrintStmt {
	return &PrintStmt{
		expression: expression,
	}
}

// PrintStmt Stmt implementation
type PrintStmt struct {
	expression Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *PrintStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.visitPrintStmt(e)
}

