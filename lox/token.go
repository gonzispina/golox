package lox

import "fmt"

type TokenType string

const (
	// Single-character tokens.

	LEFT_PAREN  TokenType = "("
	RIGHT_PAREN TokenType = ")"
	LEFT_BRACE  TokenType = "{"
	RIGHT_BRACE TokenType = "}"
	COMMA       TokenType = ","
	DOT         TokenType = "."
	MINUS       TokenType = "-"
	PLUS        TokenType = "+"
	SEMICOLON   TokenType = ";"
	SLASH       TokenType = "/"
	STAR        TokenType = "*"

	// One or two character tokens.

	BANG          TokenType = "!"
	BANG_EQUAL    TokenType = "!="
	EQUAL         TokenType = "="
	EQUAL_EQUAL   TokenType = "=="
	GREATER       TokenType = ">"
	GREATER_EQUAL TokenType = ">="
	LESS          TokenType = "<"
	LESS_EQUAL    TokenType = "<="

	// Literals.

	IDENTIFIER TokenType = "identifier"
	STRING     TokenType = "string"
	NUMBER     TokenType = "number"

	// Keywords.

	AND    TokenType = "&&"
	OR     TokenType = "||"
	CLASS  TokenType = "class"
	IF     TokenType = "if"
	ELSE   TokenType = "else"
	FALSE  TokenType = "false"
	TRUE   TokenType = "true"
	FUN    TokenType = "fun"
	WHILE  TokenType = "while"
	FOR    TokenType = "for"
	PRINT  TokenType = "print"
	RETURN TokenType = "return"
	SUPER  TokenType = "super"
	THIS   TokenType = "this"
	NIL    TokenType = "nil"
	VAR    TokenType = "var"

	EOF TokenType = "eof"
)

var Reserved = map[string]TokenType{
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

// NewToken constructor
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line, column int) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
		column:    column,
	}
}

// Token representation
type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
	column    int
}

// String cast
func (t *Token) String() string {
	return fmt.Sprintf("[Line: %v] %s %s %v", t.line, t.tokenType, t.lexeme, t.literal)
}

// Is the of the type passed
func (t *Token) Is(tt TokenType) bool {
	return t.tokenType == tt
}

// OneOf the types passed
func (t *Token) OneOf(tts ...TokenType) bool {
	for _, tt := range tts {
		if t.Is(tt) {
			return true
		}
	}
	return false
}
