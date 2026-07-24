package main

import (
	"Makary01/boot-dev-blog-aggregator/internal/config"
	"fmt"
	"os"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v", err.Error())
		os.Exit(1)
	}

	s := &state{config: cfg}

	cmds := &commands{
		handlersByName: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("Expected 2 argument, got: %v", len(args))
		os.Exit(1)
	}

	err = cmds.run(s, command{name: args[0], args: args[1:]})
	if err != nil {
		fmt.Printf("Error running command: %v", err.Error())
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("The login handler expects a single argument, the username")
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set!")
	return nil
}
