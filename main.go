package main

import (
	"context"
	"fmt"

	log "github.com/charmbracelet/log"
	"github.com/leandrolopesm/scheme/core"
	"github.com/leandrolopesm/scheme/engine"
	"github.com/leandrolopesm/scheme/util"

	aurora "github.com/logrusorgru/aurora/v4"

	"github.com/nyaosorg/go-readline-ny"
)

const VERSION = "0.11.3"

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
			if err := lispEngine.ExecuteFile(file); err != nil {
				fmt.Print(aurora.Red("Execution failed:"), "\n", err)
			}
		}
	}
}
