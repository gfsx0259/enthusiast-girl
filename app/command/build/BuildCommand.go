package build

import (
	"deployRunner/app/command"
	"deployRunner/config"
	"errors"
	"fmt"
	"strings"
)

const (
	GetCrumbCommand   string = "curl -u %s:%s -s \"https://ci.platformtests.net/crumbIssuer/api/xml\" | xmllint --format --xpath \"concat(//crumbRequestField,':',//crumb)\" -"
	TriggerJobCommand string = "curl -X POST https://%s:%s@ci.platformtests.net/job/ecommpay/job/pp/job/concept-%s/job/master/buildWithParameters -d \"FORCE_TAG=%s\" -H \"%s\""
)

type Command struct {
	params *command.ApplicationParams
	sdlc   *config.Sdlc
}

func New(application string, tag string, sdlc *config.Sdlc) Command {
	return Command{
		params: &command.ApplicationParams{Application: application, Tag: tag},
		sdlc:   sdlc,
	}
}

func (c Command) Run() (string, error) {
	crumb, err := command.Execute(fmt.Sprintf(GetCrumbCommand, c.sdlc.User, c.sdlc.Password), "")
	if err != nil {
		return "", errors.New(fmt.Sprintf("error occurred while fetching jenknns crumb: %s", err.Error()))
	}

	response, err := command.Execute(fmt.Sprintf(TriggerJobCommand, c.sdlc.User, c.sdlc.Token, c.params.Application, c.params.Tag, crumb), "")
	if err != nil {
		return "", errors.New(fmt.Sprintf("error occurred while triggering jenkins job: %s", err.Error()))
	}

	if strings.Contains(response, "404") {
		return "", errors.New("job does not exist")
	}

	return "Image building started, please wait", nil
}
