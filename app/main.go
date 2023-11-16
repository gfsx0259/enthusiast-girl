package main

import (
	"deployRunner/event"
	"deployRunner/processor"
	"errors"
	"fmt"
	telegramClient "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		if messageEvent, err := getMessage(update); err == nil {
			err := messageProcessor.Process(messageEvent)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func getMessage(update telegramClient.Update) (*event.Event, error) {
	if update.Message != nil {
		return &event.Event{
			ChatId:       update.Message.Chat.ID,
			FromId:       update.Message.From.ID,
			FromUsername: update.Message.From.UserName,
			Message:      update.Message.Text,
		}, nil
	}
	if update.CallbackQuery.Message != nil {
		callback := update.CallbackQuery

		return &event.Event{
			ChatId:       callback.Message.Chat.ID,
			FromId:       callback.From.ID,
			FromUsername: callback.From.UserName,
			Message:      callback.Data,
		}, nil
	}

	return nil, errors.New("can not resolve message")
}
