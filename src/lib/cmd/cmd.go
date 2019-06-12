package cmd

import (
	"errors"
	"strings"

	"github.com/Pingflow/devtools/src/lib/slice"
)

type Cmd []string

func New(in string) Cmd {
	return slice.RemoveEmpty(strings.Split(in, " "))
}

func (cmd Cmd) HasCmd() bool {
	return len(cmd) > 0
}

func (cmd Cmd) HasCmdE() error {
	if !cmd.HasCmd() {
		return errors.New("command not found")
	}
	return nil
}

func (cmd Cmd) Root() string {
	if cmd.HasCmd() {
		return cmd[0]
	}
	return ""
}

func (cmd Cmd) Next() Cmd {

	if cmd.HasCmd() {
		return cmd[1:]
	}

	return []string{}
}
