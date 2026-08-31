package engine

import (
	"fmt"

	. "github.com/leandrolopesm/scheme-go/core"
)

func (self *Engine) RegisterBuiltins() {
	self.AddFunc("display", []TypeFilter{Any}, false, display)
	self.AddFunc("newline", []TypeFilter{}, false, newline)

	self.AddFunc("boolean?", []TypeFilter{Any}, false, isX(Bool))
	self.AddFunc("integer?", []TypeFilter{Any}, false, isX(Integer))
	self.AddFunc("rational?", []TypeFilter{Any}, false, isX(Float))

	self.AddFunc("+", []TypeFilter{Any}, true, VarOp('+'))
	self.AddFunc("-", []TypeFilter{Any}, true, VarOp('-'))
	self.AddFunc("*", []TypeFilter{Any}, true, VarOp('*'))
	self.AddFunc("/", []TypeFilter{Any}, true, VarOp('/'))

	self.AddFunc("eqv", []TypeFilter{Any, Any}, true, eqv)
}

func VarOp(kind rune) BuiltinExec {
	floatOp := func (a float64, b float64) float64 {
		switch kind {
			case '+': return a + b
			case '-': return a - b
			case '*': return a * b
			case '/': return a / b
		}

		panic(fmt.Sprintf("Undefined operation %c", kind))
	}
	
	intOp := func (a int64, b int64) int64 {
		switch kind {
			case '+': return a + b
			case '-': return a - b
			case '*': return a * b
			case '/': return a / b
		}

		panic(fmt.Sprintf("Undefined operation %c", kind))
	}

	return func (e *Engine) error {
		filter := NumberFt
		var numbers []Unit
		var outType Type = Integer // We can be optimistic, right?

		v, err := e.Pop();

		for err == nil {
			if !filter.Matches(v.Type) {
				return fmt.Errorf("Expected integer or float, got %s", TypeNames[v.Type])
			}

			if v.Type == Float {
				outType = Float
			}

			numbers = append(numbers, v)
			
			v, err = e.Pop()
		}

		switch outType {
		case Float:
			var out float64 = 0.
			for _,num := range numbers {
				switch num.Type {
				case Float:
					out = floatOp(out, num.Value.(float64))
				default:
					out = floatOp(out, float64(num.Value.(int64)))
				}
			}

			e.Push(MkFloat(out))
		default:
			var out int64 = 0.
			for _,num := range numbers {
				out = intOp(out, num.Value.(int64))
			}

			e.Push(MkInt(out))
		}

		return nil
	}
}

func isX(which Type) BuiltinExec {
	return func(e *Engine) error {
		if val, err := e.Pop(); e != nil {
			return err
		} else {
			e.Push(MkBool(val.Type == which))
		}

		return nil
	}
}

func newline(e *Engine) error {
	print("\n")

	return nil
}

func eqv(e *Engine) error {
	lhs, lErr := e.Pop()
	rhs, rErr := e.Pop()

	switch {
	case lErr != nil:
		return lErr
	case rErr != nil:
		return rErr
	}

	if lhs.Type != rhs.Type {
		e.Push(MkBool(false))
	}

	return nil
}

func display(e *Engine) error {
	if val, err := e.Pop(); err != nil {
		return err
	} else {
		PrintUnit(val)
	}

	return nil
}