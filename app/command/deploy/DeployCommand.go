package deploy

import (
	"deployRunner/command"
	"deployRunner/executor"
	"errors"
	"fmt"
)

const (
	Workdir             string = "okd-pp/pp-%s/overlays/%s/"
	Repository          string = "ssh://git@stash.ecommpay.com:7999/okd/okd-pp.git"
	CustomizeCommand    string = "kustomize edit set image concept-%s=quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
	GitUser             string = "k.popov"
	GitEmail            string = "k.popov@it.ecommpay.com"
	GitConfigureCommand        = "git config --global user.name \"%s\" && git config --global user.email \"%s\""
	CommitMessage       string = "Update version via bot: %s => %s"
)

type Command struct {
	params *command.ApplicationParams
	target string
}

func New(application string, tag string, target string) Command {
	return Command{
		params: &command.ApplicationParams{Application: application, Tag: tag},
		target: target,
	}
}

func (c Command) Run() error {
	workdir := fmt.Sprintf(Workdir, c.params.Application, c.target)

	if _, err := executor.Execute(fmt.Sprintf(GitConfigureCommand, GitUser, GitEmail), ""); err != nil {
		return err
	}

	if _, err := executor.Execute("rm -rf ./okd-pp", ""); err != nil {
		return err
	}
	if _, err := executor.Execute(fmt.Sprintf("git clone %s", Repository), ""); err != nil {
		return err
	}

	if _, err := executor.Execute(fmt.Sprintf(CustomizeCommand, c.params.Application, c.params.Application, c.params.Tag), workdir); err != nil {
		return err
	}
	if _, err := executor.Execute("git add kustomization.yaml", workdir); err != nil {
		return err
	}

	if _, err := executor.Execute(fmt.Sprintf("git commit -m \"%s\"", fmt.Sprintf(CommitMessage, c.params.Application, c.params.Tag)), workdir); err != nil {
		return errors.New("tag already applied, nothing to do")
	}
	if _, err := executor.Execute("git push -u origin HEAD:master -f", workdir); err != nil {
		return err
	}

	return nil
}
