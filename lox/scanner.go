package lox

import (
	"strconv"
)

func NewScanner(source string) *Scanner {
	return &Scanner{
		iterator: newIterator([]rune(source)),
		tokens:   []*Token{},
	}
}

type Scanner struct {
	tokens   []*Token
	iterator *Iterator
}

func (s *Scanner) ScanTokens() ([]*Token, error) {
	if len(s.iterator.source) == 0 {
		return []*Token{}, nil
	}

	for !(s.iterator.isAtEnd()) {
		s.iterator.start = s.iterator.current
		if err := s.scanToken(); err != nil {
			report(s.iterator.line, "", err.Error())
			return nil, err
		}
	}

	s.addTokenByType(EOF)
	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	var t TokenType
	var c = s.iterator.advance()
	switch c {
	case '(':
		s.addTokenByType(LEFT_PAREN)
		break
	case ')':
		s.addTokenByType(RIGHT_PAREN)
		break
	case '{':
		s.addTokenByType(LEFT_BRACE)
		break
	case '}':
		s.addTokenByType(RIGHT_BRACE)
		break
	case ',':
		s.addTokenByType(COMMA)
		break
	case '.':
		s.addTokenByType(DOT)
		break
	case '-':
		s.addTokenByType(MINUS)
		break
	case '+':
		s.addTokenByType(PLUS)
		break
	case ';':
		s.addTokenByType(SEMICOLON)
		break
	case '*':
		if s.iterator.match('/') {
			return UnexpectedTokenError(string(s.iterator.peek()), s.iterator.line, s.iterator.column)
		}
		s.addTokenByType(STAR)
		break
	case '!':
		if s.iterator.match('=') {
			t = BANG_EQUAL
		} else {
			t = BANG
		}
		s.addTokenByType(t)
		break
	case '=':
		if s.iterator.match('=') {
			t = EQUAL_EQUAL
		} else {
			t = EQUAL
		}
		s.addTokenByType(t)
		break
	case '<':
		if s.iterator.match('=') {
			t = LESS_EQUAL
		} else {
			t = LESS
		}
		s.addTokenByType(t)
		break
	case '>':
		if s.iterator.match('=') {
			t = GREATER_EQUAL
		} else {
			t = GREATER
		}
		s.addTokenByType(t)
		break
	case '/':
		s.comment()
		break
	case '|':
		if s.iterator.match('|') {
			s.addTokenByType(OR)
		} else {
			return UnexpectedTokenError(string(s.iterator.next()), s.iterator.line, s.iterator.column+1)
		}
	case '&':
		if s.iterator.match('&') {
			s.addTokenByType(AND)
		} else {
			return UnexpectedTokenError(string(s.iterator.next()), s.iterator.line, s.iterator.column+1)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		break
	case '"':
		if err := s.string(); err != nil {
			return err
		}
		break
	default:
		if isDigit(c) {
			if err := s.number(); err != nil {
				return err
			}
			break
		} else if isAlpha(c) {
			s.identifier()
			break
		}

		return UnexpectedTokenError(string(s.iterator.peek()), s.iterator.line, s.iterator.column)
	}

	return nil
}

func (s *Scanner) comment() {
	if s.iterator.match('/') {
		for s.iterator.peek() != '\n' && !s.iterator.isAtEnd() {
			s.iterator.advance()
		}
	} else if s.iterator.match('*') {
		for !s.iterator.isAtEnd() {
			if s.iterator.advance() == '*' && s.iterator.advance() == '/' {
				break
			}
		}
	} else {
		s.addTokenByType(SLASH)
	}
}

func (s *Scanner) number() error {
	for isDigit(s.iterator.peek()) {
		s.iterator.advance()
	}

	if s.iterator.peek() == '.' && isDigit(s.iterator.next()) {
		s.iterator.advance()

		for isDigit(s.iterator.peek()) {
			s.iterator.advance()
		}
	}

	v, err := strconv.ParseFloat(string(s.iterator.source[s.iterator.start:s.iterator.current]), 64)
	if err != nil {
		return err
	}

	s.addToken(NUMBER, v)
	return nil
}

func (s *Scanner) identifier() {
	for isAlphanumeric(s.iterator.peek()) {
		s.iterator.advance()
	}

	token, ok := Reserved[s.iterator.currentLexeme()]
	if !ok {
		token = IDENTIFIER
	}

	s.addTokenByType(token)
}

func (s *Scanner) string() error {
	for s.iterator.peek() != '"' && !s.iterator.isAtEnd() {
		s.iterator.advance()
	}

	if s.iterator.isAtEnd() {
		return UnterminatedStringError(s.iterator.line, s.iterator.column)
	}

	s.iterator.advance()
	s.addToken(STRING, s.iterator.source[s.iterator.start+1:s.iterator.current-1])
	return nil
}

func (s *Scanner) addTokenByType(t TokenType) {
	s.addToken(t, nil)
}

func (s *Scanner) addToken(t TokenType, literal interface{}) {
	s.tokens = append(s.tokens, &Token{
		tokenType: t,
		lexeme:    s.iterator.currentLexeme(),
		literal:   literal,
		line:      s.iterator.line,
	})
}
