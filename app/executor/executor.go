package executor

import (
	"bytes"
	"log"
	"os/exec"
)

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
