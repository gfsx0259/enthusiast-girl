package release

import (
	"deployRunner/app/command"
	"deployRunner/config"
	"fmt"
)

const (
	DockerLoginCommand string = "docker login -u=\"%s\" -p=\"%s\" quay.ecpdss.net"
	DockerPullCommand  string = "docker pull quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
	DockerTagCommand   string = "docker tag quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
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

func (c Command) Run() (string, error) {
	if output, err := command.Execute(fmt.Sprintf(DockerLoginCommand, c.quay.User, c.quay.Password), ""); err != nil {
		return output, err
	}

	if output, err := command.Execute(fmt.Sprintf(DockerPullCommand, c.params.Application, c.params.Tag), ""); err != nil {
		return output, err
	}

	finalReleaseTag, _ := command.ResolveFinalTag(c.params.Tag)

	if output, err := command.Execute(fmt.Sprintf(DockerTagCommand, c.params.Application, c.params.Tag, c.params.Application, finalReleaseTag), ""); err != nil {
		return output, err
	}

	if output, err := command.Execute(fmt.Sprintf(DockerPushCommand, c.params.Application, finalReleaseTag), ""); err != nil {
		return output, err
	}

	return fmt.Sprintf("👌 Make final tag %s for %s application", finalReleaseTag, c.params.Application), nil
}

func (c Command) String() string {
	return fmt.Sprintf("/image release %s#%s", c.params.Application, c.params.Tag)
}
