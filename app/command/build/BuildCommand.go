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
	crumbCommand := fmt.Sprintf(GetCrumbCommand, c.sdlc.User, c.sdlc.Password)

	crumb, err := command.Execute(crumbCommand, "")
	if err != nil {
		return "", errors.New(fmt.Sprintf("error occurred while fetching jenknns crumb: %s", err.Error()))
	}

	triggerCommand := fmt.Sprintf(TriggerJobCommand, c.sdlc.User, c.sdlc.Token, c.params.Application, c.params.Tag, strings.TrimSuffix(crumb, "\n"))

	response, err := command.Execute(triggerCommand, "")
	if err != nil {
		return "", errors.New(fmt.Sprintf("error occurred while triggering jenkins job: %s", err.Error()))
	}

	if strings.Contains(response, "404") {
		return "", errors.New("job does not exist")
	}

	return "Image building started, please wait SDLC notification ⏱", nil
}

func (c Command) String() string {
	return fmt.Sprintf("/image build %s#%s", c.params.Application, c.params.Tag)
}
