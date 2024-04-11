package main

import (
	"deployRunner/app/alert"
	"deployRunner/app/telegram"
	"deployRunner/config"
	"errors"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	go listenHooks(cfg)
	listenCommands(cfg)
}

func listenCommands(cfg *config.Config) {
	telegram.NewListener(cfg).Listen()
}

func listenHooks(cfg *config.Config) {
	http.HandleFunc("/hook", func(writer http.ResponseWriter, request *http.Request) {
		alert.NewProcessor(cfg).AcceptHook(writer, request)
	})

	err := http.ListenAndServe(":"+cfg.Alert.HookPort, nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server closed: %s", err)
	} else if err != nil {
		log.Fatalf("Can not start server: %s", err)
	} else {
		log.Println("Server started")
	}
}
