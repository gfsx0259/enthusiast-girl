package inspect

import (
	"deployRunner/app/command"
	"deployRunner/config"
	"errors"
	"fmt"
)

const (
	GetSummary string = "curl -s https://%s:%s@ci.platformtests.net/job/ecommpay/job/pp/job/concept-%s/job/master/lastFailedBuild/ | xmllint --html --xpath \"string(//a[contains(@id, \\\"description-link\\\")]/@data-description)\" 2>/dev/null -"
)

type Command struct {
	params *command.ApplicationParams
	sdlc   *config.Sdlc
}

func New(application string, sdlc *config.Sdlc) Command {
	return Command{
		params: &command.ApplicationParams{Application: application},
		sdlc:   sdlc,
	}
}

func (c Command) Run() (string, error) {
	runCommand := fmt.Sprintf(GetSummary, c.sdlc.User, c.sdlc.Password, c.params.Application)
	output, err := command.Execute(runCommand, "")

	if err != nil {
		return "", errors.New(err.Error())
	}

	return output, nil
}

func (c Command) String() string {
	return ""
}
