package processor

import (
	"deployRunner/command/deploy"
	"deployRunner/command/image"
	"fmt"
	telegramClient "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
)

const (
	DeployCommandRegexp string = `^/stage (api|spa):(\d{1,4}\.\d{1,4}\.\d{1,4})$`
	ImageCommandRegexp  string = `^/image (api|spa):([\.\d\w-]+)$`
)

type Processor struct {
	bot *telegramClient.BotAPI
}

func New(bot *telegramClient.BotAPI) *Processor {
	return &Processor{bot: bot}
}

func (p *Processor) Process(message *telegramClient.Message) error {
	switch {
	case strings.HasPrefix(message.Text, "/image"):
		expression, _ := regexp.Compile(ImageCommandRegexp)
		if !expression.MatchString(message.Text) {
			p.message(message.Chat.ID, "Can`t understand image command")
			return nil
		}

		arguments := expression.FindAllStringSubmatch(message.Text, -1)

		command := image.New(arguments[0][1], arguments[0][2])

		if command.Run() == nil {
			p.message(message.Chat.ID, "Image building was triggered successfully, please wait notification")
		} else {
			p.message(message.Chat.ID, "Can`t trigger image building")
		}

		return nil
	case strings.HasPrefix(message.Text, "/stage"):
		expression, _ := regexp.Compile(DeployCommandRegexp)
		if !expression.MatchString(message.Text) {
			p.message(message.Chat.ID, "Can`t understand deploy command")
			return nil
		}

		arguments := expression.FindAllStringSubmatch(message.Text, -1)

		command := deploy.New(arguments[0][1], arguments[0][2])

		if commandError := command.Run(); commandError == nil {
			p.message(message.Chat.ID, fmt.Sprintf("New tag %s for %s application successfully applied", arguments[0][2], arguments[0][1]))
		} else {
			p.message(message.Chat.ID, commandError.Error())
		}

		return nil
	default:
		return nil
	}
}

func (p *Processor) message(chatId int64, message string) {
	_, err := p.bot.Send(telegramClient.NewMessage(chatId, message))

	if err != nil {
		fmt.Println(err)
	}
}
