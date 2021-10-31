package lox

import (
	"fmt"
)

// NewInterpreter constructor
func NewInterpreter() *Interpreter {
	e := NewEnvironment(nil)
	e.define("clock", NewClockFunction())
	return &Interpreter{environment: e, locals: map[Expression]int{}}
}

// Interpreter of the lox language
type Interpreter struct {
	environment *Environment
	locals      map[Expression]int
}

// Interpret the given expression
func (i *Interpreter) Interpret(s []Stmt) error {
	for _, stmt := range s {
		v, err := i.execute(stmt)
		if err != nil {
			return err
		}
		if v != nil {
			if _, ok := v.(Callable); !ok {
				fmt.Printf("%v\n", v)
			}
		}
	}
	return nil
}

func (i *Interpreter) Resolve(e Expression, distance int) {
	i.locals[e] = distance
}

func (i *Interpreter) lookUpVariable(lexeme string, variable *Variable) (interface{}, error) {
	var v interface{}
	var found bool

	distance, ok := i.locals[variable]
	if ok {
		v, found = i.environment.getAt(lexeme, distance)
	} else {
		v, found = i.environment.get(lexeme)
	}

	if !found {
		return nil, UndefinedVariable(lexeme, variable.token)
	}

	return v, nil

}

func (i *Interpreter) stringify(v interface{}) string {
	dt := getDataType(v)
	if dt == object {
		if v == nil {
			return "nil"
		}
	} else if dt == boolean {
		if v.(bool) {
			return "true"
		} else {
			return "false"
		}
	}
	return fmt.Sprintf("%v", v)
}

func (i *Interpreter) execute(s Stmt) (interface{}, error) {
	return s.Accept(i)
}

func (i *Interpreter) evaluate(e Expression) (interface{}, error) {
	return e.Accept(i)
}

func (i *Interpreter) visitGrouping(e *Grouping) (interface{}, error) {
	return i.evaluate(e.expression)
}

func (i *Interpreter) visitLiteral(e *Literal) (interface{}, error) {
	return e.value, nil
}

func (i *Interpreter) visitUnary(e *Unary) (interface{}, error) {
	right, err := i.evaluate(e.right)
	if err != nil {
		return nil, err
	}

	switch e.operator.tokenType {
	case MINUS:
		v, ok := right.(float64)
		if !ok {
			return nil, InvalidDataTypeError(e.operator, getDataType(right), number)
		}
		return -v, nil
	case BANG:
		return isTruthy(right), nil
	default:
		return nil, nil
	}
}

func (i *Interpreter) visitBinary(e *Binary) (interface{}, error) {
	left, err := i.evaluate(e.left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(e.right)
	if err != nil {
		return nil, err
	}

	switch e.operator.tokenType {
	case GREATER:
		return greaterThan(left, right, e.operator)
	case GREATER_EQUAL:
		return greaterEqual(left, right, e.operator)
	case LESS:
		return lesserThan(left, right, e.operator)
	case LESS_EQUAL:
		return lesserEqual(left, right, e.operator)
	case EQUAL_EQUAL:
		return isEqual(left, right, e.operator)
	case BANG_EQUAL:
		return notEqual(left, right, e.operator)
	case MINUS:
		v1, v2, err := both2Float(left, right, e.operator)
		if err != nil {
			return nil, err
		}

		return v1 - v2, nil
	case STAR:
		v1, v2, err := both2Float(left, right, e.operator)
		if err != nil {
			return nil, err
		}

		return v1 * v2, nil
	case SLASH:
		v1, v2, err := both2Float(left, right, e.operator)
		if err != nil {
			return nil, err
		}

		if v2 == 0 {
			return nil, DivisionByZeroError(e.operator)
		}

		return v1 / v2, nil
	case PLUS:
		return addValues(left, right, e.operator)
	}

	return nil, nil
}

func (i *Interpreter) visitVariable(e *Variable) (interface{}, error) {
	return i.lookUpVariable(e.token.lexeme, e)
}

// assignment → IDENTIFIER "=" assignment | equality ;
func (i *Interpreter) visitAssign(e *Assign) (interface{}, error) {
	value, err := i.evaluate(e.value)
	if err != nil {
		return nil, err
	}

	ok := i.environment.assign(e.name.lexeme, value)
	if !ok {
		return nil, UndefinedVariable(e.name.lexeme, e.name)
	}

	return nil, nil
}

func (i *Interpreter) visitLogical(e *Logical) (interface{}, error) {
	v, err := i.evaluate(e.left)
	if err != nil {
		return nil, err
	}

	if e.operator.Is(OR) {
		if isTruthy(v) {
			return v, nil
		}
	} else {
		if !isTruthy(v) {
			return v, nil
		}
	}

	return i.evaluate(e.right)
}

func (i *Interpreter) visitCall(e *Call) (interface{}, error) {
	callee, err := i.evaluate(e.callee)
	if err != nil {
		return nil, err
	}

	var arguments []interface{}
	for _, argument := range e.arguments {
		v, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, v)
	}

	c, ok := callee.(Callable)
	if !ok {
		return nil, ExpressionIsNotCallable(e.paren)
	}

	return c.Call(i, e.paren, arguments)
}

