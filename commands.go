package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	handlersMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.handlersMap[cmd.name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.name)
	}
	err := function(s, cmd)
	if err != nil {
		return fmt.Errorf("failed running %s command: %w", cmd.name, err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlersMap[name] = f
}
