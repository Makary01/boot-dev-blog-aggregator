package main

import (
	"Makary01/boot-dev-blog-aggregator/internal/config"
	"Makary01/boot-dev-blog-aggregator/internal/database"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
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
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", ensureLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", ensureLoggedIn(handlerFollow))
	cmds.register("following", ensureLoggedIn(handlerFollowing))
	cmds.register("unfollow", ensureLoggedIn(handlerUnfollow))
	cmds.register("browse", ensureLoggedIn(handlerBrowse))

	args := os.Args[1:]
	err = cmds.run(s, command{name: args[0], args: args[1:]})
	if err != nil {
		fmt.Printf("Error running command: %v\n", err.Error())
		os.Exit(1)
	}
}
