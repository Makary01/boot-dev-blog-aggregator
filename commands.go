package main

import (
	"Makary01/boot-dev-blog-aggregator/internal/database"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

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

func handlerLogin(s *state, cmd command) error {
	if err := checkArgsLen(cmd.args, 1); err != nil {
		return err
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if err := checkArgsLen(cmd.args, 1); err != nil {
		return err
	}
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("User created: %v\n", user)

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := checkArgsLen(cmd.args, 0); err != nil {
		return err
	}
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if err := checkArgsLen(cmd.args, 0); err != nil {
		return err
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range users {
		line := u.Name
		if u.Name == s.config.CurrentUserName {
			line += " (current)"
		}
		fmt.Println(line)
	}
	return nil
}

func handlerAgg(_ *state, cmd command) error {
	if err := checkArgsLen(cmd.args, 0); err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if err := checkArgsLen(cmd.args, 2); err != nil {
		return err
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if err := checkArgsLen(cmd.args, 0); err != nil {
		return err
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, f := range feeds {
		fmt.Printf("`%s` - %s, created by: %s\n", f.Name, f.Url, f.CreatedByName.String)
	}
	return nil
}

func checkArgsLen(args []string, expected int) error {
	if len(args) != expected {
		return fmt.Errorf("expected '%v' argument, got '%v'\n", expected, len(args))
	}
	return nil
}
