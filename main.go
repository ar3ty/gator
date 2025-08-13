package main

import (
	"log"
	"os"

	"github.com/ar3ty/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := state{
		cfg: &cfg,
	}
	cmds := commands{
		handlersMap: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	inline := os.Args
	if len(inline) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	command := command{
		name: inline[1],
		args: inline[2:],
	}

	err = cmds.run(&programState, command)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
