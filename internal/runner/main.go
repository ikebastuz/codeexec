package runner

import (
	"bytes"
	"codeexec/internal/config"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	command := makeExecCommand(lg, tmpFile.Name())
	log.Infof("Executing: %s\n", strings.Join(command, " "))

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

func makeExecCommand(lg LangDefinition, tempFileName string) []string {
	file := execFile(lg.sourceFileName)
	var command = []string{
		"run", "--rm", "--pull=never",
		"-v", fmt.Sprintf("%s:%s", tempFileName, file),
	}
	command = append(command, lg.image)
	command = append(command, lg.execCommand...)
	command = append(command, file)
	return command
}
