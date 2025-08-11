package main

import (
	"fmt"
	"log"

	"github.com/ar3ty/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	err = cfg.SetUser("koenig")
	if err != nil {
		log.Fatalf("cannot set user : %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Print(cfg)
}
