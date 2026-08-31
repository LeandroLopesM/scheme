package engine

import (
	"errors"

	. "github.com/leandrolopesm/scheme-go/core"
)

func (self *Engine) Pop() (Unit, error) {
	if self.stackPtr-1 < 0 {
		return Unit{}, errors.New("Stack underflow")
	}

    self.stackPtr -= 1;
	return self.stack[self.stackPtr + 1],nil
}

func (self *Engine) Push(v Unit) {
    self.stackPtr++;
    self.stack[self.stackPtr] = v;
}