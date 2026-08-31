package engine

import (
	"fmt"

	. "github.com/leandrolopesm/scheme-go/core"
	"github.com/leandrolopesm/scheme-go/parser"
)

type BuiltinExec func(e *Engine) error

type Builtin struct {
	Args []Type
	VarArgs bool

	Ret Type

	Call BuiltinExec
}

type Engine struct {
	file string

	stack    []Unit
	stackPtr int

	vars  map[string](Unit)
	funcs map[string](Builtin)
}

func New() Engine {
	ret := Engine{
		file: "#ENGINE",

		stack: make([]Unit, 128),
        stackPtr: 0,

		vars: make(map[string]Unit),
		funcs: make(map[string]Builtin),
	}

	ret.RegisterBuiltins()

	return ret
}

// func (self *Engine) ExecuteFile(file string) {
//     self.ExecuteStr(os.read)
// }

func (self *Engine) ExecuteStr(code string) error {
    if self.file == "#ENGINE" { // If this wasn't called by ExecuteFile
        self.file = "<anonymous>"
    }

    schemes, err := parser.Lex(code)

	if err != nil {
		return err
	}

	for _,scheme := range schemes {
		DebugUnit(scheme);
	}
	
	for _,scheme := range schemes {
		if err := self.checkScheme(scheme.Value.(Scheme)); err != nil {
			return err
		}

		if err := self.runScheme(scheme.Value.(Scheme)); err != nil {
			return err
		}
	}
	return nil
}

func (self *Engine) checkScheme(scheme Scheme) error {
	var actualFn Builtin
	if fn,ok := self.funcs[scheme.Name]; !ok {
		return fmt.Errorf("Undefined function '%s'", scheme.Name)
	} else {
		actualFn = fn
	}

	if len(scheme.Args) != len(actualFn.Args) && !actualFn.VarArgs {
		return fmt.Errorf("Scheme '%s': Expected %d args, got %d", scheme.Name, actualFn.Args, scheme.Args)
	}

	for idx := range actualFn.Args {
		inType := scheme.Args[idx].Type

		if scheme.Args[idx].Type == SchemeType {
			asScheme := scheme.Args[idx].Value.(Scheme);
			if e := self.checkScheme(asScheme); e != nil {
				return e
			}

			inType = self.funcs[asScheme.Name].Ret
		}

		if inType != actualFn.Args[idx] && actualFn.Args[idx] != Any {
			return fmt.Errorf(
				"Incorrect argument type for '%s'. Expected '%s' got '%s'",
				scheme.Name,
				TypeNames[actualFn.Args[idx]],
				TypeNames[inType],
			)
		}
	}
	
	return nil
}

func (self *Engine) runScheme(scheme Scheme) error {
	for _,arg := range scheme.Args {
		if arg.Type == SchemeType {
			if err := self.runScheme(arg.Value.(Scheme)); err != nil {
				return err
			}
		} else {
			self.Push(arg)
		}
	}

	if fn,ok := self.funcs[scheme.Name]; !ok {
		return fmt.Errorf("Undefined function '%s'", scheme.Name)
	} else {
		return fn.Call(self)
	}
}

func (self *Engine) AddFunc(name string, args []Type, isVarArg bool, call BuiltinExec) error {
	for k := range self.funcs {
		if k == name {
			return fmt.Errorf("Attemt to redeclare function %s", name)
		}
	}

    self.funcs[name] = Builtin{
		Args: args,
		VarArgs: isVarArg,
		Call: call,
	};

    return nil;
}
