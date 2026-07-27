package main

import (
	"Makary01/boot-dev-blog-aggregator/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func handlerAgg(s *state, cmd command) error {
	if err := checkArgsLen(cmd.args, 1); err != nil {
		return err
	}
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s\n", cmd.args[0])

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if err := checkArgsLen(cmd.args, 2); err != nil {
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

	followCmd := command{
		name: "follow",
		args: []string{cmd.args[1]},
	}
	return handlerFollow(s, followCmd, user)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if err := checkArgsLen(cmd.args, 1); err != nil {
		return err
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("`%s` is now followed by %s\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if err := checkArgsLen(cmd.args, 0); err != nil {
		return err
	}

	followings, err := s.db.GetFeedFollowings(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	for _, f := range followings {
		fmt.Println(f.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if err := checkArgsLen(cmd.args, 1); err != nil {
		return err
	}

	return s.db.DeleteFeedFollowing(
		context.Background(),
		database.DeleteFeedFollowingParams{Name: user.Name, Url: cmd.args[0]},
	)
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if err := checkArgsLen(cmd.args, 0, 1); err != nil {
		return err
	}
	limit := 2
	if len(cmd.args) == 1 {
		l, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = l
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Printf("Post %s:\n", p.Title)
		if p.Description.Valid {
			fmt.Printf("Description: %s\n", p.Description.String)
		}
		if p.PublishedAt.Valid {
			fmt.Printf("Published At: %s\n", p.PublishedAt.Time.Format(time.Stamp))
		}
		fmt.Println()

		resp, err := http.Get(p.Url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("Body: \n%s\n", string(body))
	}
	return nil
}

func scrapeFeeds(s *state) {
	f, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("Error when querying feed: '%v'", err.Error())
		return
	}

	rss, err := fetchFeed(context.Background(), f.Url)
	if err != nil {
		fmt.Printf("Error when fetching feed: '%v'", err.Error())
		return
	}

	params := database.MarkFeedFetchedParams{
		ID:        f.ID,
		UpdatedAt: time.Now(),
		LastFetchAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		fmt.Printf("Error when marking feed fetched: '%v'", err.Error())
		return
	}

	for _, item := range rss.Channel.Item {
		t, err := formatTime(item.PubDate)
		pubAt := sql.NullTime{
			Time:  t,
			Valid: err == nil,
		}
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			FeedID:      f.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			PublishedAt: pubAt,
			Title:       item.Title,
			Description: sql.NullString{String: item.Description, Valid: item.Description == ""},
			Url:         item.Link,
		}
		err = s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				switch pqErr.Code {
				case "23505":
				default:
					fmt.Printf("Error: %v\n", err.Error())
				}
			}
		}
	}
}

func checkArgsLen(args []string, expected ...int) error {
	for _, ex := range expected {
		if len(args) == ex {
			return nil
		}
	}
	expectedStrs := []string{}
	for _, ex := range expected {
		expectedStrs = append(expectedStrs, strconv.Itoa(ex))
	}

	return fmt.Errorf("expected: '%s' argument(s), got: '%n'\n", strings.Join(expectedStrs, "' OR '"), len(args))
}

func ensureLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		u, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			fmt.Println("You must be logged in to use this command!")
			return err
		}
		return handler(s, cmd, u)
	}
}

func formatTime(val string) (t time.Time, err error) {
	for _, layout := range []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
	} {
		t, err = time.Parse(layout, val)
		if err == nil {
			return t, nil
		}
	}
	return t, err
}
