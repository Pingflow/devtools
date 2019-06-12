package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/Pingflow/devtools/src/lib"
)

var ErrCommandNotFound = errors.New("command not found")

func Executor(in string) {

	cmd := Cmd(in)

	switch cmd.Root() {

	case "ps":
		executorPs(cmd.Next())
		return

	case "logs":
		executorLogs(cmd.Next())
		return

	case "exec":
		executorExec(cmd.Next())
		return

	case "consul":
		executorConsul(cmd.Next())
		return

	case "vault":
		executorVault(cmd.Next())
		return

	case "clear":
		executorClear(cmd.Next())
		return

	case "exit":
		Stop()
		return

	default:
		newError(ErrCommandNotFound)
		return

	}
}

func executorPs(cmd cmd) {

	if e := dc.Ps(); e != nil {
		newError(e)
		return
	}
}

func executorLogs(cmd cmd) {

	if e := dc.Logs(cmd...); e != nil {
		newError(e)
		return
	}
}

func executorExec(cmd cmd) {

	if e := cmd.HasCmdE(); e != nil {
		newError(e)
		return
	}

	var TTYCmd []string
	if !cmd.Next().HasCmd() {
		TTYCmd = []string{"/bin/sh"}
	} else {
		TTYCmd = cmd.Next()
	}

	fmt.Printf("\n# Connecting to %v...\n", cmd[0])
	if e := dc.TTY(cmd.Root(), TTYCmd...); e != nil {
		newError(e)
		return
	} else {
		fmt.Printf("\n# Disconnected from %v\n", cmd[0])
	}
}

func executorConsul(cmd cmd) {

	if e := cmd.HasCmdE(); e != nil {
		newError(e)
		return
	}

	switch cmd.Root() {

	case "ui":
		executorConsulUi(cmd.Next())
		return

	default:
		executorConsulDefault(cmd.Next())
		return

	}
}

func executorConsulDefault(cmd cmd) {
	newError(ErrCommandNotFound)
}

func executorConsulUi(cmd cmd) {

	if e := cmd.HasCmdE(); e != nil {
		newError(e)
		return
	}

	for _, v := range cmd {
		s, e := dc.Service(v)
		if e != nil {
			newError(e)
			return
		}

		for _, p := range s.ListPorts() {
			if p.Local == 8500 {
				if e := lib.Exec("xdg-open", fmt.Sprintf("http://127.0.0.1:%d", p.Host)); e != nil {
					newError(e)
					return
				}
			}
		}
	}
}

func executorVault(cmd cmd) {

	if e := cmd.HasCmdE(); e != nil {
		newError(e)
		return
	}

	switch cmd.Root() {

	case "ui":
		executorVaultUi(cmd.Next())
		return

	default:
		executorVaultDefault(cmd.Next())
		return

	}
}

func executorVaultDefault(cmd cmd) {
	newError(ErrCommandNotFound)
}

func executorVaultUi(cmd cmd) {

	if e := cmd.HasCmdE(); e != nil {
		newError(e)
		return
	}

	for _, v := range cmd {
		s, e := dc.Service(v)
		if e != nil {
			newError(e)
			return
		}

		for _, p := range s.ListPorts() {
			if p.Local == 8200 {
				if e := lib.Exec("xdg-open", fmt.Sprintf("http://127.0.0.1:%d", p.Host)); e != nil {
					newError(e)
					return
				}
			}
		}
	}
}

func executorClear(cmd cmd) {
	if e := dc.Stop(); e != nil {
		newError(e)
		return
	}

	if e := dc.Down(); e != nil {
		newError(e)
		return
	}

	if e := ResetPF4(); e != nil {
		newError(e)
		return
	}

	os.Exit(0)
}


func newError(e error) {
	fmt.Printf("\033[31mERROR\033[0m: %v\n", e)
}
