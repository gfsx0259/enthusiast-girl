package telegram

import (
	"deployRunner/command"
	"deployRunner/command/build"
	"deployRunner/command/deploy"
	"deployRunner/command/release"
	"deployRunner/event"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oriser/regroup"
)

const (
	CommandRegexp        = `^/?P<command>(image|deploy) (?P<sub>\w+) (?P<app>api|spa)#(?P<tag>[\.\d\w-]+)$`
	CommandHelp   string = "/help"
	CommandImage  string = "image"
	CommandDeploy string = "deploy"
	ActionBuild   string = "build"
	ActionRelease string = "release"
	EnvStage      string = "stage"
	EnvProd       string = "prod"
)

type Processor struct {
	bot *tgbotapi.BotAPI
}

func NewProcessor(bot *tgbotapi.BotAPI) *Processor {
	return &Processor{bot: bot}
}

func (p *Processor) Process(message *event.Event) error {
	if message.Message == CommandHelp {
		p.message(message.ChatId, p.help())

		return nil
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
		commandInstance := build.New(app, tag)

		if err := commandInstance.Run(); err == nil {
			p.message(message.ChatId, "Image building started, please wait")
		} else {
			p.message(message.ChatId, fmt.Sprintf("Can`t trigger `image` command: %s", err.Error()))
		}

		return nil
	case cmd == CommandImage && sub == ActionRelease:
		if !p.isCommandAvailable(message.FromUsername) {
			p.repeat(message)
			return nil
		}

		commandInstance := release.New(app, tag)

		if err := commandInstance.Run(); err == nil {
			p.message(message.ChatId, fmt.Sprintf("Make final tag %s for %s application", command.ResolveFinalTag(tag), app))
		} else {
			p.message(message.ChatId, fmt.Sprintf("Can`t trigger `image` command: %s", err.Error()))
		}

		return nil
	case cmd == CommandDeploy:
		if sub == EnvProd {
			if !p.isCommandAvailable(message.FromUsername) {
				p.repeat(message)
				return nil
			}
		}

		commandInstance := deploy.New(app, tag, p.normalizeEnvironment(sub))

		if err := commandInstance.Run(); err == nil {
			p.message(message.ChatId, fmt.Sprintf("Deploy for application %s with tag %s runned successfully", app, tag))
		} else {
			p.message(message.ChatId, fmt.Sprintf("Can`t trigger `image` command: %s", err.Error()))
		}

		return nil
	default:
		return nil
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
	_, err := p.bot.Send(messageConfig)

	if err != nil {
		fmt.Println(err)
	}
}

func (p *Processor) isCommandAvailable(username string) bool {
	maintainers := map[string]bool{
		"kopopov": true,
		"fsafsd":  true,
	}

	if _, ok := maintainers[username]; ok {
		return true
	}

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

func (p *Processor) help() string {
	return `
	Hello! I'll tell you how to build and deliver your code.
	
	Firstly, you should create release at jira with your tasks.
	Next, tell the build bot to build your release and return release candidate tag:
    <b>/build#{RELEASE ID}</b>

	Use release candidate tag to trigger image building:
    <b>/image build {APP}#{RC TAG}</b>
	Next, you can use image to deploy it on stage environment:
    <b>/deploy stage {APP}#{RC TAG}</b>

	When testing is finished, complete the release build by placing the final tag in the repository:
	<b>/build#{RELEASE ID}</b>
	Also put the final tag to image register:
	<b>/image release {APP}#{RC TAG}</b>

	Last step, deliver the image to the production environment:
	<b>/deploy prod {APP}#{TAG}</b>
	`
}
