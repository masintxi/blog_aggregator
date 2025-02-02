package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	regisCommands map[string]func(*state, command) error
}

func (cmds *commands) register(name string, fun func(*state, command) error) {
	cmds.regisCommands[name] = fun
}

func (cmds *commands) run(s *state, cmd command) error {
	fun, ok := cmds.regisCommands[cmd.name]
	if !ok {
		return fmt.Errorf("the command <%v> was not found", cmd.name)
	}
	return fun(s, cmd)
}
