package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// An Env that mapping var -> value
type Env map[Var]float64

// An Expr interface is an arithmetic expression.
type Expr interface {
	Eval(env Env) float64
	String() string
	Check(vars map[Var]bool) error
}

// A Var identifies a variable, e.g., x.
type Var string

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (v Var) String() string {
	return "变量:" + string(v)
}

// A literal is a numeric constant, e.g., 3.141.
type literal float64

func (l literal) Eval(env Env) float64 {
	return float64(l)
}

func (literal) Check(vars map[Var]bool) error {
	return nil
}

func (l literal) String() string {
	return "常量:" + strconv.FormatFloat(float64(l), 'f', -1, 64)

}

// A unary represents a unary operator expression, e.g., -x.
type unary struct {
	op rune // one of '+' | '-'
	x  Expr
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

func (u unary) String() string {
	return "(操作符号:" + strconv.QuoteRuneToASCII(u.op) + " | " +
		u.x.String() + ")"

}

// A binary represents a binary operator expression, e.g., x+y.
type binary struct {
	op   rune // one of '+', '-', '*', '/'
	x, y Expr
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", b.op))
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (b binary) String() string {
	return "(" + b.x.String() + " | 操作符号:" +
		strconv.QuoteRuneToASCII(b.op) + " | " + b.y.String() + ")"

}

// A call represents a function call expression, e.g., sin(x).
type call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (c call) String() string {
	var args string
	for i, v := range c.args {
		args += v.String()
		if i < len(c.args)-1 {
			args += ", "
		}
	}
	//fmt.Println(args)
	return "函数:" + c.fn + "(" + args + ")"
}

// Test Eval
//func TestEval(t *testing.T) {
//	tests := []struct {
//		expr string
//		env  Env
//		want string
//	}{
//		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
//		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
//		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
//		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
//		{"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
//		{"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
//	}
//	var prevExpr string
//	for _, test := range tests {
//		if test.expr != prevExpr {
//			fmt.Printf("\n%s\n", test.expr)
//			prevExpr = test.expr
//		}
//		expr, err := Parse(test.expr)
//		if err != nil {
//			t.Error(err)
//			continue
//		}
//		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
//		fmt.Printf("\t%v => %s\n", test.env, got)
//		if got != test.want {
//			t.Errorf("%s.Eval() in %v = %q, want %q\n",
//				test.expr, test.env, got, test.want)
//		}
//	}
//}

func main() {
	env := Env{"x": 3, "y": 4}
	xy := unary{
		op: '-',
		x:  Var("y"),
	}

	add := binary{
		op: '+',
		x:  Var("x"),
		y:  Var("y"),
	}

	// 乘
	mul := binary{
		op: '*',
		x:  Var("x"),
		y:  Var("y"),
	}

	// pow
	pow := call{
		fn:   "pow",
		args: []Expr{Var("x"), Var("y")},
	}

	fmt.Println("xy:", add.Eval(env))
	fmt.Println(xy.Eval(env))
	fmt.Println(mul.Eval(env))
	fmt.Println(pow.Eval(env))

	var val Expr = Var("x")
	fmt.Println(val)
	val = literal(234.323233)
	fmt.Println(val)
	fmt.Println(xy)
	fmt.Println(add)
	fmt.Println(mul)
	fmt.Println(pow)

	wawa := Env{"q": 1, "w": 2, "e": 3, "x": 3, "y": 4}
	fmt.Println("变量列表：", wawa)
	taowa := call{
		fn: "pow",
		args: []Expr{
			add,
			call{
				fn: "sqrt",
				args: []Expr{
					binary{
						op: '+',
						x:  Var("e"),
						y:  Var("q"),
					},
				},
			},
		},
	}
	fmt.Println(taowa)            // 函数:pow((变量:x | 操作符号:'+' | 变量:y), 函数:sqrt((变量:e | 操作符号:'+' | 变量:q)))
	fmt.Println(taowa.Eval(wawa)) // 49
}
