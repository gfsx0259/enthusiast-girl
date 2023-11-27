package command

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type Command interface {
	Run() (string, error)
}

type ApplicationParams struct {
	Application string
	Tag         string
}

func ResolveFinalTag(tag string) string {
	return tag[:strings.IndexByte(tag, '-')]
}

func Execute(command string, dir string) (output string, err error) {
	var outBuffer bytes.Buffer

	cmd := exec.Command("sh", "-c", command)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = &outBuffer

	err = cmd.Run()

	if err != nil {
		log.Println(err)
	}

	return outBuffer.String(), err
}
