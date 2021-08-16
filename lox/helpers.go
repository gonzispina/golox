package lox

import "fmt"

func isDigit(v rune) bool {
	return v >= '0' && v <= '9'
}

func isAlpha(v rune) bool {
	return v >= 'a' && v <= 'z' || v >= 'A' && v <= 'Z' || v == '_'
}

func isAlphanumeric(v rune) bool {
	return isAlpha(v) || isDigit(v)
}

func report(line int, where, message string) {
	fmt.Printf("[line %v] Error %s: %s", line, where, message)
}
