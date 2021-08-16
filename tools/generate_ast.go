package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	definitions := map[string]string{
		"Binary":   "left Expression, operator *Token, right Expression",
		"Grouping": "expression Expression",
		"Literal":  "value interface{}",
		"Unary":    "operator *Token, right Expression",
	}

	dir, _ := os.Getwd()
	path := filepath.Join(dir, "/lox/ast.go")
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating the file: %s", err.Error())
		os.Exit(64)
	}

	f.WriteString("package lox\n\n")
	f.WriteString("// Expression representation\n")
	f.WriteString("type Expression interface {\n")
	f.WriteString("	Accept(v ExpressionVisitor) (interface{}, error)\n")
	f.WriteString("}\n\n")

	f.WriteString("// ExpressionVisitor defines the visit method of every Expression\n")
	f.WriteString("type ExpressionVisitor interface {\n")
	for name, _ := range definitions {
		f.WriteString(fmt.Sprintf("	visit%s(e *%s) (interface{}, error)\n", name, name))
	}
	f.WriteString("}\n\n")

	for name, value := range definitions {
		f.WriteString(fmt.Sprintf("// New%s Expression constructor\n", name))
		f.WriteString(fmt.Sprintf("func New%s(%s) *%s {\n", name, value, name))
		f.WriteString(fmt.Sprintf("	return &%s{\n", name))
		for _, line := range strings.Split(value, ", ") {
			prop := strings.Split(line, " ")[0]
			f.WriteString(fmt.Sprintf("		%s: %s,\n", prop, prop))
		}
		f.WriteString("	}\n")
		f.WriteString("}\n\n")

		f.WriteString(fmt.Sprintf("// %s Expression implementation\n", name))
		f.WriteString(fmt.Sprintf("type %s struct {\n", name))
		for _, line := range strings.Split(value, ", ") {
			f.WriteString("	" + line + "\n")
		}
		f.WriteString("}\n\n")
		f.WriteString("// Accept method of the visitor pattern it calls the proper visit method\n")
		f.WriteString(fmt.Sprintf("func(e *%s) Accept(v ExpressionVisitor) (interface{}, error) {\n", name))
		f.WriteString(fmt.Sprintf("	return v.visit%s(e)\n", name))
		f.WriteString("}\n\n")
	}

	if err := f.Close(); err != nil {
		fmt.Printf("Error closing the file: %s", err.Error())
		os.Exit(64)
	}
}
