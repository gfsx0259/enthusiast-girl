package main

import (
	"deployRunner/config"
	"deployRunner/telegram"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	telegram.NewListener(cfg).Listen()
}
