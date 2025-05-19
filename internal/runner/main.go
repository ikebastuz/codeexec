package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func Run(lang Lang, code string) (stdoutStr string, stderrStr string, err error) {
	var lg, ok = langDefinition[lang]
	if !ok {
		return "", "", fmt.Errorf("unknown language: %s", lang)
	}

	// File
	tmpFile, err := os.CreateTemp("", "tmp-*")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp file: %s", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(code); err != nil {
		return "", "", fmt.Errorf("failed to write code to temp file: %s", err)
	}

	var stdout, stderr bytes.Buffer
	var command = []string{"run", "--rm", "-v", fmt.Sprintf("%s:%s", tmpFile.Name(), execFile(lg.ext))}
	command = append(command, lg.image)
	command = append(command, lg.command...)
	command = append(command, execFile(lg.ext))

	log.Infof("Executing: docker %s\n", command)

	cmd := exec.Command("docker", command...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// TODO: pull images on init
	err = cmd.Run()

	return stdout.String(), stderr.String(), err
}

func execFile(ext string) string {
	return fmt.Sprintf("%s/%s", WORKDIR, execName(ext))
}

func execName(ext string) string {
	return fmt.Sprintf("%s.%s", EXEC_FILE, ext)
}
