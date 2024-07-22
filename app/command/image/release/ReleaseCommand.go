package release

import (
	"deployRunner/app/command"
	"deployRunner/config"
	"fmt"
	"strings"
)

const (
	DockerLoginCommand string = "docker login %s -u=\"%s\" -p=\"%s\""
	DockerPullCommand  string = "docker pull %s/ecommpay/pp/%s:%s"
	DockerTagCommand   string = "docker tag %s/ecommpay/pp/%s:%s %s/ecommpay/pp/%s:%s"
	DockerPushCommand  string = "docker push %s/ecommpay/pp/%s:%s"
)

type Command struct {
	params *command.ApplicationParams
	registry *config.Registry
}

func New(application string, tag string, registry *config.Registry) Command {
	return Command{
		params: &command.ApplicationParams{Application: application, Tag: tag},
		registry: registry,
	}
}

func (c Command) Run() (string, error) {
    registryPaths := strings.Split(c.registry.Host, "/")

	if output, err := command.Execute(fmt.Sprintf(DockerLoginCommand, c.registry.User, c.registry.Password, registryPaths[0]), ""); err != nil {
		return output, err
	}

	if output, err := command.Execute(fmt.Sprintf(DockerPullCommand, c.registry.Host, c.params.Application, c.params.Tag), ""); err != nil {
		return output, err
	}

	finalReleaseTag, _ := command.ResolveFinalTag(c.params.Tag)

	if output, err := command.Execute(fmt.Sprintf(DockerTagCommand, c.registry.Host, c.params.Application, c.params.Tag, c.registry.Host, c.params.Application, finalReleaseTag), ""); err != nil {
		return output, err
	}

	if output, err := command.Execute(fmt.Sprintf(DockerPushCommand, c.registry.Host, c.params.Application, finalReleaseTag), ""); err != nil {
		return output, err
	}

	return fmt.Sprintf("👌 Make final tag %s for %s application", finalReleaseTag, c.params.Application), nil
}

func (c Command) String() string {
	return fmt.Sprintf("/image release %s#%s", c.params.Application, c.params.Tag)
}
