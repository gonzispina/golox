package lox

import (
	"fmt"
	"strings"
)

const (
	// UnexpectedTokenCode error
	UnexpectedTokenCode = "UnexpectedToken"
	// UnterminatedStringCode error
	UnterminatedStringCode = "UnterminatedString"
	// UnhandledTokenCode error
	UnhandledTokenCode = "UnhandledToken"
	// UnclosedParenthesisCode error
	UnclosedParenthesisCode = "UnclosedParenthesis"
	// ExpectedIdentifierCode error
	ExpectedIdentifierCode = "ExpectedIdentifier"
	// BreakStatementOutsideLoopCode error
	BreakStatementOutsideLoopCode = "BreakStatementOutsideLoop"
	// ContinueStatementOutsideLoopCode error
	ContinueStatementOutsideLoopCode = "ContinueStatementOutsideLoop"
	// ReturnStatementOutsideFunctionCode error
	ReturnStatementOutsideFunctionCode = "ReturnStatementOutsideFunction"
	// ArgumentSizeExceededCode error
	ArgumentSizeExceededCode = "ArgumentSizeExceeded"
	// InvalidTargetCode error
	InvalidTargetCode = "InvalidTarget"

	// InvalidDataTypeCode error
	InvalidDataTypeCode = "InvalidDataType"
	// InvalidOperationCode error
	InvalidOperationCode = "InvalidOperation"
	// DivisionByZeroCode error
	DivisionByZeroCode = "DivisionByZero"
	// UndefinedVariableCode error
	UndefinedVariableCode = "UndefinedVariable"
	// ExpressionIsNotCallableCode error
	ExpressionIsNotCallableCode = "ExpressionIsNotCallable"
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

// UnexpectedLexeme error
func UnexpectedLexeme(t rune, line, column int) *SyntaxError {
	return &SyntaxError{
		Error{
			description: fmt.Sprintf("unexpected token '%s'", string(t)),
			code:        UnexpectedTokenCode,
			line:        &line,
			column:      &column,
		},
	}
}

// UnexpectedToken error
func UnexpectedToken(unexpected *Token, expected ...TokenType) *SyntaxError {
	description := fmt.Sprintf("unexpected token '%s'", unexpected.lexeme)

	if len(expected) > 0 {
		var es []string
		for _, e := range expected[0:] {
			es = append(es, string(e))
		}

		expectation := fmt.Sprintf("'%s'", strings.Join(es, "', '"))
		if len(expected) > 1 {
			expectation = fmt.Sprintf("%s or '%s'", expectation, expected[len(expected)-1])
		}

		description = fmt.Sprintf("%s. Expecting %s", description, expectation)
	}

	return &SyntaxError{
		Error{
			description: description,
			code:        UnexpectedTokenCode,
			line:        &unexpected.line,
			column:      &unexpected.column,
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
	return UnexpectedToken(t, SEMICOLON)
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
	return UnexpectedToken(t, LEFT_BRACE)
}

// ExpectedEndingBrace error
func ExpectedEndingBrace(t *Token) *SyntaxError {
	return UnexpectedToken(t, RIGHT_BRACE)
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

// ReturnStatementOutsideFunction error
func ReturnStatementOutsideFunction(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "return statements must be inside a function or method",
			code:        ReturnStatementOutsideFunctionCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// ArgumentLimitExceeded error
func ArgumentLimitExceeded(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "a function can't have more than 255 arguments",
			code:        ArgumentSizeExceededCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// InvalidTarget raises when an assignment bad targeted
func InvalidTarget(t *Token) *SyntaxError {
	return &SyntaxError{
		Error{
			description: "invalid assignment target",
			code:        InvalidTargetCode,
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
			code:        InvalidDataTypeCode,
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
			code:        InvalidOperationCode,
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
			code:        DivisionByZeroCode,
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
			code:        UndefinedVariableCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// ExpressionIsNotCallable raises when not callable expression is treated as a callable one
func ExpressionIsNotCallable(t *Token) *RuntimeError {
	return &RuntimeError{
		Error{
			description: "expression is not callable",
			code:        ExpressionIsNotCallableCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}

// WrongNumberOfArguments raises when the wrong amount of arguments is passed to a function or method
func WrongNumberOfArguments(t *Token, got, expected int) *RuntimeError {
	return &RuntimeError{
		Error{
			description: fmt.Sprintf("got %v arguments but function expects %v parameters", got, expected),
			code:        ExpressionIsNotCallableCode,
			line:        &t.line,
			column:      &t.column,
		},
	}
}
