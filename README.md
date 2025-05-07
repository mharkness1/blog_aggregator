# blog_aggregator

This is a cli tool for aggregating RSS feeds and storing their details/content. It was built as part of an introduction to Go course provided by boot.dev

# Overview
# Set-up
# Available Commands
# Future Extensions

gator needs: psql and go installed.

Add more detail.



	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerDeleteAllUsers)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", agg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetAllFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerGetPostsForUser))

Extending the Project

You've done all the required steps, but if you'd like to make this project your own, here are some ideas:

    Add sorting and filtering options to the browse command
    Add pagination to the browse command
    Add concurrency to the agg command so that it can fetch more frequently
    Add a search command that allows for fuzzy searching of posts
    Add bookmarking or liking posts
    Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
    Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
    Write a service manager that keeps the agg command running in the background and restarts it if it crashes
