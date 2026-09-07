package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	. "github.com/leandrolopesm/eunuch/core"
	"github.com/leandrolopesm/eunuch/parser"
	"github.com/leandrolopesm/eunuch/util"
	"github.com/logrusorgru/aurora/v4"
)

type BuiltinExec func(e *Engine) error

type Builtin struct {
	Args    []Type
	VarArgs bool

	Ret Type

	Call BuiltinExec
}

type Engine struct {
	file string

	stack        Stack[Unit]
	stackHistory Stack[int] // Defines the lower bounds for the current stackPtr

	vars  map[string](Unit)
	funcs map[string](Builtin)
}

func New() Engine {
	ret := Engine{
		file: "#ENGINE",

		stack:        NewStack[Unit](),
		stackHistory: NewStack[int](),

		vars:  make(map[string]Unit),
		funcs: make(map[string]Builtin),
	}
	
	return ret
}

func formatStackTrace(strace error) error {
	calls := strings.Split(strace.Error(), "|")
	var out string

	for i := range calls {
		out = fmt.Sprintf("%s%s%s\n", out, strings.Repeat(". ", i+1), calls[i])
	}

	return errors.New(out)
}

func (self *Engine) ExecuteFile(file string) error {
	if v, e := os.ReadFile(file); e != nil {
		return e
	} else {
		self.file = file
		return self.ExecuteStr(string(v))
	}
}

func (self *Engine) ExecuteStr(code string) error {
	if self.file == "#ENGINE" { // If this wasn't called by ExecuteFile
		self.file = "<anonymous>"
	}

	schemes, err := parser.Lex(self.file, code)

	if err != nil {
		return err
	}

	if log.GetLevel() == log.DebugLevel {
		for _, scheme := range schemes {
			DebugUnit(scheme)
		}
	}

	for _, scheme := range schemes {
		if err := self.checkScheme(scheme.Value.(Scheme)); err != nil {
			return formatStackTrace(err)
		}

		if err := self.runScheme(scheme.Value.(Scheme)); err != nil {
			return formatStackTrace(err)
		}
	}

	return nil
}

func (self *Engine) GetVar(name string) (Unit, error) {
	if v,ok := self.vars[name]; !ok {
		return Unit{}, fmt.Errorf("Undefined variable '%s'", name)
	} else {
		return v, nil
	}
}

func (self *Engine) SetVar(name string, val Unit) {
	self.vars[name] = val;
}

func (self *Engine) checkScheme(scheme Scheme) error {
	var actualFn Builtin
	if fn, ok := self.funcs[scheme.Name]; !ok {
		return fmt.Errorf("Undefined function '%s'", scheme.Name)
	} else {
		actualFn = fn
	}

	if len(scheme.Args) != len(actualFn.Args) && !actualFn.VarArgs {
		return fmt.Errorf("Scheme '%s': Expected %d args, got %d", scheme.Name, len(actualFn.Args), len(scheme.Args))
	}

	idx := 0
	for range actualFn.Args {
		inType := scheme.Args[idx].Type

		if scheme.Args[idx].Type == SchemeType {
			asScheme := scheme.Args[idx].Value.(Scheme)
			if e := self.checkScheme(asScheme); e != nil {
				return e
			}

			// If the function exists, use it's return value as the type
			inType = self.funcs[asScheme.Name].Ret
		
		// If the passed argument is a symbol and we dont want a symbol, get its actual type
		} else if scheme.Args[idx].Type == Symbol && !actualFn.Args[idx].Matches(Symbol) {
			if v,e := self.GetVar(scheme.Args[idx].Value.(string)); e != nil {
				return e
			} else {
				inType = v.Type
			}
		}

		if !actualFn.Args[idx].Matches(inType) {
			return fmt.Errorf(
				"Incorrect argument type for '%s'. Expected '%s' got '%s'",
				scheme.Name,
				TypeNames[actualFn.Args[idx]],
				TypeNames[inType],
			)
		}

		if !actualFn.VarArgs { // We only match against the first arg, repeating
			idx++
		}
	}

	return nil
}

func (self *Engine) runScheme(scheme Scheme) error {
	fn := self.funcs[scheme.Name] // Function must exist (Already checked with checkScheme)
	self.saveStack()

	for i, arg := range scheme.Args {
		switch arg.Type {
		case SchemeType:
			if err := self.runScheme(arg.Value.(Scheme)); err != nil {
				return fmt.Errorf("%s %s|%v", scheme.Position.ToString(), scheme.Name, err)
			}
		case Symbol:
			fnArgIdx := util.If(fn.VarArgs, 0, i)

			if !fn.Args[fnArgIdx].Matches(Symbol) || fn.Args[fnArgIdx] == Any { // Any is shorthand for anything OTHER THAN Symbol
				v,_ := self.GetVar(arg.Value.(string))
				self.Push(v)
			} else {
				self.Push(arg)
			}
		default:
			self.Push(arg)
		}
	}

	ret := fn.Call(self)

	self.loadStack()
	if ret != nil {
		return fmt.Errorf("%s %s: %s", scheme.Position.ToString(), scheme.Name, aurora.Red(ret))
	}

	return nil
}

func (self *Engine) AddFunc(name string, args []Type, isVarArg bool, ret Type, call BuiltinExec) error {
	for k := range self.funcs {
		if k == name {
			return fmt.Errorf("Attemt to redeclare function %s", name)
		}
	}

	self.funcs[name] = Builtin{
		Args:    args,
		VarArgs: isVarArg,
		Call:    call,

		Ret: ret,
	}

	return nil
}

func (self *Engine) Pop() (Unit, error) {
	return self.stack.Pop()
}

func (self *Engine) Push(v Unit) {
	self.stack.Push(v)
}
