package lox

// Expression representation
type Expression interface {
	Accept(v ExpressionVisitor) (interface{}, error)
}

// ExpressionVisitor defines the visit method of every Expression
type ExpressionVisitor interface {
	visitGrouping(e *Grouping) (interface{}, error)
	visitLiteral(e *Literal) (interface{}, error)
	visitThis(e *This) (interface{}, error)
	visitUnary(e *Unary) (interface{}, error)
	visitAssign(e *Assign) (interface{}, error)
	visitBinary(e *Binary) (interface{}, error)
	visitCall(e *Call) (interface{}, error)
	visitGet(e *Get) (interface{}, error)
	visitSet(e *Set) (interface{}, error)
	visitLogical(e *Logical) (interface{}, error)
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

// NewThis Expression constructor
func NewThis(keyword *Token) *This {
	return &This{
		keyword: keyword,
	}
}

// This Expression implementation
type This struct {
	keyword *Token
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *This) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitThis(e)
}

// NewGet Expression constructor
func NewGet(object Expression, name *Token) *Get {
	return &Get{
		object: object,
		name: name,
	}
}

// Get Expression implementation
type Get struct {
	object Expression
	name *Token
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Get) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitGet(e)
}

// NewSet Expression constructor
func NewSet(object Expression, name *Token, value Expression) *Set {
	return &Set{
		object: object,
		name: name,
		value: value,
	}
}

// Set Expression implementation
type Set struct {
	object Expression
	name *Token
	value Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Set) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitSet(e)
}

// NewLogical Expression constructor
func NewLogical(left Expression, operator *Token, right Expression) *Logical {
	return &Logical{
		left: left,
		operator: operator,
		right: right,
	}
}

// Logical Expression implementation
type Logical struct {
	left Expression
	operator *Token
	right Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Logical) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitLogical(e)
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

// NewCall Expression constructor
func NewCall(callee Expression, paren *Token, arguments []Expression) *Call {
	return &Call{
		callee: callee,
		paren: paren,
		arguments: arguments,
	}
}

// Call Expression implementation
type Call struct {
	callee Expression
	paren *Token
	arguments []Expression
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *Call) Accept(v ExpressionVisitor) (interface{}, error) {
	return v.visitCall(e)
}

// Stmt representation
type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}

// StmtVisitor defines the visit method of every Stmt
type StmtVisitor interface {
	visitClassStmt(e *ClassStmt) (interface{}, error)
	visitExpressionStmt(e *ExpressionStmt) (interface{}, error)
	visitFunctionStmt(e *FunctionStmt) (interface{}, error)
	visitVarStmt(e *VarStmt) (interface{}, error)
	visitBlockStmt(e *BlockStmt) (interface{}, error)
	visitIfStmt(e *IfStmt) (interface{}, error)
	visitForStmt(e *ForStmt) (interface{}, error)
	visitPrintStmt(e *PrintStmt) (interface{}, error)
	visitCircuitBreakStmt(e *CircuitBreakStmt) (interface{}, error)
}

// NewForStmt Stmt constructor
func NewForStmt(initializer Stmt, condition Expression, increment Expression, body *BlockStmt, br *bool, cont *bool) *ForStmt {
	return &ForStmt{
		initializer: initializer,
		condition: condition,
		increment: increment,
		body: body,
		br: br,
		cont: cont,
	}
}

// ForStmt Stmt implementation
type ForStmt struct {
	initializer Stmt
	condition Expression
	increment Expression
	body *BlockStmt
	br *bool
	cont *bool
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *ForStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.visitForStmt(e)
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

// NewCircuitBreakStmt Stmt constructor
func NewCircuitBreakStmt(value *bool, statement Stmt) *CircuitBreakStmt {
	return &CircuitBreakStmt{
		value: value,
		statement: statement,
	}
}

// CircuitBreakStmt Stmt implementation
type CircuitBreakStmt struct {
	value *bool
	statement Stmt
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *CircuitBreakStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.visitCircuitBreakStmt(e)
}

// NewIfStmt Stmt constructor
func NewIfStmt(expression Expression, thenBranch *BlockStmt, elseBranch *BlockStmt) *IfStmt {
	return &IfStmt{
		expression: expression,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

// IfStmt Stmt implementation
type IfStmt struct {
	expression Expression
	thenBranch *BlockStmt
	elseBranch *BlockStmt
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *IfStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.visitIfStmt(e)
}

// NewFunctionStmt Stmt constructor
func NewFunctionStmt(name *Token, params []*Token, body *BlockStmt, rt *bool) *FunctionStmt {
	return &FunctionStmt{
		name: name,
		params: params,
		body: body,
		rt: rt,
	}
}

// FunctionStmt Stmt implementation
type FunctionStmt struct {
	name *Token
	params []*Token
	body *BlockStmt
	rt *bool
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *FunctionStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.visitFunctionStmt(e)
}

// NewVarStmt Stmt constructor
func NewVarStmt(name *Token, initializer Stmt) *VarStmt {
	return &VarStmt{
		name: name,
		initializer: initializer,
	}
}

// VarStmt Stmt implementation
type VarStmt struct {
	name *Token
	initializer Stmt
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

// NewClassStmt Stmt constructor
func NewClassStmt(name *Token, super *Variable, methods []*FunctionStmt) *ClassStmt {
	return &ClassStmt{
		name: name,
		super: super,
		methods: methods,
	}
}

// ClassStmt Stmt implementation
type ClassStmt struct {
	name *Token
	super *Variable
	methods []*FunctionStmt
}

// Accept method of the visitor pattern it calls the proper visit method
func(e *ClassStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.visitClassStmt(e)
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

