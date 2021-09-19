package lox

// NewParser constructor
func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens: tokens}
}

// Parser implementation
type Parser struct {
	tokens []*Token
	index  int
}

func (p *Parser) current() *Token {
	return p.tokens[p.index]
}

func (p *Parser) advance() *Token {
	if p.isAtEnd() {
		return p.current()
	}
	p.index++
	return p.previous()
}

func (p *Parser) previous() *Token {
	if p.index == 0 {
		return nil
	}
	return p.tokens[p.index-1]
}

func (p *Parser) match(tt TokenType) bool {
	if p.current().Is(tt) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) isAtEnd() bool {
	return p.current().Is(EOF)
}

func (p *Parser) Parse() ([]Stmt, []error) {
	var s []Stmt
	var errs []error
	for !p.isAtEnd() {
		statement, err := p.declaration(nil, nil)
		if err != nil {
			errs = append(errs, err)
			p.synchronize()
			continue
		}
		s = append(s, statement)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return s, nil
}

func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		if p.advance().OneOf(
			CLASS, FUN,
			VAR, FOR,
			PRINT, RETURN,
			IF,
		) {
			break
		}
	}
}

// declaration → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) declaration(br *CircuitBreakStmt, cont *CircuitBreakStmt) (Stmt, error) {
	if !p.match(VAR) {
		return p.statement(br, cont)
	}

	if !p.match(IDENTIFIER) {
		return nil, ExpectedIdentifier(p.current())
	}

	var err error
	var initializer Expression

	identifier := p.previous()
	initializer = NewLiteral("nil")
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if !p.match(SEMICOLON) {
		return nil, ExpectedSemicolonError(p.current())
	}

	return NewVarStmt(identifier, initializer), nil
}

// statement → exprStmt | forStmt | ifStmt | printStmt | block ;
// exprStmt → expression ";" ;
// forStmt → "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;
// ifStmt → "if" "(" expression ")" statement ( "else" statement )? ;
// printStmt → "print" expression ";"
// block → "{" declaration* "}" ;
func (p *Parser) statement(br *CircuitBreakStmt, cont *CircuitBreakStmt) (Stmt, error) {
	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(IF) {
		return p.ifStatement(br, cont)
	}

	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		return p.blockStatement(br, cont)
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (*ForStmt, error) {
	var err error

	cont := NewCircuitBreakStmt(false)
	br := NewCircuitBreakStmt(false)

	if p.match(LEFT_BRACE) {
		body, err := p.blockStatement(br, cont)
		if err != nil {
			return nil, err
		}
		return NewForStmt(nil, nil, nil, body, br, cont), nil
	}

	var initializer Stmt
	if p.current().Is(VAR) {
		initializer, err = p.declaration(br, cont)
		if err != nil {
			return nil, err
		}
	}

	if p.match(LEFT_BRACE) {
		body, err := p.blockStatement(br, cont)
		if err != nil {
			return nil, err
		}
		return NewForStmt(initializer, nil, nil, body, br, cont), nil
	}

	var conditional Expression
	var increment Expression
	if p.match(SEMICOLON) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	} else {
		conditional, err = p.expression()
		if err != nil {
			return nil, err
		}

		if p.match(SEMICOLON) {
			increment, err = p.expression()
			if err != nil {
				return nil, err
			}
		}
	}

	if !p.match(LEFT_BRACE) {
		return nil, ExpectedOpeningBrace(p.current())
	}

	body, err := p.blockStatement(br, cont)
	if err != nil {
		return nil, err
	}

	return NewForStmt(initializer, conditional, increment, body, br, cont), nil
}

func (p *Parser) ifStatement(br *CircuitBreakStmt, cont *CircuitBreakStmt) (*IfStmt, error) {
	expression, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(LEFT_BRACE) {
		return nil, ExpectedOpeningBrace(p.current())
	}

	thenBranch, err := p.blockStatement(br, cont)
	if err != nil {
		return nil, err
	}

	var elseBranch *BlockStmt
	if p.match(ELSE) {
		if !p.match(LEFT_BRACE) {
			return nil, ExpectedOpeningBrace(p.current())
		}

		elseBranch, err = p.blockStatement(br, cont)
		if err != nil {
			return nil, err
		}
	}

	return NewIfStmt(expression, thenBranch, elseBranch), nil
}

func (p *Parser) printStatement() (*PrintStmt, error) {
	e, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.advance().Is(SEMICOLON) {
		return nil, ExpectedSemicolonError(p.current())
	}

	return NewPrintStmt(e), nil
}

