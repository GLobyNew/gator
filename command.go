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

const (
	feedURL = "https://www.wagslane.dev/index.xml"
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

func existInDB(s *state, name string) (bool, error) {
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("command 'login' expects only one argument")
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
	if len(cmd.args) != 1 {
		return errors.New("command 'register' expects only one argument")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return err
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
		return err
	}

	err = s.db.DeleteFeeds(context.Background())

	if err != nil {
		return err
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
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("command 'agg' doesn't expect args")
	}

	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	printFeed(feed)
	return nil

}

func handleAddFeed(s *state, cmd command) error {
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	if len(cmd.args) != 2 {
		return errors.New("command 'addfeed' expect 2 args: <name> <url>")
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return err
	}

	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	printFeed(fetchedFeed)

	return nil

}
