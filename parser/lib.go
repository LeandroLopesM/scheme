package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	. "github.com/leandrolopesm/scheme-go/core"
)

type Lexer struct {
	iter Iterator[rune]
	line int
}

func (lex *Lexer) Error(msg string) error {
	return fmt.Errorf("%v:%v: %v", lex.line, lex.iter.Tell(), msg)
}

func Lex(raw string) ([]Unit, error) {
	lex := Lexer {
		iter: NewIterator([]rune(raw)),
		line: 1,
	}

	var out []Unit
	currVal, err := lex.iter.Curr()

	for err == nil {
		switch *currVal {
		case '\n', '\r':
			lex.line++
		case ';':
		for err == nil {
			if *currVal == '\n' {
				lex.iter.Prev()
				break
			}

			currVal, err = lex.iter.Next()
		}
		case '(':
			lex.iter.Consume()

			if group, e := lex.parseGroup(); e != nil {
				return out, e
			} else {
				out = append(out, Unit{
					Type: SchemeType,
					Value: group,
				})
			}
		}

		currVal,err = lex.iter.Next()
	}

	return out, nil
}

func (lex *Lexer) parseGroup() (Scheme, error) {
	var out Scheme;

	out.Name = lex.parseId().Value.(string)

	val, err := lex.iter.Next();

	for err == nil {
		switch {
		case *val == ')': return out, nil 
		case *val == '\n':
			lex.line++
			goto nxt
		case *val == ';':
			for err == nil {
				if *val == '\n' {
					lex.iter.Prev()
					break
				}

				val, err = lex.iter.Next()
			}
		case isWhitespace(*val):
			goto nxt

		case unicode.IsDigit(*val):
			if val, err := lex.parseNum(); err == nil {
				out.Args = append(out.Args, val)
			} else {
				return Scheme{}, err
			}

		case *val == '"':
			lex.iter.Consume()
			if val,err := lex.parseStr(); err == nil {
				out.Args = append(out.Args, val)
			} else {
				return Scheme{},err
			}
		
		case *val == '(': 
			lex.iter.Consume()
			if val,err := lex.parseGroup(); err == nil {
				out.Args = append(out.Args, Unit{ Type: SchemeType, Value: val })
			} else {
				return Scheme{},err
			}

		case *val == '#':
			out.Args = append(out.Args, lex.parseBool())

		case isIdentVal(*val):
			out.Args = append(out.Args, lex.parseId())
		
		default:
			return Scheme{},lex.Error(fmt.Sprintf("Unexpected parseme %c", *val))
		}

		nxt:
		val, err = lex.iter.Next()
	}

	// return out, nil
	return Scheme{}, lex.Error(fmt.Sprintf("Scheme '%s': Expected ')'", out.Name))
}

func (lex *Lexer) parseBool() Unit {
	next := lex.iter.PeekOr(' ');

	if next != 't' && next != 'f' {
		return lex.parseId() // It wasn't a bool
	}

	lex.iter.Consume() // Eat the postfix

	val := next == 't' 
	return Unit{
		Type: Bool,
		Value: val,
	}
}

func (lex *Lexer) parseId() Unit {
	var buffer []rune

	v, e := lex.iter.Curr()
	for e == nil {
		if !isIdentVal(*v) {
			_,_ = lex.iter.Prev();
			break
		}

		
		buffer = append(buffer, *v)
		v, e = lex.iter.Next()
	}

	
	return Unit{
		Value: string(buffer),
		Type: Integer,
	}
}

func (lex *Lexer) parseNum() (Unit, error) {
	var buffer []rune

	v, e := lex.iter.Curr()
	for e == nil {
		if !unicode.IsDigit(*v) && !isHex(*v) && !isPrefix(*v) && *v != '.' {
			_,_ = lex.iter.Prev();
			break
		}

		buffer = append(buffer, *v)
		v, e = lex.iter.Next()
	}

	asStr := string(buffer)

	
	if strings.Contains(asStr, ".") {
		if val, err := strconv.ParseFloat(asStr, 64); err != nil {
			return Unit{}, lex.Error(fmt.Sprintf("Invalid floating point literal '%v'", asStr))
		} else {
			return Unit{
				Value: val,
				Type: Float,
			}, nil
		}
	}

	var radix int

	if startsWith(asStr, "0x") {
		radix = 16
	} else if startsWith(asStr, "0b") {
		radix = 2
	} else {
		radix = 10
	}

	if val, err := strconv.ParseInt(asStr, radix, 64); err != nil {
		return Unit{}, lex.Error(fmt.Sprintf("Invalid integer literal '%v'", asStr))
	} else {
		return Unit{
			Value: val,
			Type: Integer,
		}, nil
	}
}

func (lex *Lexer) parseStr() (Unit, error) {
	var out []rune

	start := lex.iter.Tell();

	v,e := lex.iter.Curr()
	for e == nil {
		if *v == '"' && lex.iter.PeekOr(' ') == '"' { // Strings end in ""
			lex.iter.Consume() // Glomp '"'
			break
		}

		out = append(out, *v)
		v,e = lex.iter.Next()
	}

	if e != nil {
		return Unit{}, lex.Error(fmt.Sprintf("Unclosed string ..\"%s\"..", stringAround(5, lex.iter.elems, start)))
	}

	return Unit {
		Type: String,
		Value: string(out),
	}, nil
}

func startsWith(str string, str2 string) bool {
	var cmp string = str2

	if len(str) < len(str2) {
		cmp = str
	}

	for i := range cmp {
		if str[i] != str2[i] {
			return false
		}
	}

	return true
}

func isWhitespace(c rune) bool {
	return c == '\t' || c == ' ' || c == '\n' || c == '\r'
}

func isHex(c rune) bool {
	return strings.Contains("ABCDEFG", string(c))
}

func isPrefix(c rune) bool {
	return c == 'b' || c == 'x'
}

func isIdentVal(c rune) bool {
	return !(c == '(' || c == ')' || isWhitespace(c))
}
