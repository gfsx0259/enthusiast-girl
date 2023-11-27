package telegram

import (
	"deployRunner/app/command"
	"deployRunner/app/command/build"
	"deployRunner/app/command/deploy"
	"deployRunner/app/command/help"
	"deployRunner/app/command/release"
	"deployRunner/app/event"
	"deployRunner/config"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oriser/regroup"
	"golang.org/x/exp/slices"
)

const (
	CommandRegexp        = `^/(?P<command>image|deploy) (?P<sub>\w+) (?P<app>api|spa)#(?P<tag>[\.\d\w-]+)$`
	CommandHelp   string = "/help"
	CommandImage  string = "image"
	CommandDeploy string = "deploy"
	ActionBuild   string = "build"
	ActionRelease string = "release"
	EnvStage      string = "stage"
	EnvProd       string = "prod"
)

type Processor struct {
	bot    *tgbotapi.BotAPI
	config *config.Config
}

func NewProcessor(bot *tgbotapi.BotAPI, config *config.Config) *Processor {
	return &Processor{bot, config}
}

func (p *Processor) Process(message *event.Event) error {
	if message.Message == CommandHelp {
		return p.executeCommand(message, help.New())
	}

	expression := regroup.MustCompile(CommandRegexp)
	groups, err := expression.Groups(message.Message)
	if err != nil {
		p.message(message.ChatId, fmt.Sprintf("Can`t understand command, use format %s", CommandRegexp))
		return nil
	}

	cmd := groups["command"]
	sub := groups["sub"]
	app := groups["app"]
	tag := groups["tag"]

	switch {
	case cmd == CommandImage && sub == ActionBuild:
		return p.executeCommand(message, build.New(app, tag, &p.config.Sdlc))
	case cmd == CommandImage && sub == ActionRelease:
		if !p.isCommandAvailable(message) {
			return nil
		}

		return p.executeCommand(message, release.New(app, tag, &p.config.Quay))
	case cmd == CommandDeploy:
		if sub == EnvProd && !p.isCommandAvailable(message) {
			return nil
		}

		return p.executeCommand(message, deploy.New(app, tag, &p.config.Stash, p.normalizeEnvironment(sub)))
	default:
		return nil
	}
}

func (p *Processor) executeCommand(message *event.Event, command command.Command) error {
	if output, err := command.Run(); err == nil {
		p.message(message.ChatId, output)

		return nil
	} else {
		p.message(message.ChatId, fmt.Sprintf("Can`t trigger command: %s", err.Error()))

		return err
	}
}

func (p *Processor) normalizeEnvironment(env string) string {
	if env == EnvProd {
		return "prod-fi1"
	}

	return env
}

func (p *Processor) message(chatId int64, message string) {
	messageConfig := tgbotapi.NewMessage(chatId, message)
	messageConfig.ParseMode = "HTML"

	if _, err := p.bot.Send(messageConfig); err != nil {
		fmt.Println(err)
	}
}

func (p *Processor) isCommandAvailable(message *event.Event) bool {
	if slices.Contains(p.config.Maintainers, message.FromUsername) {
		return true
	}

	p.repeat(message)
	return false
}

func (p *Processor) repeat(message *event.Event) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Approve", message.Message),
		),
	)

	keyboardMessage := tgbotapi.NewMessage(message.ChatId, "Only maintainers @kopopov, @fsafsd can work with a production environment")
	keyboardMessage.ReplyMarkup = keyboard

	if _, err := p.bot.Send(keyboardMessage); err != nil {
		fmt.Println(err)
	}
}
