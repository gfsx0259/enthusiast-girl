package telegram

import (
	"deployRunner/app/event"
	"deployRunner/config"
	"errors"
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Listener struct {
	config *config.Config
}

func NewListener(config *config.Config) *Listener {
	return &Listener{config}
}

func (l *Listener) Listen() {
	bot := l.createBot()

	l.subscribe(
		bot.GetUpdatesChan(telegram.NewUpdate(0)),
		NewProcessor(bot, l.config),
	)
}

func (l *Listener) createBot() *telegram.BotAPI {
	bot, err := telegram.NewBotAPI(l.config.Telegram.Token)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	return bot
}

func (l *Listener) subscribe(updates telegram.UpdatesChannel, eventProcessor *Processor) {
	for update := range updates {
		if receivedEvent, err := toEvent(update); err == nil {
			err := eventProcessor.Process(receivedEvent)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func toEvent(update telegram.Update) (*event.Event, error) {
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
