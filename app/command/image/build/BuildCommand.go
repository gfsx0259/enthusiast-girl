package build

import (
	"deployRunner/app/command"
	"deployRunner/config"
	"errors"
	"fmt"
)

const (
	RunJobCommand string = "java -jar /bin/jenkins-cli.jar -s https://ci.platformtests.net/ -auth %s:%s build ecommpay/pp/concept-%s/master -p FORCE_TAG=%s"
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
	runCommand := fmt.Sprintf(RunJobCommand, c.sdlc.User, c.sdlc.Password, c.params.Application, c.params.Tag)
	_, err := command.Execute(runCommand, "")

	if err != nil {
		return "", errors.New(err.Error())
	}

	return "Image building started, please wait SDLC notification ⏱", nil
}

func (c Command) String() string {
	return fmt.Sprintf("/image build %s#%s", c.params.Application, c.params.Tag)
}
