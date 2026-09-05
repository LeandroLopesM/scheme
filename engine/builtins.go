package engine

import (
	"errors"
	"fmt"
	"math"

	. "github.com/leandrolopesm/scheme/core"
	"github.com/leandrolopesm/scheme/util"
)

func (self *Engine) RegisterBuiltins() {
	self.AddFunc("display", []TypeFilter{Any}, true, display)
	self.AddFunc("newline", []TypeFilter{}, false, newline)

	self.AddFunc("string?", []TypeFilter{Any}, false, isX(String))
	self.AddFunc("symbol?", []TypeFilter{Any}, false, isX(Symbol))
	self.AddFunc("boolean?", []TypeFilter{Any}, false, isX(Bool))
	self.AddFunc("integer?", []TypeFilter{Any}, false, isX(Integer))
	self.AddFunc("rational?", []TypeFilter{Any}, false, isX(Float))

	self.AddFunc("+", []TypeFilter{NumberFt}, true, MathOp('+'))
	self.AddFunc("-", []TypeFilter{NumberFt}, true, MathOp('-'))
	self.AddFunc("*", []TypeFilter{NumberFt}, true, MathOp('*'))
	self.AddFunc("/", []TypeFilter{NumberFt}, true, MathOp('/'))

	self.AddFunc("max", []TypeFilter{NumberFt}, true, OrdOp('>'))
	self.AddFunc("min", []TypeFilter{NumberFt}, true, OrdOp('<'))

	self.AddFunc("expt", []TypeFilter{NumberFt}, true, expt)

	self.AddFunc("eqv", []TypeFilter{Any}, true, eqv)
	
	self.AddFunc("define", []TypeFilter{SymbolFt, Any}, false, define)
	
	self.AddFunc("string", []TypeFilter{CharFt}, true, stringize)
	self.AddFunc("string-ref", []TypeFilter{StringFt, IntegerFt}, false, stringRef)
	self.AddFunc("string-append", []TypeFilter{StringFt}, true, stringConcat)
	self.AddFunc("make-string", []TypeFilter{IntegerFt}, false, stringCreate) // This shouldnt get used much
	self.AddFunc("string-set!", []TypeFilter{StrVarFt, IntegerFt, CharFt}, false, stringSet) // This shouldnt get used much
}

func stringSet(e *Engine) error {
	char := util.Assert(e.Pop()).Value.(rune)
	idx := util.Assert(e.Pop()).Value.(int64)
	strVar := util.Assert(e.Pop())

	var trueStrVal string

	if strVar.Type == Symbol {
		if val, ok := e.vars[strVar.Value.(string)]; !ok {
			return fmt.Errorf("Undefined variable %s", strVar.Value.(string))
		} else {
			trueStrVal = val.Value.(string)
		}
	} else {
		trueStrVal = strVar.Value.(string)
	}

	asArr := []rune(trueStrVal);
	if int(idx) > len(asArr) || idx < 0 {
		return fmt.Errorf("Index %d out of bounds for %d ", idx, len(asArr))
	}

	asArr[idx] = char;

	if strVar.Type == Symbol {
		e.vars[strVar.Value.(string)] = MkString(string(asArr))
	} else {
		e.Push(MkString(string(asArr)))
	}

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
		if val,err := e.Pop(); err != nil{
			break
		} else {
			strs = append(strs, val.Value.(string))
		}
	}

	var out string
	idx := len(strs) - 1;

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
		if val,err := e.Pop(); err != nil{
			break
		} else {
			chars = append(chars, val.Value.(rune))
		}
	}

	var out []rune
	idx := len(chars) - 1;

	for idx >= 0 {
		out = append(out, chars[idx])

		idx--
	}

	e.Push(MkString(string(out)))
	return nil
}

func define(e *Engine) error {
	if value, err := e.Pop(); err != nil {
		return errors.New("Expected variable value")
	} else {
		if name, err := e.Pop(); err != nil {
			return errors.New("Expected variable name")
		} else {
			e.vars[name.Value.(string)] = value
		}
	}

	return nil
}

func OrdOp(kind rune) BuiltinExec {
	return func(e *Engine) error {
		var nums []Unit
		var overallType Type

		val, err := e.Pop()
		for err == nil {
			if len(nums) == 0 {
				overallType = val.Type
			} else if val.Type != overallType {
				return fmt.Errorf(
					"%s expects all arguments to be the same type (Was %s, now %s)",
					util.If(kind == '>',
						"(max)",
						"(min)",
					),
					TypeNames[overallType],
					TypeNames[val.Type],
				)
			}

			nums = append(nums, val)
			val, err = e.Pop()
		}

		switch overallType {
		case Float:
			var curr = numAsF(nums[0])
			for _, v := range nums {
				if util.If(kind == '>',
					numAsF(v) > curr,
					numAsF(v) < curr,
				) {
					curr = numAsF(v)
				}
			}

			e.Push(MkFloat(curr))
		default:
			var curr = nums[0].Value.(int64)
			for _, v := range nums {

				if util.If(
					kind == '>',
					v.Value.(int64) > curr,
					v.Value.(int64) < curr,
				) {
					curr = v.Value.(int64)
				}
			}

			e.Push(MkInt(curr))
		}

		return nil
	}
}

func numAsF(num Unit) float64 {
	switch num.Type {
	case Float:
		return num.Value.(float64)
	default:
		return float64(num.Value.(int64))
	}
}

func expt(e *Engine) error {
	lhs, lErr := e.Pop()
	rhs, rErr := e.Pop()

	if lErr != nil {
		return lErr
	} else if rErr != nil {
		return rErr
	}

	var lFloat float64 = numAsF(lhs)
	var rFloat float64 = numAsF(rhs)

	e.Push(MkFloat(math.Pow(lFloat, rFloat)))

	return nil
}

// TODO: (- 4) => -4
// TODO: (/ 4) => 1/4
func MathOp(kind rune) BuiltinExec {
	floatOp := func(a float64, b float64) float64 {
		switch kind {
		case '+':
			return a + b
		case '-':
			return a - b
		case '*':
			return a * b
		case '/':
			return a / b
		}

		panic(fmt.Sprintf("Undefined operation %c", kind))
	}

	intOp := func(a int64, b int64) int64 {
		switch kind {
		case '+':
			return a + b
		case '-':
			return a - b
		case '*':
			return a * b
		case '/':
			return a / b
		}

		panic(fmt.Sprintf("Undefined operation %c", kind))
	}

	return func(e *Engine) error {
		filter := NumberFt
		var numbers []Unit
		var outType Type = Integer // We can be optimistic, right?

		v, err := e.Pop()

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
			var out float64 = numAsF(numbers[0])
			for _, num := range numbers[1:] {
				out = floatOp(out, numAsF(num))
			}

			e.Push(MkFloat(out))
		default:
			var out int64 = numbers[0].Value.(int64)
			for _, num := range numbers[1:] {
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
	lhs := util.Assert(e.Pop()) // We can assert because the argCount was checked
	rhs := util.Assert(e.Pop())

	if lhs.Type != rhs.Type {
		e.Push(MkBool(false))
	}

	e.Push(MkBool(lhs == rhs))
	return nil
}

func display(e *Engine) error {
	var args []Unit

	for {
		if val, err := e.Pop(); err != nil {
			break
		} else {
			args = append(args, val)
		}
	}

	idx := len(args) - 1;

	for idx >= 0 {
		print(SprintUnit(args[idx]))

		idx--
	}

	return nil
}
