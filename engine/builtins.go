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

	self.AddFunc("+", []TypeFilter{Any}, true, add)

	self.AddFunc("eqv", []TypeFilter{Any, Any}, true, eqv)
}

func add(e *Engine) error {
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
		var out float32 = 0.
		for _,num := range numbers {
			switch num.Type {
			case Float:
				out += num.Value.(float32)
			default:
				out += float32(num.Value.(int64))
			}
		}

		e.Push(MkFloat(out))
	default:
		var out int64 = 0.
		for _,num := range numbers {
			out += num.Value.(int64)
		}

		e.Push(MkInt(out))
	}

	return nil
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