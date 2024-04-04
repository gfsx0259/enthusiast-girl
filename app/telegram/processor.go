package telegram

import (
	"deployRunner/app/command"
	"deployRunner/app/command/deploy"
	"deployRunner/app/command/help"
	"deployRunner/app/command/image/build"
	"deployRunner/app/command/image/inspect"
	"deployRunner/app/command/image/logs"
	"deployRunner/app/command/image/release"
	"deployRunner/app/event"
	"deployRunner/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oriser/regroup"
	"golang.org/x/exp/slices"
	"strings"
)

const (
	CommandRegexp        = `^/(?P<command>image|deploy) (?P<sub>\w+) (?P<app>api|spa)#?(?P<tag>[\.\d\w-]+)?$`
	CommandHelp   string = "/help"
	CommandImage  string = "image"
	CommandDeploy string = "deploy"
	ActionBuild   string = "build"
	ActionInspect string = "inspect"
	ActionLogs    string = "logs"
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
		return p.executeCommand(message, help.New(), "")
	}

	expression := regroup.MustCompile(CommandRegexp)
	groups, err := expression.Groups(message.Message)
	if err != nil {
		p.message(message.ChatId, fmt.Sprintf("Can`t understand command, use format %s", CommandRegexp), "", "")
		return nil
	}

	cmd := groups["command"]
	sub := groups["sub"]
	app := groups["app"]
	tag := groups["tag"]

	switch {
	case cmd == CommandImage && sub == ActionBuild:
		buildCommand := build.New(app, tag, &p.config.Sdlc)
		deployCommand := deploy.New(app, tag, &p.config.Stash, EnvStage)

		return p.executeCommand(
			message,
			buildCommand,
			deployCommand.String(),
		)
	case cmd == CommandImage && sub == ActionInspect:
		inspectCommand := inspect.New(app, &p.config.Sdlc)
		logsCommand := logs.New(app, &p.config.Sdlc)

		return p.executeCommand(
			message,
			inspectCommand,
			logsCommand.String(),
		)
	case cmd == CommandImage && sub == ActionLogs:
		logsCommand := logs.New(app, &p.config.Sdlc)
		logsResult, _ := logsCommand.Run()

		p.file(message.ChatId, logsResult)

		return nil
	case cmd == CommandImage && sub == ActionRelease:
		if !p.isCommandAvailable(message) {
			return nil
		}

		var finalTag string

		if finalTag, err = command.ResolveFinalTag(tag); err != nil {
			p.message(message.ChatId, err.Error(), "", "")
			return err
		}

		releaseCommand := release.New(app, tag, &p.config.Quay)
		deployCommand := deploy.New(app, finalTag, &p.config.Stash, EnvProd)

		return p.executeCommand(
			message,
			releaseCommand,
			deployCommand.String(),
		)
	case cmd == CommandDeploy:
		deployCommand := deploy.New(app, tag, &p.config.Stash, p.normalizeEnvironment(sub))

		if sub == EnvProd {
			if !p.isCommandAvailable(message) {
				return nil
			}

			return p.executeCommand(
				message,
				deployCommand,
				"",
			)
		} else {
			return p.executeCommand(
				message,
				deployCommand,
				release.New(app, tag, &p.config.Quay).String(),
			)
		}
	default:
		return nil
	}
}

func (p *Processor) executeCommand(message *event.Event, command command.Command, nextCommand string) error {
	if output, err := command.Run(); err == nil {
		p.message(message.ChatId, output, nextCommand, "")
		return nil
	} else {
		p.message(message.ChatId, fmt.Sprintf("Can`t trigger command: %s", err.Error()), "", "")
		return err
	}
}

func (p *Processor) normalizeEnvironment(env string) string {
	if env == EnvProd {
		return "prod-fi1"
	}

	return env
}

func (p *Processor) message(chatId int64, message string, nextCommand string, nextCommandTitle string) {
	messageConfig := tgbotapi.NewMessage(chatId, message)
	messageConfig.ParseMode = "HTML"

	if nextCommand != "" {
		messageConfig.ReplyMarkup = p.createKeyboard(nextCommandTitle, nextCommand)
	}

	if _, err := p.bot.Send(messageConfig); err != nil {
		fmt.Println(err)
	}
}

func (p *Processor) file(chatId int64, fileUrl string) {
	file := tgbotapi.FilePath(fileUrl)
	messageConfig := tgbotapi.NewDocument(chatId, file)

	if _, err := p.bot.Send(messageConfig); err != nil {
		fmt.Println(err)
	}
}

func (p *Processor) isCommandAvailable(message *event.Event) bool {
	if slices.Contains(p.config.Maintainers, message.FromUsername) {
		return true
	}

	p.message(
		message.ChatId,
		p.buildApproveRequiredString(),
		message.Message,
		"Approve",
	)
	return false
}

func (p *Processor) buildApproveRequiredString() string {
	var maintainers []string

	for _, maintainer := range p.config.Maintainers {
		maintainers = append(maintainers, fmt.Sprintf("@%s", maintainer))
	}

	maintainersString := strings.Join(maintainers, ", ")

	return fmt.Sprintf("Only maintainers %s can work with a production environment", maintainersString)
}

func (p *Processor) createKeyboard(commandTitle string, command string) tgbotapi.InlineKeyboardMarkup {
	var title string

	if commandTitle != "" {
		title = commandTitle
	} else {
		title = "Run " + command
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(title, command),
		),
	)
}
