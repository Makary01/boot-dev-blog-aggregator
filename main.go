package main

import (
	"Makary01/boot-dev-blog-aggregator/internal/config"
	"Makary01/boot-dev-blog-aggregator/internal/database"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err.Error())
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err.Error())
		os.Exit(1)
	}

	s := &state{
		config: cfg,
		db:     database.New(db),
	}

	cmds := &commands{
		handlersByName: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("Expected 2 argument, got: %v\n", len(args))
		os.Exit(1)
	}

	err = cmds.run(s, command{name: args[0], args: args[1:]})
	if err != nil {
		fmt.Printf("Error running command: %v\n", err.Error())
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected '1' argument, got '%v'", len(cmd.args))
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
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected '1' argument, got '%v'", len(cmd.args))
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

	fmt.Printf("User created: %v", user)

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	return nil
}
