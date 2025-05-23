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
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func PullAllImages() {
	var wg sync.WaitGroup

	for _, def := range LangDefinitions {
		image := def.image
		wg.Add(1)
		go func(img string) {
			defer wg.Done()
			log.Infof("Pulling image %s", img)
			cmd := exec.Command("docker", "pull", img)
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Warnf("Failed to pull image %s: %v\nOutput: %s", img, err, string(out))
			} else {
				log.Infof("Successfully pulled image %s", img)
			}
		}(image)
	}
	wg.Wait()
}

func StartImageMonitor() {
	ticker := time.NewTicker(config.CHECK_IMAGES_INTERVAL)
	defer ticker.Stop()
	for {
		log.Infof("Checking docker images...")
		missing := false
		for _, def := range LangDefinitions {
			img := def.image
			cmd := exec.Command("docker", "image", "inspect", img)
			if err := cmd.Run(); err != nil {
				log.Warnf("Image %s not found, will pull all images...", img)
				missing = true
				break
			}
		}
		if missing {
			PullAllImages()
		} else {
			log.Infof("All images are in place")
		}
		<-ticker.C
	}
}

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
		return stdout.String(), stderr.String(), errors.New(ERR_TIMEOUT_EXCEEDED)
	}

	return stdout.String(), stderr.String(), err
}

func mkDockerBaseCommand(tempDir string) []string {
	return []string{
		"run", "--rm",
		"--pull=never",
		"-w", WORKDIR,
		"-v", fmt.Sprintf("%s:%s", tempDir, WORKDIR),
	}
}

func mkBuildCommand(ld LangDefinition, tempDir string) []string {
	if ld.buildCommand == nil {
		return nil
	}

	command := mkDockerBaseCommand(tempDir)

	command = append(command, ld.image)
	command = append(command, ld.buildCommand...)
	return command
}

func mkExecCommand(ld LangDefinition, tempDir string) []string {
	command := mkDockerBaseCommand(tempDir)

	command = append(command, ld.image)
	command = append(command, ld.execCommand...)

	if ld.buildCommand == nil {
		command = append(command, ld.sourceFileName)
	}

	return command
}
