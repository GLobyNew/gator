package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/GLobyNew/gator/internal/config"
	"github.com/GLobyNew/gator/internal/database"
	"github.com/google/uuid"
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

func checkOnlyOneArg(cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no args in command")
	} else if len(cmd.args) > 1 {
		return errors.New("too many args, expect one argument")
	}
	return nil
}

func existInDB(s *state, name string) (bool, error) {
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func handlerLogin(s *state, cmd command) error {
	err := checkOnlyOneArg(cmd)
	if err != nil {
		return err
	}

	exist, err := existInDB(s, cmd.args[0])
	if err != nil {
		return err
	}
	if !exist {
		os.Exit(1)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %q has been added in config!\n", s.cfg.CurrentUserName)
	return nil

}

func handlerRegister(s *state, cmd command) error {
	err := checkOnlyOneArg(cmd)
	if err != nil {
		return err
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("User %q has been created!\n", user.Name)

	s.cfg.SetUser(user.Name)
	fmt.Printf("User %q has been added in config!\n", s.cfg.CurrentUserName)

	return nil
}

func handlerReset(s *state, cmd command) error {

	if len(cmd.args) > 0 {
		return errors.New("command 'reset' doesn't expect args")
	}

	err := s.db.DeleteUsers(context.Background())

	if err != nil {
		os.Exit(1)
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {

	if len(cmd.args) > 0 {
		return errors.New("command 'users' doesn't expect args")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	currentUser := s.cfg.CurrentUserName

	for _, user := range users {
		if user == currentUser {
			fmt.Printf("* %v (current)\n", user)
		} else {
			fmt.Printf("* %v\n", user)
		}
	}

	return nil
}
