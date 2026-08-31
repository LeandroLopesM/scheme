package engine

import (
	. "github.com/leandrolopesm/scheme-go/core"
)

func (self *Engine) RegisterBuiltins() {
	self.AddFunc("display", []Type{Any}, false, display)
	self.AddFunc("newline", []Type{}, false, newline)

	self.AddFunc("boolean?", []Type{Any, Any}, true, isX(Bool))
	self.AddFunc("integer?", []Type{Any}, false, isX(Integer))
	self.AddFunc("rational?", []Type{Any}, false, isX(Float))

	self.AddFunc("eqv", []Type{Any, Any}, true, eqv)
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