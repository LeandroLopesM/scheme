package core

import (
	"fmt"
)

type Scheme struct {
	Name string
	Args []Unit

	Position Position
}

func (pos Position) ToString() string {
	return fmt.Sprintf("%s:%d:%d", pos.File, pos.Line, pos.Char)
}

var TypeFilterNames = map[TypeFilter](string){
	SymbolFt:   "Symbolifier",
	IntegerFt: "Integer",
	FloatFt:   "Float",
	BoolFt:    "Bool",
	StringFt:  "String",
	Any:       "Any",
}

var TypeNames = map[Type](string){
	Symbol:   "Symbolifier",
	Integer: "Integer",
	Float:   "Float",
	Bool:    "Bool",
	String:  "String",
	Char:  "Char",
}

type TypeFilter int

const (
	NumberFt TypeFilter = iota

	SymbolFt
	IntegerFt
	FloatFt
	BoolFt
	StringFt
	Any
)

func (tf TypeFilter) Matches(ty Type) bool {
	switch tf {
	case Any:
		return true
	case NumberFt:
		return ty == Integer || ty == Float
	default:
		return ty == Type(tf)
	}
}

type Type int

const (
	SchemeType Type = iota
	Symbol
	Integer
	Float
	Bool
	String
	Char
)

type Position struct {
	Line, Char int
	File       string
}

type Unit struct {
	Type Type

	Value any
}

func MkInt(v int64) Unit {
	return Unit{
		Type:  Integer,
		Value: v,
	}
}
func MkFloat(v float64) Unit {
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
func MkChar(v rune) Unit {
	return Unit{
		Type:  Char,
		Value: v,
	}
}