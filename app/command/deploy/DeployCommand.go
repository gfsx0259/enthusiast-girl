package deploy

import (
	"deployRunner/app/command"
	"deployRunner/config"
	"errors"
	"fmt"
	"strings"
)

const (
	Workdir                   string = "okd-pp/pp-%s/overlays/%s/"
	Repository                string = "ssh://git@stash.ecommpay.com:7999/okd/okd-pp.git"
	CustomizeImageCommand     string = "kustomize edit set image concept-%s=quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
	CustomizeConfigmapCommand string = "kustomize edit set configmap pp-api-config --from-literal=%s=%s"
	GitConfigureCommand       string = "git config --global user.name \"%s\" && git config --global user.email \"%s\""
	MakeEnvironmentJs         string = "PP_SPA_SENTRY_RELEASE=%s envsubst < environment.js.envsubst > environment.js"
	CommitMessage             string = "Update version via bot: %s => %s"
)

type Command struct {
	params *command.ApplicationParams
	stash  *config.Stash
	target string
}

func New(application string, tag string, stash *config.Stash, target string) Command {
	return Command{
		params: &command.ApplicationParams{Application: application, Tag: tag},
		target: target,
		stash:  stash,
	}
}

func (c Command) Run() (string, error) {
	if output, err := command.Execute(fmt.Sprintf(GitConfigureCommand, c.stash.User, c.stash.Email), ""); err != nil {
		return output, err
	}
	if output, err := command.Execute("rm -rf ./okd-pp", ""); err != nil {
		return output, err
	}
	if output, err := command.Execute(fmt.Sprintf("git clone %s", Repository), ""); err != nil {
		return output, err
	}

	workdir := fmt.Sprintf(Workdir, c.params.Application, c.target)

	if output, err := command.Execute(fmt.Sprintf(CustomizeImageCommand, c.params.Application, c.params.Application, c.params.Tag), workdir); err != nil {
		return output, err
	}

	finalTag := c.getSentryReleaseTag(c.params.Tag)

	if c.params.Application == "spa" {
		if output, err := command.Execute(fmt.Sprintf(MakeEnvironmentJs, finalTag), workdir); err != nil {
			return output, err
		}
	}
	if c.params.Application == "api" {
		if output, err := command.Execute(fmt.Sprintf(CustomizeConfigmapCommand, "PP_API_RELEASE", finalTag), workdir); err != nil {
			return output, err
		}
	}

	if output, err := command.Execute("git add --all", workdir); err != nil {
		return output, err
	}
	if output, err := command.Execute(fmt.Sprintf("git commit -m \"%s\"", fmt.Sprintf(CommitMessage, c.params.Application, c.params.Tag)), workdir); err != nil {
		return output, errors.New("tag already applied, nothing to do")
	}
	if output, err := command.Execute("git push -u origin HEAD:master -f", workdir); err != nil {
		return output, err
	}

	return fmt.Sprintf("Deploy for application %s with tag %s runned successfully, please wait ARGO notification ⏱", c.params.Application, c.params.Tag), nil
}

func (c Command) getSentryReleaseTag(tag string) string {
	if strings.Contains(tag, "-rc") {
		sentryReleaseTag, _ := command.ResolveFinalTag(c.params.Tag)
		return sentryReleaseTag
	} else {
		return tag
	}
}

func (c Command) String() string {
	return fmt.Sprintf("/deploy %s %s#%s", c.target, c.params.Application, c.params.Tag)
}
