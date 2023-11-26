package release

import (
	"deployRunner/command"
	"deployRunner/config"
	"fmt"
)

const (
	DockerLoginCommand string = "docker login -u=\"%s\" -p=\"%s\" quay.ecpdss.net"
	DockerPullCommand  string = "docker pull quay.ecpdss.net/platform/ecommpay/pp/concept-%s:latest-%s"
	DockerTagCommand   string = "docker tag quay.ecpdss.net/platform/ecommpay/pp/concept-%s:latest-%s quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
	DockerPushCommand  string = "docker push quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
)

type Command struct {
	params *command.ApplicationParams
	quay   *config.Quay
}

func New(application string, tag string, quay *config.Quay) Command {
	return Command{
		params: &command.ApplicationParams{Application: application, Tag: tag},
		quay:   quay,
	}
}

func (c Command) Run() error {
	finalTag := command.ResolveFinalTag(c.params.Tag)

	if _, err := command.Execute(fmt.Sprintf(DockerLoginCommand, c.quay.User, c.quay.Password), ""); err != nil {
		return err
	}

	if _, err := command.Execute(fmt.Sprintf(DockerPullCommand, c.params.Application, c.params.Tag), ""); err != nil {
		return err
	}

	if _, err := command.Execute(fmt.Sprintf(DockerTagCommand, c.params.Application, c.params.Tag, c.params.Application, finalTag), ""); err != nil {
		return err
	}

	if _, err := command.Execute(fmt.Sprintf(DockerPushCommand, c.params.Application, finalTag), ""); err != nil {
		return err
	}

	return nil
}
