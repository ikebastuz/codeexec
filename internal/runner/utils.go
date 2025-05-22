package runner

import (
	"codeexec/internal/config"
	"fmt"
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

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
