package core

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
)

func DebugUnit(gUnit Unit) {
	debugUnit(gUnit, 0)
}

func PrintUnit(unit Unit) {
	switch unit.Type {
	case SchemeType:
		fmt.Printf("Scheme<%s>", unit.Value.(Scheme).Name)
	case Float:
		fmt.Printf("%f", unit.Value.(float64))
	case Integer:
		fmt.Printf("%d", unit.Value.(int64))
	case String:
		fmt.Printf("%s", unit.Value.(string))
	case Bool:
		fmt.Printf("%v", unit.Value.(bool))
	default:
		panic(fmt.Sprintf("Unreachable type %v", unit))
	}
}

func debugUnit(unit Unit, depth int) {
	switch unit.Type {
	case SchemeType:
		asGroup := unit.Value.(Scheme)
		log.Debugf("%sScheme '%s'", repeat(depth), asGroup.Name)
		for _, member := range asGroup.Args {
			debugUnit(member, depth+1)
		}
	case Float:
		log.Debugf("%sFloat %.1f", repeat(depth), unit.Value.(float64))
	case Integer:
		log.Debugf("%sInt %v", repeat(depth), unit.Value.(int64))
	case String:
		log.Debugf("%sString %s", repeat(depth), unit.Value.(string))
	default:
		log.Warnf("%sUnknown %v", repeat(depth), unit.Value)
	}
}

func repeat(depth int) string {
	return strings.Repeat(". ", depth)
}