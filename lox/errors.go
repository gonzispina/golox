package lox

import "errors"

var (
	// ErrUnexpectedCharacter error
	ErrUnexpectedCharacter = errors.New("unexpected character")
	// ErrUnterminatedString error
	ErrUnterminatedString = errors.New("unterminated string")
)
