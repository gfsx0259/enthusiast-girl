package processor

import (
	"deployRunner/command"
	"deployRunner/command/build"
	"deployRunner/command/deploy"
	"deployRunner/command/release"
	"deployRunner/event"
	"fmt"
	telegramClient "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
)

const (
	DeployCommandRegexp string = `^/deploy (stage|prod) (api|spa)#([\.\d\w-]+)$`
	ImageCommandRegexp  string = `^/image (build|release) (api|spa)#([\.\d\w-]+)$`
	ActionBuild         string = "build"
	ActionRelease       string = "release"
	EnvProd             string = "prod"
)

type Processor struct {
	bot *telegramClient.BotAPI
}

func New(bot *telegramClient.BotAPI) *Processor {
	return &Processor{bot: bot}
}

func (p *Processor) Process(message *event.Event) error {
	switch {
	case strings.HasPrefix(message.Message, "/image"):
		expression, _ := regexp.Compile(ImageCommandRegexp)
		if !expression.MatchString(message.Message) {
			p.message(message.ChatId, fmt.Sprintf("Can`t understand `image` command, use format: %s", expression.String()))
			return nil
		}

		arguments := expression.FindAllStringSubmatch(message.Message, -1)

		action := arguments[0][1]
		app := arguments[0][2]
		tag := arguments[0][3]

		var cmd command.Command
		var successMessage string

		switch action {
		case ActionBuild:
			cmd = build.New(app, tag)
			successMessage = "Image building started, please wait"
		case ActionRelease:
			if !p.isCommandAvailable(message.FromUsername) {
				p.repeat(message)
				return nil
			}

			cmd = release.New(app, tag)
			successMessage = fmt.Sprintf("Make final tag %s for %s application", command.ResolveFinalTag(tag), app)
		}

		if err := cmd.Run(); err == nil {
			p.message(message.ChatId, successMessage)
		} else {
			p.message(message.ChatId, fmt.Sprintf("Can`t trigger `image` command: %s", err.Error()))
		}

		return nil
	case strings.HasPrefix(message.Message, "/deploy"):
		expression, _ := regexp.Compile(DeployCommandRegexp)
		if !expression.MatchString(message.Message) {
			p.message(message.ChatId, fmt.Sprintf("Can`t understand `deploy` command, use format %s", expression.String()))
			return nil
		}

		arguments := expression.FindAllStringSubmatch(message.Message, -1)

		if arguments[0][1] == EnvProd {
			if !p.isCommandAvailable(message.FromUsername) {
				p.repeat(message)
				return nil
			}
		}

		env := p.normalizeEnvironment(arguments[0][1])
		app := arguments[0][2]
		tag := arguments[0][3]

		cmd := deploy.New(app, tag, env)

		if err := cmd.Run(); err == nil {
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
	_, err := p.bot.Send(telegramClient.NewMessage(chatId, message))

	if err != nil {
		fmt.Println(err)
	}
}

func (p *Processor) isCommandAvailable(username string) bool {
	maintainers := map[string]bool{
		"kopopov":   true,
		"r.petunin": true,
	}

	if _, ok := maintainers[username]; ok {
		return true
	}

	return false
}

func (p *Processor) repeat(message *event.Event) {
	keyboard := telegramClient.NewInlineKeyboardMarkup(
		telegramClient.NewInlineKeyboardRow(
			telegramClient.NewInlineKeyboardButtonData("Approve", message.Message),
		),
	)

	keyboardMessage := telegramClient.NewMessage(message.ChatId, "Only maintainers @kopopov, @fsafsd can work with a production environment")
	keyboardMessage.ReplyMarkup = keyboard

	if _, err := p.bot.Send(keyboardMessage); err != nil {
		fmt.Println(err)
	}
}
