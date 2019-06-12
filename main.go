package main

import (
	"fmt"
	"log"

	"github.com/Pingflow/devtools/src/services"
	"github.com/c-bata/go-prompt"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	fmt.Printf("[DT] GoLang DevTools for GoMicro %v (%v) built at %v\n", version, commit, date)
	if e := services.Start(); e != nil {
		log.Fatal(e)
	}
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer services.Stop()
	p := prompt.New(
		services.Executor,
		services.Completer,
		prompt.OptionPrefix("dt> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionMaxSuggestion(10),
	)
	p.Run()
}
