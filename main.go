package main

import (
	"log"
	"os"

	"github.com/mharkness1/blog_aggregator/internal/config"
)

type state struct {
	cfg *config.ConfigStruct
}

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	cliState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(cliState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
