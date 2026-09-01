package main

import (
	"context"
	"fmt"
	"os"

	log "github.com/charmbracelet/log"
	"github.com/leandrolopesm/scheme-go/core"
	"github.com/leandrolopesm/scheme-go/engine"
	"github.com/leandrolopesm/scheme-go/util"

	aurora "github.com/logrusorgru/aurora/v4"

	"github.com/nyaosorg/go-readline-ny"
)

const VERSION = "0.7.1"

func main() {
	opt := util.Args()
	util.Logger(opt)

	lispEngine := engine.New()

	if opt.Repl {
		var editor readline.Editor

		for {
			if text, err := editor.ReadLine(context.Background()); err != nil {
				if err.Error() == "EOF" {
					break
				}

				log.Errorf("Failed to read input: %s", err)
			} else {
				if err := lispEngine.ExecuteStr(text); err != nil {
					log.Errorf("Execution failed: %s", err)
				}

				if v, e := lispEngine.Pop(); e == nil { // If the last call pushed a value, print it
					core.PrintUnit(v)
				}
			}
		}
	} else {
		for _, file := range opt.Files {
			if err := lispEngine.ExecuteStr(string(util.Assert(os.ReadFile(file)))); err != nil {
				fmt.Print(aurora.Red("Execution failed:"), "\n", err)
			}
		}
	}
}
