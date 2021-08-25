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
		return nil
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

func (p *Parser) isAtEnd() bool {
	return p.current().Is(EOF)
}

func (p *Parser) Parse() (Expression, error) {
	return p.expression()
}

func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		if p.advance().OneOf(
			CLASS, FUN,
			VAR, FOR,
			IF, WHILE,
			PRINT, RETURN,
		) {
			break
		}
	}
}

// expression → equality
func (p *Parser) expression() (Expression, error) {
	return p.equality()
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

// unary → ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() (Expression, error) {
	if p.current().OneOf(BANG, MINUS) {
		operator := p.advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return NewUnary(operator, right), nil
	}

	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() (Expression, error) {
	if p.current().OneOf(NIL, TRUE, FALSE, NUMBER, STRING) {
		return NewLiteral(p.advance().literal), nil
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