func (i *Interpreter) visitGet(e *Get) (interface{}, error) {
	o, err := i.evaluate(e.object)
	if err != nil {
		return nil, err
	}

	if instance, ok := o.(*Instance); ok {
		return instance.Get(e.name)
	}

	return nil, NotAnObject(e.name)
}

func (i *Interpreter) visitSet(e *Set) (interface{}, error) {
	o, err := i.evaluate(e.object)
	if err != nil {
		return nil, err
	}

	instance, ok := o.(*Instance)
	if ok {
		return nil, NotAnObject(e.name)
	}

	v, err := i.evaluate(e.value)
	if err != nil {
		return nil, err
	}

	instance.Set(e.name, v)
	return nil, nil
}

func (i *Interpreter) visitPrintStmt(s *PrintStmt) (interface{}, error) {
	value, err := i.evaluate(s.expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(i.stringify(value))
	return nil, nil
}

func (i *Interpreter) visitExpressionStmt(s *ExpressionStmt) (interface{}, error) {
	v, err := i.evaluate(s.expression)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (i *Interpreter) visitVarStmt(e *VarStmt) (interface{}, error) {
	var value interface{}
	var err error

	if e.initializer != nil {
		value, err = i.execute(e.initializer)
		if err != nil {
			return nil, err
		}
	}

	i.environment.define(e.name.lexeme, value)
	return nil, nil
}

func (i *Interpreter) visitBlockStmt(e *BlockStmt) (interface{}, error) {
	prev := *i.environment
	defer func() {
		i.environment = &prev
	}()

	i.environment = NewEnvironment(i.environment)
	for _, statement := range e.statements {
		v, err := i.execute(statement)
		if err != nil {
			return nil, err
		}

		if v != nil {
			if _, ok := statement.(*CircuitBreakStmt); ok {
				return v, nil
			}
		}
	}

	return nil, nil
}

func (i *Interpreter) visitIfStmt(e *IfStmt) (interface{}, error) {
	v, err := i.evaluate(e.expression)
	if err != nil {
		return nil, err
	}

	if isTruthy(v) {
		return i.execute(e.thenBranch)
	}

	if e.elseBranch != nil {
		return i.execute(e.elseBranch)
	}

	return nil, nil
}

func (i *Interpreter) visitForStmt(e *ForStmt) (interface{}, error) {
	prev := *i.environment
	defer func() {
		*e.cont = false
		*e.br = false
		i.environment = &prev
	}()

	i.environment = NewEnvironment(i.environment)
	if e.initializer != nil {
		_, err := i.execute(e.initializer)
		if err != nil {
			return nil, err
		}
	}

	for {
		if e.condition != nil {
			v, err := i.evaluate(e.condition)
			if err != nil {
				return nil, err
			}

			if !isTruthy(v) {
				break
			}
		}

		for _, stmt := range e.body.statements {
			_, err := i.execute(stmt)
			if err != nil {
				return nil, err
			}

			if *e.cont {
				*e.cont = false
				break
			}

			if *e.br {
				return nil, nil
			}
		}

		if e.increment != nil {
			_, err := i.evaluate(e.increment)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (i *Interpreter) visitFunctionStmt(e *FunctionStmt) (interface{}, error) {
	f := NewFunction(e, i.environment)
	if e.name != nil {
		i.environment.define(e.name.lexeme, f)
	}
	return f, nil
}

func (i *Interpreter) visitCircuitBreakStmt(e *CircuitBreakStmt) (interface{}, error) {
	*e.value = true
	if e.statement != nil {
		return i.execute(e.statement)
	}
	return nil, nil
}

func (i *Interpreter) visitClassStmt(e *ClassStmt) (interface{}, error) {
	i.environment.define(e.name.lexeme, nil)

	methods := map[string]*Function{}
	for _, method := range e.methods {
		methods[method.name.lexeme] = NewFunction(method, i.environment)
	}

	c := NewClass(e, methods, i.environment)
	i.environment.assign(e.name.lexeme, c)
	return nil, nil
}
