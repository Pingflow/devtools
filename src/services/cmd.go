package services

import (
	"errors"
	"github.com/Pingflow/devtools/src/lib"
	"strings"
)

type cmd []string


func Cmd(in string) cmd {
	return lib.RemoveEmptySlice(strings.Split(in, " "))
}

func (cmd cmd) HasCmd() bool {
	return len(cmd) > 0
}

func (cmd cmd) HasCmdE() error {
	if !cmd.HasCmd() {
		return errors.New("command not found")
	}
	return nil
}

func (cmd cmd) Root() string {
	if cmd.HasCmd() {
		return cmd[0]
	}
	return ""
}

func (cmd cmd) Next() cmd {

	if cmd.HasCmd() {
		return cmd[1:]
	}

	return []string{}
}
