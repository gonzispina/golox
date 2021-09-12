package lox

import (
	"fmt"
)

const (
	// ErrUnexpectedTokenCode error
	ErrUnexpectedTokenCode = "UnexpectedToken"
	// UnterminatedStringCode error
	UnterminatedStringCode = "UnterminatedString"
	// UnhandledTokenCode error
	UnhandledTokenCode = "UnhandledToken"
	// UnclosedParenthesisCode error
	UnclosedParenthesisCode = "UnclosedParenthesis"
	// ExpectedSemicolonCode error
	ExpectedSemicolonCode = "ExpectedSemicolon"
	// ExpectedIdentifierCode error
	ExpectedIdentifierCode = "ExpectedIdentifier"
	// ExpectedOpeningBraceCode error
	ExpectedOpeningBraceCode = "ExpectedOpeningBrace"
	// ExpectedEndingBraceCode error
	ExpectedEndingBraceCode = "ExpectedEndingBrace"
	// BreakStatementOutsideLoopCode error
	BreakStatementOutsideLoopCode = "BreakStatementOutsideLoop"
	// ContinueStatementOutsideLoopCode error
	ContinueStatementOutsideLoopCode = "ContinueStatementOutsideLoop"

	// ErrInvalidDataTypeCode error
	ErrInvalidDataTypeCode = "InvalidDataType"
	// ErrInvalidOperationCode error
	ErrInvalidOperationCode = "InvalidOperation"
	// ErrDivisionByZeroCode error
	ErrDivisionByZeroCode = "DivisionByZero"
	// ErrUndefinedVariableCode error
	ErrUndefinedVariableCode = "UndefinedVariable"
	// ErrInvalidTargetCode error
	ErrInvalidTargetCode = "InvalidTarget"
)

// Error representation
type Error struct {
	description string
	code        string
	line        *int
	column      *int
}

// Error string
func (e *Error) Error(t string) string {
	on := ""
	if e.line != nil && e.column != nil {
		on += fmt.Sprintf("[Line: %v, Column: %v] ", *e.line, *e.column)
	}
	return fmt.Sprintf("%s %s: %s. Code %v", t, on, e.description, e.code)
}

// SyntaxError representation
type SyntaxError struct {
	err Error
}

func (e *SyntaxError) Error() string {
	return e.err.Error("SyntaxError")
}

// UnexpectedTokenError error
func UnexpectedTokenError(t string, line, column int) *SyntaxError {
	return &SyntaxError{
		Error{
			description: fmt.Sprintf("unexpected token %s", t),
			code:        ErrUnexpectedTokenCode,
			line:        &line,
			column:      &column,
		},
	}
}

// UnterminatedStringError error
func UnterminatedStringError(line, column int) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "unterminated string",
			code:        UnterminatedStringCode,
			line:        &line,
			column:      &column,
		},
	}
}

// UnhandledTokenError error
func UnhandledTokenError(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: fmt.Sprintf("unhandled token %s", t.lexeme),
			code:        UnhandledTokenCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// UnclosedParenthesisError error
func UnclosedParenthesisError(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "parenthesis is not closed",
			code:        UnclosedParenthesisCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// ExpectedSemicolonError error
func ExpectedSemicolonError(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "expected semicolon",
			code:        ExpectedSemicolonCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// ExpectedIdentifier error
func ExpectedIdentifier(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "expected identifier",
			code:        ExpectedIdentifierCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// ExpectedOpeningBrace error
func ExpectedOpeningBrace(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "expected opening brace",
			code:        ExpectedOpeningBraceCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// ExpectedEndingBrace error
func ExpectedEndingBrace(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "expected ending brace",
			code:        ExpectedEndingBraceCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// BreakStatementOutsideLoop error
func BreakStatementOutsideLoop(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "break statements must be inside a for block",
			code:        BreakStatementOutsideLoopCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// ContinueStatementOutsideLoop error
func ContinueStatementOutsideLoop(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "continue statements must be inside a for block",
			code:        ContinueStatementOutsideLoopCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// RuntimeError representation
type RuntimeError struct {
	err Error
}

func (e *RuntimeError) Error() string {
	return e.err.Error("RuntimeError")
}

// InvalidDataTypeError raises when the interpreter receives an unexpected data type
func InvalidDataTypeError(t *Token, got dataType, expected dataType) *RuntimeError {
	return &RuntimeError{
		Error{
			description: fmt.Sprintf("expected %s, got %s", expected, got),
			code:        ErrInvalidDataTypeCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// InvalidOperationError raises when the interpreter receives an uncomputable operation between data types
func InvalidOperationError(t *Token, left dataType, right dataType) *RuntimeError {
	return &RuntimeError{
		Error{
			description: fmt.Sprintf("invalid operation between %s and %s", right, left),
			code:        ErrInvalidOperationCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// DivisionByZeroError raises when it tries to divide by zero
func DivisionByZeroError(t *Token) *RuntimeError {
	return &RuntimeError{
		Error{
			description: "division by zero is not supported",
			code:        ErrDivisionByZeroCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// UndefinedVariable raises when an undefined variable is called
func UndefinedVariable(name string, t *Token) *RuntimeError {
	return &RuntimeError{
		Error{
			description: fmt.Sprintf("undefined variable '%s'", name),
			code:        ErrUndefinedVariableCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// InvalidTarget raises when an assignment bad targeted
func InvalidTarget(t *Token) *RuntimeError {
	return &RuntimeError{
		Error{
			description: "invalid assignment target",
			code:        ErrInvalidTargetCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}
