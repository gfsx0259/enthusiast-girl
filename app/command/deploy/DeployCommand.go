package deploy

import (
	"deployRunner/app/command"
	"deployRunner/config"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	TempDir                   string = "/tmp/projects"
	Workdir                   string = "okd-pp/pp-%s/overlays/%s/"
	Repository                string = "ssh://git@stash.ecommpay.com:7999/okd/okd-pp.git"
	WorkdirValuesStage        string = "nl1/lui-common-stage/app-alpha"
	RepositoryValues          string = "git@%s:applications/platform/payment-page/values/%s.git"
	ValuesGeneratedFile       string = "values.generated.yaml"
	CustomizeImageCommand     string = "kustomize edit set image concept-%s=quay.ecpdss.net/platform/ecommpay/pp/concept-%s:%s"
	CustomizeConfigmapCommand string = "kustomize edit set configmap pp-api-config --from-literal=%s=%s"
	GitConfigureCommand       string = "git config --global user.name \"%s\" && git config --global user.email \"%s\""
	MakeEnvironmentJs         string = "PP_SPA_SENTRY_RELEASE=%s envsubst < environment.js.envsubst > environment.js"
	CommitMessage             string = "PP-123 Update version via bot: %s => %s"
)

type Api struct {
	Api ImageDefinition `yaml:"api"`
}

type ImageDefinition struct {
	Image TagDefinition `yaml:"image"`
}

type TagDefinition struct {
	Tag string `yaml:"tag"`
}

type Command struct {
	params *command.ApplicationParams
	stash  *config.Git
	target string
}

func New(application string, tag string, stash *config.Git, target string) Command {
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

	if c.target == "prod" {
		if err := c.fetch("okd-pp", Repository); err != nil {
			return "", err
		}

		workdir := TempDir + "/" + fmt.Sprintf(Workdir, c.params.Application, c.target)

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

		if err := c.push(workdir); err != nil {
			return "", err
		}
	} else {
		if err := c.fetch(c.params.Application, fmt.Sprintf(RepositoryValues, c.stash.Host, c.params.Application)); err != nil {
			return "", err
		}

		workdir := TempDir + "/" + c.params.Application + "/" + WorkdirValuesStage + "/"

		if err := c.overrideValuesFile(workdir); err != nil {
			return "", err
		}

		if err := c.push(workdir); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("Deploy for application %s with tag %s had run successfully, please wait for ARGO notification ⏱", c.params.Application, c.params.Tag), nil
}

func (c Command) fetch(targetDir string, repo string) error {
	if _, err := command.Execute(fmt.Sprintf("mkdir -p %s", TempDir), ""); err != nil {
		return err
	}
	if _, err := command.Execute(fmt.Sprintf("rm -rf ./%s", targetDir), TempDir); err != nil {
		return err
	}
	if _, err := command.Execute(fmt.Sprintf("git clone %s", repo), TempDir); err != nil {
		return err
	}

	return nil
}

func (c Command) push(workdir string) error {
	if _, err := command.Execute("git add --all", workdir); err != nil {
		return err
	}
	if _, err := command.Execute(fmt.Sprintf("git commit -m \"%s\"", fmt.Sprintf(CommitMessage, c.params.Application, c.params.Tag)), workdir); err != nil {
		return errors.New("tag had been already applied, there is nothing to do")
	}
	if _, err := command.Execute("git push -u origin HEAD:master -f", workdir); err != nil {
		return err
	}

	return nil
}

func (c Command) overrideValuesFile(filePath string) error {
	configStructure := c.buildConfigStructure()

	configYaml, err := yaml.Marshal(&configStructure)
	if err != nil {
		return errors.New("Error while marshaling: " + err.Error())
	}

	err = os.WriteFile(filePath+ValuesGeneratedFile, configYaml, 0664)

	if err != nil {
		return errors.New("Unable to write data into the file: " + err.Error())
	}

	return nil
}

func (c Command) buildConfigStructure() interface{} {
	if c.params.Application == "api" {
		return Api{
			Api: ImageDefinition{
				Image: TagDefinition{
					Tag: c.params.Tag,
				},
			},
		}
	} else {
		return ImageDefinition{
			Image: TagDefinition{
				Tag: c.params.Tag,
			},
		}
	}
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
