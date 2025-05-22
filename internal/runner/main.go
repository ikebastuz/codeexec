package runner

import (
	"bytes"
	"codeexec/internal/config"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

func Run(lang Lang, code string) (string, string, float64, error) {
	var lg, ok = LangDefinitions[lang]
	if !ok {
		return "", "", 0, fmt.Errorf("unknown language: %s", lang)
	}

	// File
	tmpFile, err := os.CreateTemp("", "tmp-*")
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to create temp file: %s", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(code); err != nil {
		return "", "", 0, fmt.Errorf("failed to write code to temp file: %s", err)
	}

	var stdout, stderr bytes.Buffer
	var command = []string{"run", "--rm", "--pull=never", "-v", fmt.Sprintf("%s:%s", tmpFile.Name(), execFile(lg.fileName))}
	command = append(command, lg.image)
	command = append(command, lg.command...)
	command = append(command, execFile(lg.fileName))

	log.Infof("Executing: docker %s\n", command)

	ctx, cancel := context.WithTimeout(context.Background(), config.PROCESS_TIMEOUT)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", command...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	// TODO: calculate duration only for execution if there is a build step
	duration := time.Since(start).Seconds()

	if ctx.Err() == context.DeadlineExceeded {
		log.Errorf("Timeout exceeded: %s", ctx.Err())
		return "", "", duration, errors.New("timeout exceeded")
	}
	return stdout.String(), stderr.String(), duration, err
}
