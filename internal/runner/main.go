package runner

import (
	"bytes"
	"codeexec/internal/config"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func mkWorkDir(lg LangDefinition, sourceCode string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "tmp-app-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	sourceFilePath := filepath.Join(tmpDir, lg.sourceFileName)
	if err := os.WriteFile(sourceFilePath, []byte(sourceCode), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to write code to source file: %w", err)
	}
	return tmpDir, nil
}

func runDockerCommand(command []string, timeout time.Duration) (string, string, error) {
	var stdout, stderr bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", command...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return stdout.String(), stderr.String(), errors.New("timeout exceeded")
	}
	return stdout.String(), stderr.String(), err
}

func mkCommand(lg LangDefinition, tempDir string, isBuild bool) []string {
	command := []string{
		"run", "--rm",
		"--pull=never",
		"-w", WORKDIR,
		"-v", fmt.Sprintf("%s:%s", tempDir, WORKDIR),
	}
	command = append(command, lg.image)
	if isBuild {
		command = append(command, lg.buildCommand...)
	} else {
		command = append(command, lg.execCommand...)
	}
	return command
}

func Run(lang Lang, code string) (string, string, float64, error) {
	lg, ok := LangDefinitions[lang]
	if !ok {
		return "", "", 0, fmt.Errorf("unknown language: %s", lang)
	}

	tmpDir, err := mkWorkDir(lg, code)
	if err != nil {
		return "", "", 0, err
	}
	defer os.RemoveAll(tmpDir)

	if lg.buildCommand != nil {
		stdout, stderr, err := runDockerCommand(mkCommand(lg, tmpDir, true), config.PROCESS_TIMEOUT)
		if err != nil {
			log.Errorf("Failed to build: %s", err)
			return stdout, stderr, 0, err
		}
	}

	start := time.Now()
	stdout, stderr, err := runDockerCommand(mkCommand(lg, tmpDir, false), config.PROCESS_TIMEOUT)
	duration := time.Since(start).Seconds()
	if err != nil {
		log.Errorf("Execution error: %s", err)
	}
	return stdout, stderr, duration, err
}
