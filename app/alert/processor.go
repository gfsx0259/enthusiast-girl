package alert

import (
	"deployRunner/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
)

type Processor struct {
	config *config.Config
}

func NewProcessor(config *config.Config) *Processor {
	return &Processor{config}
}

func (p *Processor) AcceptHook(w http.ResponseWriter, r *http.Request) {
	bot := p.createBot()

	messageConfig := tgbotapi.NewPhoto(
		int64(p.config.Alert.ChatId),
		tgbotapi.FilePath("/static/sentry.png"),
	)

	hookStructure := NewStructure(r.Body)
	if hookStructure != nil {
		rows := []string{
			"Sentry alert\n",
			"Alert resolved for projects: " + strings.ToUpper(strings.Join(hookStructure.Metric.Projects, ",")),
			hookStructure.Title,
			hookStructure.Text,
			hookStructure.Url,
		}

		messageConfig.Caption = strings.Join(rows, "\n")

		if _, err := bot.Send(messageConfig); err != nil {
			fmt.Println(err)
		}
		if _, err := io.WriteString(w, "ok"); err != nil {
			fmt.Println(err)
		}
	}
}

func (p *Processor) createBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(p.config.Telegram.Token)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	return bot
}
