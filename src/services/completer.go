package services

import (
	"github.com/c-bata/go-prompt"
)

var resetSuggest []prompt.Suggest

func Completer(in prompt.Document) []prompt.Suggest {

	cmd := Cmd(in.Text)

	switch cmd.Root() {

	case "ps":
		return completerPs(in, cmd.Next())

	case "logs":
		return completerLogs(in, cmd.Next())

	case "exec":
		return completerExec(in, cmd.Next())

	case "consul":
		return completerConsul(in, cmd.Next())

	case "vault":
		return completerVault(in, cmd.Next())

	case "clear":
		return completerClear(in, cmd.Next())

	case "exit":
		return completerExit(in, cmd.Next())

	default:
		return completerDefault(in, cmd.Next())
	}
}

func completerDefault(in prompt.Document, cmd cmd) []prompt.Suggest {
	s := []prompt.Suggest{
		{
			Text:        "ps",
			Description: "List services",
		},
		{
			Text:        "logs",
			Description: "View output from service(s)",
		},
		{
			Text:        "exec",
			Description: "Execute a command in a running service",
		},
		{
			Text:        "consul",
			Description: "Consul commands",
		},
		{
			Text:        "vault",
			Description: "Vault commands",
		},
		{
			Text:        "clear",
			Description: "Stop and remove services",
		},
		{
			Text:        "exit",
			Description: "Exit prompt",
		},
	}

	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func completerPs(in prompt.Document, cmd cmd) []prompt.Suggest {
	s := resetSuggest
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func completerLogs(in prompt.Document, cmd cmd) []prompt.Suggest {
	s := suggest(cmd)
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func completerExec(in prompt.Document, cmd cmd) []prompt.Suggest {
	if !cmd.HasCmd() {
		s := resetSuggest
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
	}
	s := suggest(cmd)
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func completerConsul(in prompt.Document, cmd cmd) []prompt.Suggest {
	switch cmd.Root() {

	case "ui":
		s := suggestStartWith(cmd.Next(), "consul")
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)

	default:
		s := []prompt.Suggest{
			{
				Text:        "ui",
				Description: "Open Web UI",
			},
		}
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)

	}
}

func completerVault(in prompt.Document, cmd cmd) []prompt.Suggest {
	switch cmd.Root() {

	case "ui":
		s := suggestStartWith(cmd.Next(), "vault")
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)

	default:
		s := []prompt.Suggest{
			{
				Text:        "ui",
				Description: "Open Web UI",
			},
		}
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)

	}
}

func completerClear(in prompt.Document, cmd cmd) []prompt.Suggest {
	s := resetSuggest
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func completerExit(in prompt.Document, cmd cmd) []prompt.Suggest {
	s := resetSuggest
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}
