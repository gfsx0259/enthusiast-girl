package release

import (
	"deployRunner/command"
	"fmt"
	"os"
)

const (
	DockerLoginCommand string = "docker login -u=\"%s\" -p=\"%s\" quay.ecpdss.net"
	DockerPullCommand  string = "docker pull quay.ecpdss.net/platform/ecommpay/pp/concept-%s:latest-%s"
	DockerTagCommand   string = "docker tag quay.ecpdss.net/platform/ecommpay/pp/concept-%s:latest-%s quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
	DockerPushCommand  string = "docker push quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
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
	finalTag := command.ResolveFinalTag(c.params.Tag)

	if _, err := command.Execute(fmt.Sprintf(DockerLoginCommand, os.Getenv("QUAY_USER"), os.Getenv("QUAY_PASSWORD")), ""); err != nil {
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
