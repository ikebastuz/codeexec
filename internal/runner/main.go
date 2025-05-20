package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"

	log "github.com/sirupsen/logrus"
)

func Run(lang Lang, code string) (string, string, error) {
	var lg, ok = LangDefinitions[lang]
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

	return stdout.String(), stderr.String(), err
}

func execFile(fileName string) string {
	return fmt.Sprintf("%s/%s", WORKDIR, fileName)
}

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
