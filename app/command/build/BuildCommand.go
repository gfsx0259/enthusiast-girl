package build

import (
	"deployRunner/command"
	"deployRunner/executor"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	GetCrumbCommand   string = "curl -u %s:%s -s \"https://ci.platformtests.net/crumbIssuer/api/xml\" | xmllint --format --xpath \"concat(//crumbRequestField,':',//crumb)\" -"
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
	crumb, err := executor.Execute(fmt.Sprintf(GetCrumbCommand, os.Getenv("SDLC_USER"), os.Getenv("SDLC_PASSWORD")), "")
	if err != nil {
		return err
	}

	response, err := executor.Execute(fmt.Sprintf(TriggerJobCommand, c.params.Application, c.params.Tag, crumb), "")
	if err != nil {
		return err
	}

	if strings.Contains(response, "404") {
		return errors.New("job does not exist")
	}

	return nil
}
