package logs

import (
	"deployRunner/app/command"
	"deployRunner/config"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	GetArtifacts string = "https://%s:%s@ci.platformtests.net/job/ecommpay/job/pp/job/concept-%s/job/master/lastFailedBuild/artifact/*zip*/archive.zip"
	UploadDir    string = "uploads"
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
	fileUrl := fmt.Sprintf(GetArtifacts, c.sdlc.User, c.sdlc.Password, c.params.Application)
	filePath, _ := c.upload(fileUrl, UploadDir, "artifacts.zip")
	return filePath, nil
}

func (c Command) String() string {
	return fmt.Sprintf("/image %s %s", "logs", c.params.Application)
}

func (c Command) upload(url string, dir string, name string) (string, error) {
	err := os.RemoveAll(dir)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(strconv.Itoa(resp.StatusCode))
	}

	filePath := filepath.Join(dir, "/", name)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return filePath, nil
}
