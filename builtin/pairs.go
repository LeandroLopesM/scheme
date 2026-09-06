package builtin

import (
	. "github.com/leandrolopesm/scheme/core"
	. "github.com/leandrolopesm/scheme/engine"
	"github.com/leandrolopesm/scheme/util"
)

func pairGet(idx int) BuiltinExec {
	return func (e *Engine) error {
		pair := util.Assert(e.Pop()).Value.(PairVal)

		e.Push(pair[idx])
		return nil
	}
}

func pairSet(idx int) BuiltinExec {
	return func (e *Engine) error {
		pairVarName := util.Assert(e.Pop()).Value.(string)
		pairVarVal := util.Assert(e.GetVar(pairVarName)).Value.(PairVal)
		
		newVal := util.Assert(e.Pop())
		pairVarVal[idx] = newVal

		e.SetVar(pairVarName, MkPair(pairVarVal))
		return nil
	}
}

func newPair(e *Engine) error {
	cdr := util.Assert(e.Pop())
	car := util.Assert(e.Pop())

	e.Push(MkPair(PairVal{ car, cdr }))

	return nil
}