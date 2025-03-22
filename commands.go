package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mharkness1/blog_aggregator/internal/api"
	"github.com/mharkness1/blog_aggregator/internal/database"
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

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("failed to set username: %w", err)
	}

	fmt.Println("Username changed successfully.")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("register command requires: register <username>")
	}

	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User created successfully: %s\n", user.Name)
	fmt.Printf("ID: %v", user.ID)

	return nil
}

func handlerDeleteAllUsers(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("delete all users failed: %w", err)
	}
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error in getting users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUser {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}

	return nil
}

func agg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addFeed command requires: addFeed <name> <url>")
	}

	currentUser_id := user.ID

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    currentUser_id,
	})
	if err != nil {
		return fmt.Errorf("error creating feed record: %w", err)
	}

	fmt.Println("Successfully created feed:")
	fmt.Printf("%s\n", feed.Name)
	fmt.Printf("%s\n", feed.Url)

	_, err = s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser_id,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed following record: %w", err)
	}
	fmt.Println("Successfully created feed follow record.")

	return nil
}

func handlerGetAllFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error reading feeds: %w", err)
	}
	for _, feed := range feeds {
		userName, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error matching feeds and users: %w", err)
		}
		fmt.Printf("Name: %s URL: %s  User: %s\n", feed.Name, feed.Url, userName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	feedUrl := cmd.Args[0]
	feedID, err := s.db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("failed to retrieve feed id: %w", err)
	}

	FeedFollow, err := s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedID.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Printf("Name: %s\n", FeedFollow.FeedName)
	fmt.Printf("User: %s\n", FeedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follows: %w", err)
	}

	for _, feed := range feeds {
		feedName, err := s.db.GetFeedFromId(context.Background(), feed.FeedID)
		if err != nil {
			return fmt.Errorf("failed to retrieve feed name: %w", err)
		}
		fmt.Printf("Name: %s\n", feedName.Name)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	feedUrl := cmd.Args[0]
	feed, err := s.db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("failed to find feed details: %w", err)
	}

	err = s.db.DeleteFeedFollowRecord(context.Background(), database.DeleteFeedFollowRecordParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow: %w", err)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	feedFetch, err := s.db.GetNextToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next to fetch: %w", err)
	}

	result, err := api.FetchFeed(context.Background(), feedFetch.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch rss feed: %w", err)
	}

	_, err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:        feedFetch.ID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}
	fmt.Printf("Channel Title: %s\n", result.Channel.Title)
	for _, item := range result.Channel.Item {
		fmt.Printf("%s\n", item.Title)
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
		if err != nil {
			return fmt.Errorf("failed to retrieve user id: %w", err)
		}
		return handler(s, cmd, user)
	}
}

func handlerGetPostsForUser(s *state, cmd command, user database.User) error {
	var limit int32
	if len(cmd.Args) != 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("error parsing limit arguments: %w", err)
		}
		limit = int32(parsedLimit)
	} else {
		limit = 2
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("error retrieving posts for %s: %w", user.ID, err)
	}
	for i, post := range posts {
		fmt.Printf("%d: %s\n", i+1, post.Title)
	}

	return nil
}
