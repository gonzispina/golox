package lox

import "errors"

var (
	// ErrUnexpectedToken error
	ErrUnexpectedToken = errors.New("unexpected character")
	// ErrUnterminatedString error
	ErrUnterminatedString = errors.New("unterminated string")
	// ErrUnhandledToken error
	ErrUnhandledToken = errors.New("unhandled token")
	// ErrUnclosedParenthesis error
	ErrUnclosedParenthesis = errors.New("parenthesis is not closed")
)

// NewParserError representation
func NewParserError(err error, token *Token) *ParserError {
	return &ParserError{err: err, token: token}
}

// ParserError representation
type ParserError struct {
	err   error
	token *Token
}

// Error ...
func (e *ParserError) Error() string {
	return "parser error: " + e.err.Error()
}
