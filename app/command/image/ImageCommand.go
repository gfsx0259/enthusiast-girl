package image

import (
	"deployRunner/command"
	"deployRunner/executor"
	"fmt"
	"os"
)

const (
	GetCrumbCommand   string = "wget -q --auth-no-challenge --user %s --password %s --output-document - 'https://ci.platformtests.net/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)'"
	TriggerJobCommand string = "curl -I -X POST https://sdlc:113ebce2c260e9c6137832606f7a305e06@ci.platformtests.net/job/ecommpay/job/pp/job/concept-%s/view/tags/job/%s/build -H \"%s\""
)

type Command struct {
	params *command.ApplicationParams
}

func New(application string, tag string) Command {
	return Command{
		params: &command.ApplicationParams{Application: application, Tag: tag},
	}
}

func (c Command) Run() error {
	user := os.Getenv("SDLC_USER")
	password := os.Getenv("SDLC_PASSWORD")

	crumb, err := executor.Execute(fmt.Sprintf(GetCrumbCommand, user, password), "")
	if err != nil {
		return err
	}

	if _, err = executor.Execute(fmt.Sprintf(TriggerJobCommand, c.params.Application, c.params.Tag, crumb), ""); err != nil {
		return err
	}

	return nil
}
