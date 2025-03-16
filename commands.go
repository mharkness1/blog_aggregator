package main

import (
	"errors"
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(s *state, cmd command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("too many arguments given, login command requires: login <username>")
	}

	name := cmd.Args[0]

	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("failed to set username: %w", err)
	}

	fmt.Println("Username changed successfully.")

	return nil
}
