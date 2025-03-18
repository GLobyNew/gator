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

func existInDB(s *state, name string) (bool, error) {
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
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

func handlerAddFeed(s *state, cmd command, user database.User) error {

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
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	return nil

}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("command 'feeds' doesn't expect args")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("* %s - %s - %s\n", feed.Name, feed.Url, user.Name)
	}

	return nil

}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("command 'follow' expects only one argument <url>")
	}

	url := cmd.args[0]
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	feedFollows, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("* feed %q <-> %q user\n", feedFollows.FeedName, feedFollows.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	if len(cmd.args) != 0 {
		return errors.New("command 'following' doesn't expect arguments")
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	following, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}

	for _, follow := range following {
		fmt.Printf("* %s\n", follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("command 'unfollow' expects only one argument: <feed url>")
	}

	// Check if feed is exist in db
	feedToDelete, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	// We should check that logged-in user even subscribed to feed, before removing it
	// If exist, remove it, if no - return error

	usersFF, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, uFF := range usersFF {
		if uFF.FeedUrl == feedToDelete.Url {
			err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
				Name: user.Name,
				Url:  feedToDelete.Url,
			})
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("logged-in user %q don't follow %q feed", user.Name, feedToDelete.Url)
}
