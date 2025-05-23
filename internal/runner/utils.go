package runner

import (
	"codeexec/internal/config"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
