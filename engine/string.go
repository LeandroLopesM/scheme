package engine

import (
	"fmt"

	. "github.com/leandrolopesm/scheme/core"
	"github.com/leandrolopesm/scheme/util"
)

func stringSet(e *Engine) error {
	char := util.Assert(e.Pop()).Value.(rune)
	idx := util.Assert(e.Pop()).Value.(int64)
	strVar,err := e.GetVar(util.Assert(e.Pop()).Value.(string))

	if err != nil {
		return err
	}

	asArr := []rune(strVar.Value.(string))
	if int(idx) > len(asArr) || idx < 0 {
		return fmt.Errorf("Index %d out of bounds for %d ", idx, len(asArr))
	}

	asArr[idx] = char

	e.SetVar(strVar.Value.(string), MkString(string(asArr)))

	return nil
}

func stringCreate(e *Engine) error {
	var len int64 = util.Assert(e.Pop()).Value.(int64)

	if len < 0 {
		return fmt.Errorf("Invalid index %d", len)
	}

	var tmp = make([]rune, len)
	e.Push(MkString(string(tmp)))

	return nil
}

func stringConcat(e *Engine) error {
	var strs []string
	for {
		if val, err := e.Pop(); err != nil {
			break
		} else {
			strs = append(strs, val.Value.(string))
		}
	}

	var out string
	idx := len(strs) - 1

	for idx >= 0 {
		out = fmt.Sprintf("%s%s", out, strs[idx])
		idx--
	}

	e.Push(MkString(out))
	return nil
}

func stringRef(e *Engine) error {
	idx := util.Assert(e.Pop()).Value.(int64)
	str := util.Assert(e.Pop()).Value.(string)

	if int(idx) > len([]rune(str)) || idx < 0 {
		return fmt.Errorf("Index %d out of bounds for %d", idx, len([]rune(str)))
	}

	e.Push(MkChar([]rune(str)[idx]))
	return nil
}

func stringize(e *Engine) error {
	var chars []rune
	for {
		if val, err := e.Pop(); err != nil {
			break
		} else {
			chars = append(chars, val.Value.(rune))
		}
	}

	var out []rune
	idx := len(chars) - 1

	for idx >= 0 {
		out = append(out, chars[idx])

		idx--
	}

	e.Push(MkString(string(out)))
	return nil
}