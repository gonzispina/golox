package lox

// NewResolver constructor
func NewResolver(i *Interpreter) *Resolver {
	return &Resolver{
		interpreter: i,
		scopes:      NewScopeStack(),
	}
}

// Resolver is in charge of doing the semantic analysis for variable resolving before the program is
// interpreted. Since it visits every node of the syntax tree once, its time complexity is O(n)
// where n is the amount of nodes in the tree.
type Resolver struct {
	interpreter *Interpreter
	scopes      *ScopeStack
	inClass     bool
}

// Resolve API
func (r *Resolver) Resolve(stmts []Stmt) (interface{}, error) {
	for _, s := range stmts {
		_, err := r.resolveStatement(s)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) resolveStatement(s Stmt) (interface{}, error) {
	return s.Accept(r)
}

func (r *Resolver) resolveExpression(e Expression) (interface{}, error) {
	return e.Accept(r)
}

func (r *Resolver) resolveLocal(e Expression, name *Token) (interface{}, error) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		s, err := r.scopes.Get(uint(i))
		if err != nil {
			break
		}

		v, ok := s[name.lexeme]
		if ok {
			r.interpreter.Resolve(e, r.scopes.Size()-i-1)
			v.used = true
			break
		}
	}

	return nil, nil
}

func (r *Resolver) resolveFunction(s *FunctionStmt) (interface{}, error) {
	r.beginScope()
	for _, param := range s.params {
		r.declare(param)
		r.define(param)
	}

	v, err := r.Resolve(s.body.statements)
	if err != nil {
		return nil, err
	}

	if err := r.endScope(); err != nil {
		return nil, err
	}

	return v, nil
}

func (r *Resolver) beginScope() map[string]*ScopeEntry {
	scope := map[string]*ScopeEntry{}
	r.scopes.Push(scope)
	return scope
}

func (r *Resolver) endScope() error {
	s, err := r.scopes.Pop()
	if err != nil {
		return nil
	}

	for _, entry := range s {
		if !entry.used {
			return UnusedVariable(entry.token)
		}
	}

	return nil
}

func (r *Resolver) declare(t *Token) {
	s, err := r.scopes.Peek()
	if err != nil {
		return
	}

	s[t.lexeme] = &ScopeEntry{token: t}
	return
}

func (r *Resolver) define(t *Token) {
	s, err := r.scopes.Peek()
	if err != nil {
		return
	}

	s[t.lexeme].defined = true
	return
}

func (r *Resolver) visitBinary(e *Binary) (interface{}, error) {
	_, err := r.resolveExpression(e.left)
	if err != nil {
		return nil, err
	}

	return r.resolveExpression(e.right)
}

func (r *Resolver) visitCall(e *Call) (interface{}, error) {
	_, err := r.resolveExpression(e.callee)
	if err != nil {
		return nil, err
	}

	for _, argument := range e.arguments {
		_, err = r.resolveExpression(argument)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) visitGrouping(e *Grouping) (interface{}, error) {
	return r.resolveExpression(e.expression)
}

func (r *Resolver) visitLogical(e *Logical) (interface{}, error) {
	_, err := r.resolveExpression(e.left)
	if err != nil {
		return nil, err
	}

	return r.resolveExpression(e.right)
}

func (r *Resolver) visitLiteral(_ *Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) visitUnary(e *Unary) (interface{}, error) {
	return r.resolveExpression(e.right)
}

func (r *Resolver) visitVariable(e *Variable) (interface{}, error) {
	s, err := r.scopes.Peek()
	if err == nil {
		if entry, ok := s[e.token.lexeme]; ok && !entry.defined {
			return nil, InitializerSelfReference(e.token)
		}
	}

	return r.resolveLocal(e, e.token)
}

func (r *Resolver) visitAssign(e *Assign) (interface{}, error) {
	_, err := r.resolveExpression(e.value)
	if err != nil {
		return nil, err
	}

	return r.resolveLocal(e, e.name)
}

func (r *Resolver) visitGet(e *Get) (interface{}, error) {
	return r.resolveExpression(e.object)
}

func (r *Resolver) visitSet(e *Set) (interface{}, error) {
	_, err := r.resolveExpression(e.value)
	if err != nil {
		return nil, err
	}
	return r.resolveExpression(e.object)
}

func (r *Resolver) visitThis(e *This) (interface{}, error) {
	if !r.inClass {
		return nil, ThisOutsideClass(e.keyword)
	}
	return r.resolveLocal(e, e.keyword)
}

func (r *Resolver) visitIfStmt(e *IfStmt) (interface{}, error) {
	_, err := r.resolveExpression(e.expression)
	if err != nil {
		return nil, err
	}

	_, err = r.Resolve(e.thenBranch.statements)
	if err != nil {
		return nil, err
	}

	if e.elseBranch == nil {
		return nil, nil
	}

	return r.Resolve(e.elseBranch.statements)
}

func (r *Resolver) visitForStmt(e *ForStmt) (interface{}, error) {
	r.beginScope()
	if e.initializer != nil {
		_, err := r.resolveStatement(e.initializer)
		if err != nil {
			return nil, err
		}
	}

	_, err := r.resolveExpression(e.condition)
	if err != nil {
		return nil, err
	}

	v, err := r.Resolve(e.body.statements)
	if err != nil {
		return nil, err
	}

	if err := r.endScope(); err != nil {
		return nil, err
	}

	return v, nil
}

func (r *Resolver) visitPrintStmt(e *PrintStmt) (interface{}, error) {
	return r.resolveExpression(e.expression)
}

func (r *Resolver) visitVarStmt(e *VarStmt) (interface{}, error) {
	s, err := r.scopes.Peek()
	if err != nil {
		return nil, err
	}

	if _, ok := s[e.name.lexeme]; ok {
		return nil, VariableAlreadyDeclared(e.name)
	}

	r.declare(e.name)
	if e.initializer != nil {
		_, err := r.resolveStatement(e.initializer)
		if err != nil {
			return nil, err
		}
	}
	r.define(e.name)
	return nil, nil
}

func (r *Resolver) visitBlockStmt(e *BlockStmt) (interface{}, error) {
	r.beginScope()

	v, err := r.Resolve(e.statements)
	if err != nil {
		return nil, err
	}

	if err := r.endScope(); err != nil {
		return nil, err
	}

	return v, nil
}

func (r *Resolver) visitCircuitBreakStmt(e *CircuitBreakStmt) (interface{}, error) {
	if e.statement == nil {
		return nil, nil
	}
	return r.resolveStatement(e.statement)
}

func (r *Resolver) visitExpressionStmt(e *ExpressionStmt) (interface{}, error) {
	return r.resolveExpression(e.expression)
}

func (r *Resolver) visitFunctionStmt(e *FunctionStmt) (interface{}, error) {
	if e.name != nil {
		// If is not a lambda function
		r.declare(e.name)
		r.define(e.name)
	}

	return r.resolveFunction(e)
}

func (r *Resolver) visitClassStmt(e *ClassStmt) (interface{}, error) {
	r.declare(e.name)
	r.define(e.name)

	r.beginScope()
	r.inClass = true
	defer func() {
		r.inClass = false
	}()

	for _, method := range e.methods {
		_, err := r.resolveFunction(method)
		if err != nil {
			return nil, err
		}
	}

	if err := r.endScope(); err != nil {
		return nil, err
	}

	return nil, nil
}
