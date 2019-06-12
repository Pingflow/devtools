package main

import (
	"log"

	"github.com/Pingflow/devtools/src"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if e := src.App(version, commit, date).Run(); e != nil {
		log.Fatal(e)
	}
}
