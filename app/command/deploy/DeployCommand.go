package deploy

import (
	"deployRunner/command"
	"deployRunner/config"
	"errors"
	"fmt"
)

const (
	Workdir             string = "okd-pp/pp-%s/overlays/%s/"
	Repository          string = "ssh://git@stash.ecommpay.com:7999/okd/okd-pp.git"
	CustomizeCommand    string = "kustomize edit set image concept-%s=quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
	GitConfigureCommand        = "git config --global user.name \"%s\" && git config --global user.email \"%s\""
	CommitMessage       string = "Update version via bot: %s => %s"
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

func (c Command) Run() error {
	workdir := fmt.Sprintf(Workdir, c.params.Application, c.target)

	if _, err := command.Execute(fmt.Sprintf(GitConfigureCommand, c.stash.User, c.stash.Email), ""); err != nil {
		return err
	}

	if _, err := command.Execute("rm -rf ./okd-pp", ""); err != nil {
		return err
	}
	if _, err := command.Execute(fmt.Sprintf("git clone %s", Repository), ""); err != nil {
		return err
	}

	if _, err := command.Execute(fmt.Sprintf(CustomizeCommand, c.params.Application, c.params.Application, c.params.Tag), workdir); err != nil {
		return err
	}
	if _, err := command.Execute("git add kustomization.yaml", workdir); err != nil {
		return err
	}

	if _, err := command.Execute(fmt.Sprintf("git commit -m \"%s\"", fmt.Sprintf(CommitMessage, c.params.Application, c.params.Tag)), workdir); err != nil {
		return errors.New("tag already applied, nothing to do")
	}
	if _, err := command.Execute("git push -u origin HEAD:master -f", workdir); err != nil {
		return err
	}

	return nil
}
