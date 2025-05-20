package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func Run(lang Lang, code string, c chan Job) {
	var lg, ok = langDefinition[lang]
	if !ok {
		c <- Job{
			Stdout: "",
			Stderr: "",
			Error:  fmt.Errorf("unknown language: %s", lang),
		}
		return
	}

	// File
	tmpFile, err := os.CreateTemp("", "tmp-*")
	if err != nil {
		c <- Job{
			Stdout: "",
			Stderr: "",
			Error:  fmt.Errorf("failed to create temp file: %s", err),
		}
		return
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(code); err != nil {
		c <- Job{
			Stdout: "",
			Stderr: "",
			Error:  fmt.Errorf("failed to write code to temp file: %s", err),
		}
	}

	var stdout, stderr bytes.Buffer
	var command = []string{"run", "--rm", "-v", fmt.Sprintf("%s:%s", tmpFile.Name(), execFile(lg.fileName))}
	command = append(command, lg.image)
	command = append(command, lg.command...)
	command = append(command, execFile(lg.fileName))

	log.Infof("Executing: docker %s\n", command)

	cmd := exec.Command("docker", command...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// TODO: pull images on init
	err = cmd.Run()

	c <- Job{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Error:  err,
	}
}

func execFile(fileName string) string {
	return fmt.Sprintf("%s/%s", WORKDIR, fileName)
}
