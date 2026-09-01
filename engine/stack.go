package engine

import (
	"errors"

	"github.com/leandrolopesm/scheme-go/util"
)

const STACK_SIZE = 64

type Stack[T any] struct {
	raw []T
	ptr int
	size int
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		ptr: 0,
		size: STACK_SIZE,
	}
}

func (self *Stack[T]) Pop() (T, error) {
	if self.ptr-1 < 0 {
		var def T
		return def, errors.New("Stack underflow")
	}

    self.ptr--
	return self.raw[self.ptr],nil
}

func (self *Stack[T]) grow() {
	last := self.raw
	self.raw = make([]T, self.size * 2)
	copy(last, self.raw)
}

func (self *Stack[T]) Push(v T) {
	if self.ptr + 1 >= len(self.raw) {
		self.grow()
	}

	self.ptr++;
    self.raw[self.ptr] = v
}

func (self *Engine) saveStack() {
	self.stackHistory.Push(self.stack)
}

func (self *Engine) loadStack() {
	self.stack = util.Assert(self.stackHistory.Pop())
}