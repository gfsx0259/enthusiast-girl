package processor

import (
	"bytes"
	"log"
	"os/exec"
)

func execute(command string, dir string) (output string, err error) {
	var out bytes.Buffer

	cmd := exec.Command("sh", "-c", command)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = &out

	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	return out.String(), nil
}
