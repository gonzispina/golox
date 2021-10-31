package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	expressions := map[string]string{
		"Assign":   "name *Token, value Expression",
		"Binary":   "left Expression, operator *Token, right Expression",
		"Call":     "callee Expression, paren *Token, arguments []Expression",
		"Get":      "object Expression, name *Token",
		"Set":      "object Expression, name *Token, value Expression",
		"Grouping": "expression Expression",
		"Logical":  "left Expression, operator *Token, right Expression",
		"Literal":  "value interface{}",
		"This":     "keyword *Token",
		"Unary":    "operator *Token, right Expression",
		"Variable": "token *Token",
	}

	statements := map[string]string{
		"ExpressionStmt":   "expression Expression",
		"FunctionStmt":     "name *Token, params []*Token, body *BlockStmt, rt *bool",
		"IfStmt":           "expression Expression, thenBranch *BlockStmt, elseBranch *BlockStmt",
		"ForStmt":          "initializer Stmt, condition Expression, increment Expression, body *BlockStmt, br *bool, cont *bool",
		"PrintStmt":        "expression Expression",
		"VarStmt":          "name *Token, initializer Stmt",
		"BlockStmt":        "statements []Stmt",
		"ClassStmt":        "name *Token, methods []*FunctionStmt",
		"CircuitBreakStmt": "value *bool, statement Stmt",
	}

	dir, _ := os.Getwd()
	path := filepath.Join(dir, "/lox/ast.go")
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating the file: %s", err.Error())
		os.Exit(64)
	}

	f.WriteString("package lox\n\n")
	define("Expression", expressions, f)
	define("Stmt", statements, f)

	if err := f.Close(); err != nil {
		fmt.Printf("Error closing the file: %s", err.Error())
		os.Exit(64)
	}
}

func define(iface string, types map[string]string, f *os.File) {
	f.WriteString(fmt.Sprintf("// %s representation\n", iface))
	f.WriteString(fmt.Sprintf("type %s interface {\n", iface))
	f.WriteString(fmt.Sprintf("	Accept(v %sVisitor) (interface{}, error)\n", iface))
	f.WriteString("}\n\n")

	f.WriteString(fmt.Sprintf("// %sVisitor defines the visit method of every %s\n", iface, iface))
	f.WriteString(fmt.Sprintf("type %sVisitor interface {\n", iface))
	for name, _ := range types {
		f.WriteString(fmt.Sprintf("	visit%s(e *%s) (interface{}, error)\n", name, name))
	}
	f.WriteString("}\n\n")

	for name, value := range types {
		f.WriteString(fmt.Sprintf("// New%s %s constructor\n", name, iface))
		f.WriteString(fmt.Sprintf("func New%s(%s) *%s {\n", name, value, name))
		f.WriteString(fmt.Sprintf("	return &%s{\n", name))
		for _, line := range strings.Split(value, ", ") {
			prop := strings.Split(line, " ")[0]
			f.WriteString(fmt.Sprintf("		%s: %s,\n", prop, prop))
		}
		f.WriteString("	}\n")
		f.WriteString("}\n\n")

		f.WriteString(fmt.Sprintf("// %s %s implementation\n", name, iface))
		f.WriteString(fmt.Sprintf("type %s struct {\n", name))
		for _, line := range strings.Split(value, ", ") {
			f.WriteString("	" + line + "\n")
		}
		f.WriteString("}\n\n")
		f.WriteString("// Accept method of the visitor pattern it calls the proper visit method\n")
		f.WriteString(fmt.Sprintf("func(e *%s) Accept(v %sVisitor) (interface{}, error) {\n", name, iface))
		f.WriteString(fmt.Sprintf("	return v.visit%s(e)\n", name))
		f.WriteString("}\n\n")
	}
}
