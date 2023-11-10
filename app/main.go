package main

import (
	"deployRunner/processor"
	telegramClient "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func main() {
	bot := createBot()

	subscribe(
		bot.GetUpdatesChan(telegramClient.NewUpdate(0)),
		processor.New(bot),
	)
}

func createBot() *telegramClient.BotAPI {
	bot, err := telegramClient.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	return bot
}

func subscribe(updates telegramClient.UpdatesChannel, messageProcessor *processor.Processor) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		err := messageProcessor.Process(update.Message)
		if err != nil {
			log.Print(err)
			return
		}
	}
}
