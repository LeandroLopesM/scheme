package core

import (
	"fmt"
)

type Scheme struct {
	Name string
	Args []Unit

	Position Position
}

type VectorVal = []Unit

func (pos Position) ToString() string {
	return fmt.Sprintf("%s:%d:%d", pos.File, pos.Line, pos.Char)
}

var TypeNames = map[Type](string){
	Symbol:   "Symbol",
	Integer: "Integer",
	Float:   "Float",
	Bool:    "Bool",
	String:  "String",
	Char:  "Char",
	Vector:  "Vector",
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
	Vector

	// Misc filters for arguments, not actual unit values
	None
	Any
	Number
	StrVar
)

func (tf Type) Matches(ty Type) bool {
	switch tf {
	case Any:
		return true
	case Number:
		return ty == Integer || ty == Float
	case StrVar:
		return ty == String || ty == Symbol
	default:
		return ty == Type(tf)
	}
}

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
func MkVector(v []Unit) Unit {
	return Unit {
		Type:  Vector,
		Value: v,
	}
}