package runner

import (
	"codeexec/internal/config"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func Run(lang Lang, code string) (string, string, float64, float64, error) {
	ld, ok := LangDefinitions[lang]
	if !ok {
		return "", "", 0, 0, fmt.Errorf("unknown language: %s", lang)
	}

	tmpDir, err := mkWorkDir(ld, code)
	if err != nil {
		return "", "", 0, 0, err
	}
	defer os.RemoveAll(tmpDir)

	var buildDuration float64

	if ld.buildCommand != nil {
		buildCommand := mkBuildCommand(ld, tmpDir)

		start := time.Now()
		stdout, stderr, err := runCommand(buildCommand, config.PROCESS_TIMEOUT)
		buildDuration = time.Since(start).Seconds()

		if err != nil {
			log.Errorf("Failed to build: %s", err)
			return stdout, stderr, buildDuration, 0, err
		}
	}

	execCommand := mkExecCommand(ld, tmpDir)

	start := time.Now()
	stdout, stderr, err := runCommand(execCommand, config.PROCESS_TIMEOUT)
	execDuration := time.Since(start).Seconds()

	if err != nil {
		log.Errorf("Execution error: %s", err)
	}
	return stdout, stderr, buildDuration, execDuration, err
}
