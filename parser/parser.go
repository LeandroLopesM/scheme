package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
	. "github.com/leandrolopesm/eunuch/core"
)

type Lexer struct {
	iter Iterator[rune]
	line, lastLineOff int
	file string
}

func (lex *Lexer) Error(msg string) error {
	return fmt.Errorf("%s:%v:%v: %v", lex.file, lex.line, lex.iter.Tell() - lex.lastLineOff, msg)
}

func Lex(file, raw string) ([]Unit, error) {
	lex := Lexer {
		iter: NewIterator([]rune(raw)),
		line: 1,
		file: file,
	}

	var out []Unit

	for {
		if val, err := lex.parseUnit(); err != nil {
			if err == iterEnd {
				return out, nil
			}

			return out, err
		} else {
			out = append(out, val)
		}
	}
}

var iterEnd = errors.New("Iterator ended")
var schemeEnd = errors.New("Scheme ended (')')")

func (lex *Lexer) parseUnit() (Unit, error) {
start:
	curr, err := lex.iter.Curr();

	if err == nil {
		switch curr {
		case ';':
			lex.skipComment();
			goto start
		case '\n':
			lex.iter.Consume()
			lex.line++
			lex.lastLineOff = lex.iter.Tell()
			goto start;
		case '"':
			return lex.parseString();
		case '(':
			return lex.parseScheme();
		case ')':
			return Unit{}, schemeEnd;
		case '#':
			return lex.parseTag(), nil;
		case '\'':
			return lex.parseQuote();
		default:
			switch {
			case unicode.IsDigit(curr):
				return lex.parseNum()
			case unicode.IsSpace(curr):
				lex.iter.Consume();
				goto start;
			default:
				return lex.parseSymbol(), nil
			}
		}
	}

	return Unit{}, iterEnd
}

func (lex *Lexer) parseString() (Unit, error) {
	val, err := lex.iter.Next();
	var buffer []rune
	for err == nil {
		if val == '"' && lex.iter.PrevOr(' ') != '\\' {
			lex.iter.Consume();
			lex.iter.Consume();

			return Unit {
				Type: String,
				Value: string(buffer),
			}, nil
		} 
		buffer = append(buffer, val)

		val, err = lex.iter.Next();
	}

	return Unit{}, lex.Error("Unclosed string");
}

func (lex *Lexer) parseQuote() (Unit, error) {
	lex.iter.Consume() // eat '
	
	if val, err := lex.parseUnit(); err != nil {
		return Unit{}, err
	} else {
		return Unit {
			Type: SchemeType,
			Value: Scheme {
				Name: "quote",
				Args: []Unit {
					val,
				},
			},
		}, nil
	}
}

func (lex *Lexer) parseScheme() (Unit, error) {
	lex.iter.Consume() // Discard (
	out := Scheme {
		Name: lex.parseSymbol().Value.(string),
		Position: Position{
			Line: lex.line,
			Char: lex.iter.Tell() - lex.lastLineOff,
			File: lex.file,
		},
	}

	_, err := lex.iter.Curr()
	
	var val Unit
	for err == nil {
		if val, err = lex.parseUnit(); err != nil {
			if err == schemeEnd {
				lex.iter.Consume()
				
				return Unit{
					Type: SchemeType,
					Value: out,
				}, nil
			}

			return Unit{}, err
		} else {
			out.Args = append(out.Args, val)
		}

		// _, err = lex.iter.Next();
	}

	return Unit{}, lex.Error(fmt.Sprintf("Scheme '%s': Missing closing ')'", out.Name))
}

func (lex *Lexer) parseTag() Unit {
	asIdent := lex.parseSymbol(); // #...
	asStr := asIdent.Value.(string); // #...

	if len(asStr) > 1 {
		switch []rune(asStr)[1] {
		case '\\':
			return lex.parseChar(asStr);
		case 't', 'f':
			return lex.parseBool(asStr[1:]);
		}
	}

	return asIdent
}

func (lex *Lexer) parseBool(str string) Unit {
	return Unit{
		Type: Bool,
		Value: str == "t",
	}
}

func (lex *Lexer) parseChar(str string) Unit {
	var char rune
	switch str[2:] {
	case "newline":
		char = '\n'
	case "space":
		char = ' '
	default:
		if len(str[2:]) > 1 {
			log.Warnf("Unknown character expr '%s', using '%c'", str, str[2])
		}

		char = []rune(str)[2]
	}

	return Unit{
		Type: Char,
		Value: char,
	}
}

func (lex *Lexer) parseNum() (Unit, error) {
	num := lex.parseSymbol().Value.(string)

	if strings.Contains(num, ".") {
		if val, err := strconv.ParseFloat(num, 64); err != nil {
			return Unit{}, lex.Error(fmt.Sprintf("Invalid float number '%s'", num))
		} else {
			return Unit {
				Type: Float,
				Value: val,
			}, nil
		}
	}

	radix := 10
	if strings.ContainsAny(num, "ABCDEFabcdef") {
		radix = 16
	}

	if val, err := strconv.ParseInt(num,radix, 64); err != nil {
		return Unit{}, lex.Error(fmt.Sprintf("Invalid integer '%s'", num))
	} else {
		return Unit {
			Type: Integer,
			Value: val,
		}, nil
	}
}

func (lex *Lexer) parseSymbol() Unit {
	var buff []rune
	nxt, err := lex.iter.Curr()
	
	for err == nil {
		buff = append(buff, nxt)

		if nxt,err = lex.iter.Next(); err != nil {
			lex.iter.Consume()
			break
		} else if !isIdentChar(nxt) {
			// _,_ = lex.iter.Prev()
			break
		}
	}

	return Unit {
		Type: Symbol,
		Value: string(buff),
	}
}

func isIdentChar(char rune) bool {
	return !(unicode.IsSpace(char) || char == ')' || char == '(')
}

func (lex *Lexer) skipComment() {
	curr, err := lex.iter.Curr()
	var buff []rune
	for err == nil {
		if curr, err = lex.iter.Next(); err != nil {
			return;
		} else {
			if curr == '\n' {
				return;
			}

			buff = append(buff, curr)
		}
	}
}