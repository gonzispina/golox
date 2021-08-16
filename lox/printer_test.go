package lox_test

import (
	"golox/lox"
	"testing"
)

func TestASTPrinter_Print(t *testing.T) {
	e := lox.NewBinary(
		lox.NewUnary(
			lox.NewToken(lox.MINUS, "-", nil, 1),
			lox.NewLiteral(123),
		),
		lox.NewToken(lox.STAR, "*", nil, 1),
		lox.NewGrouping(lox.NewLiteral(45.67)),
	)

	res, err := lox.NewASTPrinter().Print(e)
	if err != nil {
		t.FailNow()
	}

	if res != "(* (- 123) (group 45.67))" {
		t.Fail()
	}
}