func (p *Parser) expressionStatement() (*ExpressionStmt, error) {
	e, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.advance().Is(SEMICOLON) {
		return nil, ExpectedSemicolonError(p.current())
	}

	return NewExpressionStmt(e), nil
}

func (p *Parser) blockStatement(br *CircuitBreakStmt, cont *CircuitBreakStmt) (*BlockStmt, error) {
	var s []Stmt

	for !p.isAtEnd() && !p.current().Is(RIGHT_BRACE) {
		var statement Stmt
		var err error
		if p.match(BREAK) {
			if br == nil {
				return nil, BreakStatementOutsideLoop(p.previous())
			}

			statement = br
			if !p.match(SEMICOLON) {
				return nil, ExpectedSemicolonError(p.current())
			}
		} else if p.match(CONTINUE) {
			if cont == nil {
				return nil, ContinueStatementOutsideLoop(p.previous())
			}

			statement = cont
			if !p.match(SEMICOLON) {
				return nil, ExpectedSemicolonError(p.current())
			}
		} else {
			statement, err = p.declaration(br, cont)
			if err != nil {
				return nil, err
			}
		}

		s = append(s, statement)
	}

	if !p.match(RIGHT_BRACE) {
		return nil, ExpectedEndingBrace(p.current())
	}

	return NewBlockStmt(s), nil
}

// expression → assignment
func (p *Parser) expression() (Expression, error) {
	return p.assignment()
}

// assignment → IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() (Expression, error) {
	e, err := p.or()
	if err != nil {
		return nil, err
	}

	if !p.match(EQUAL) {
		return e, nil
	}

	equals := p.current()
	value, err := p.assignment()
	if err != nil {
		return nil, err
	}

	variable, ok := e.(*Variable)
	if !ok {
		return nil, InvalidTarget(equals)
	}

	return NewAssign(variable.token, value), nil
}

// or → and ( "or" and )* ;
func (p *Parser) or() (Expression, error) {
	e, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		e = NewLogical(e, operator, right)
	}

	return e, nil
}

// and → equality ( "and" equality )* ;
func (p *Parser) and() (Expression, error) {
	e, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		e = NewLogical(e, operator, right)
	}

	return e, nil
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (Expression, error) {
	expression, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.current().OneOf(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.advance()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expression = NewBinary(expression, operator, right)
	}

	return expression, nil
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (Expression, error) {
	expression, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.current().OneOf(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.advance()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expression = NewBinary(expression, operator, right)
	}

	return expression, nil
}

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (Expression, error) {
	expression, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.current().OneOf(MINUS, PLUS) {
		operator := p.advance()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expression = NewBinary(expression, operator, right)
	}

	return expression, nil
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (Expression, error) {
	expression, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.current().OneOf(SLASH, STAR) {
		operator := p.advance()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expression = NewBinary(expression, operator, right)
	}

	return expression, nil
}

// unary → ( "!" | "-" ) unary | call ;
func (p *Parser) unary() (Expression, error) {
	if p.current().OneOf(BANG, MINUS) {
		operator := p.advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return NewUnary(operator, right), nil
	}

	return p.call()
}

// call → primary ( "(" arguments? ")" )* ;
func (p *Parser) call() (Expression, error) {
	e, err := p.primary()
	if err != nil {
		return nil, err
	}

	for p.match(LEFT_PAREN) {
		var arguments []Expression
		if p.match(RIGHT_PAREN) {
			e = NewCall(e, p.current(), arguments)
			continue
		}

		for {
			if len(arguments) > 255 {
				return nil, ArgumentLimitExceeded(p.previous())
			}

			arg, err := p.expression()
			if err != nil {
				return nil, err
			}

			arguments = append(arguments, arg)
			if p.match(COMMA) {
				continue
			} else if p.match(RIGHT_PAREN) {
				e = NewCall(e, p.current(), arguments)
				break
			} else {
				return nil, BadFunctionSignature(p.current())
			}
		}
	}

	return e, nil
}

// primary → "true" | "false" | "nil" | NUMBER | STRING | "(" expression ")" | IDENTIFIER ;
func (p *Parser) primary() (Expression, error) {
	if p.current().OneOf(NIL, TRUE, FALSE, NUMBER, STRING) {
		return NewLiteral(p.advance().literal), nil
	}

	if p.current().Is(IDENTIFIER) {
		return NewVariable(p.advance()), nil
	}

	if p.current().Is(LEFT_PAREN) {
		p.advance()
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.advance().Is(RIGHT_PAREN) {
			return nil, UnclosedParenthesisError(p.current())
		}

		return NewGrouping(expression), nil
	}

	return nil, UnhandledTokenError(p.current())
}
