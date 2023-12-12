package command

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Command interface {
	Run() (string, error)
	String() string
}

type ApplicationParams struct {
	Application string
	Tag         string
}

func ResolveFinalTag(tag string) (string, error) {
	if !strings.Contains(tag, "-") {
		return "", errors.New("tag does not contain rc suffix")
	}

	return tag[:strings.IndexByte(tag, '-')], nil
}

func Execute(command string, dir string) (output string, err error) {
	var outBuffer bytes.Buffer
	var errorBuffer bytes.Buffer

	log.Println(fmt.Sprintf("Run command: %s", command))

	cmd := exec.Command("sh", "-c", command)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errorBuffer

	err = cmd.Run()

	if err != nil {
		return outBuffer.String(), errors.New(errorBuffer.String())
	}

	return outBuffer.String(), nil
}
