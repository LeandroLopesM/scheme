package core

type Scheme struct {
	Name string
	Args []Unit
}

var TypeNames = map[Type](string){
	Ident:   "Identifier",
	Integer: "Integer",
	Float:   "Float",
	Bool:    "Bool",
	String:  "String",
	Any:     "Any",
}

type Type int

const (
	SchemeType Type = iota
	Ident
	Integer
	Float
	Bool
	String

	Any // Only used in builtin function definitions
)

type Unit struct {
	Type Type

	Value any
}

func MkNumber(v int64) Unit {
	return Unit{
		Type:  Integer,
		Value: v,
	}
}
func MkFloat(v float32) Unit {
	return Unit{
		Type:  Float,
		Value: v,
	}
}
func MkBool(v bool) Unit {
	return Unit{
		Type:  Bool,
		Value: v,
	}
}
func MkString(v string) Unit {
	return Unit{
		Type:  String,
		Value: v,
	}
}