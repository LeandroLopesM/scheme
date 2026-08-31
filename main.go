package main

import (
	"os"

	log "github.com/charmbracelet/log"
	"github.com/leandrolopesm/scheme-go/engine"
	"github.com/leandrolopesm/scheme-go/util"
)

const VERSION = "0.2.1";

func main() {
	util.Logger()
	
	if len(os.Args) < 2 {
		log.Error("USAGE: scheme-go <file.lsp>")
		return
	}

	file := os.Args[1]
	lispEngine := engine.New()
	if err := lispEngine.ExecuteStr(string(must(os.ReadFile(file)))); err != nil {
		log.Errorf("Execution failed: %s", err)
	}
}

func must[T any](a T, err error) T {
	if err != nil {
		log.Errorf("Infallible operation failed: %s", err)
		os.Exit(1)
	}

	return a
}