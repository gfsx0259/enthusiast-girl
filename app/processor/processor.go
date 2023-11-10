package processor

import (
	"fmt"
	telegramClient "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
)

const (
	DeployCommandRegexp           = `^/stage (api|spa):(\d{1,4}\.\d{1,4}\.\d{1,4})$`
	DeployConfigurationRepository = "ssh://git@stash.ecommpay.com:7999/okd/okd-pp.git"
)

type Processor struct {
	bot *telegramClient.BotAPI
}

func New(bot *telegramClient.BotAPI) *Processor {
	return &Processor{bot: bot}
}

func (p *Processor) Process(message *telegramClient.Message) error {
	switch {
	case strings.HasPrefix(message.Text, "/stage"):
		re, _ := regexp.Compile(DeployCommandRegexp)

		if !re.MatchString(message.Text) {
			_, _ = p.bot.Send(telegramClient.NewMessage(message.Chat.ID, "Can`t understand deploy command"))
			return nil
		}

		res := re.FindAllStringSubmatch(message.Text, -1)

		tag := res[0][2]
		app := res[0][1]

		workdir := fmt.Sprintf("okd-pp/pp-%s/overlays/stage/", app)
		commitMessage := fmt.Sprintf("Update version via bot: %s => %s", app, tag)

		_, err := execute("rm -rf ./okd-pp", "")
		if err != nil {
			return err
		}

		_, err = execute("git clone "+DeployConfigurationRepository, "")
		if err != nil {
			return err
		}

		_, err = execute("kustomize edit set image concept-"+app+"=quay.ecpdss.net/platform/ecommpay/pp/concept-"+app+":"+tag, workdir)
		if err != nil {
			return err
		}

		_, err = execute("git add kustomization.yaml", "okd-pp/pp-"+app+"/overlays/stage/")
		if err != nil {
			return err
		}

		_, err = execute("git config --global user.name \"k.popov\" && git config --global user.email \"k.popov@it.ecommpay.com\"", "")
		if err != nil {
			return err
		}

		_, err = execute(fmt.Sprintf("git commit -m \"%s\"", commitMessage), "okd-pp/pp-"+app+"/overlays/stage/")
		if err != nil {
			return err
		}

		_, err = execute("git push -u origin HEAD:master -f", "okd-pp/pp-"+app+"/overlays/stage/")
		if err != nil {
			return err
		}

		_, err = p.bot.Send(telegramClient.NewMessage(message.Chat.ID, fmt.Sprintf("New tag %s for %s application successfully applied", tag, app)))
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}
