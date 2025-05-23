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

func runCommand(command []string, timeout time.Duration) (string, string, error) {
	var stdout, stderr bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", command...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	log.Infof("Running command: %s", strings.Join(command, " "))
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
		// TODO: rewrite ugly logic
		if lg.buildCommand == nil {
			command = append(command, lg.sourceFileName)
		}
	}
	return command
}

func Run(lang Lang, code string) (string, string, float64, float64, error) {
	lg, ok := LangDefinitions[lang]
	if !ok {
		return "", "", 0, 0, fmt.Errorf("unknown language: %s", lang)
	}

	tmpDir, err := mkWorkDir(lg, code)
	if err != nil {
		return "", "", 0, 0, err
	}
	defer os.RemoveAll(tmpDir)

	var buildDuration float64
	if lg.buildCommand != nil {
		start := time.Now()
		stdout, stderr, err := runCommand(mkCommand(lg, tmpDir, true), config.PROCESS_TIMEOUT)
		buildDuration = time.Since(start).Seconds()
		if err != nil {
			log.Errorf("Failed to build: %s", err)
			return stdout, stderr, buildDuration, 0, err
		}
	}

	start := time.Now()
	stdout, stderr, err := runCommand(mkCommand(lg, tmpDir, false), config.PROCESS_TIMEOUT)
	execDuration := time.Since(start).Seconds()
	if err != nil {
		log.Errorf("Execution error: %s", err)
	}
	return stdout, stderr, buildDuration, execDuration, err
}
