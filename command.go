package main

import (
	"errors"
	"fmt"

	"github.com/GLobyNew/gator/internal/config"
	"github.com/GLobyNew/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmd map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.cmd[cmd.name]; ok {
		err := f(s, cmd)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("command not found")
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no args in login command, expects username")
	} else if len(cmd.args) > 1 {
		return errors.New("too many args, expect one argument (username)")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %q has been added in config!\n", s.cfg.CurrentUserName)
	return nil

}
