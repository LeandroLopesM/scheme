package main

import (
	"context"
	"os"

	log "github.com/charmbracelet/log"
	"github.com/leandrolopesm/scheme-go/core"
	"github.com/leandrolopesm/scheme-go/engine"
	"github.com/leandrolopesm/scheme-go/util"

	"github.com/nyaosorg/go-readline-ny"
)

const VERSION = "0.6.0";

func main() {
	log.SetReportTimestamp(false)
	opt := util.Args()
	log.SetLevel(util.If(*opt.Verbose, log.DebugLevel, log.InfoLevel))

	lispEngine := engine.New()

	if opt.Repl {
		var editor readline.Editor

		for {
			if text,err := editor.ReadLine(context.Background()); err != nil {
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
		for _,file := range opt.Files {
			if err := lispEngine.ExecuteStr(string(must(os.ReadFile(file)))); err != nil {
				log.Errorf("Execution failed: %s", err)
			}
		}
	}
}

func must[T any](a T, err error) T {
	if err != nil {
		log.Errorf("Infallible operation failed: %s", err)
		os.Exit(1)
	}

	return a
}