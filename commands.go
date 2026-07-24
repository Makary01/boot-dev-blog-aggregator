package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	handlersByName map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlersByName[cmd.name]
	if !ok {
		return fmt.Errorf("Command '%s' doesn't exist", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	c.handlersByName[name] = handler
}
