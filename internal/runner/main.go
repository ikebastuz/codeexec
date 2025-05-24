package runner

import (
	"codeexec/internal/config"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func Run(lang Lang, code string) Result {
	ld, ok := LangDefinitions[lang]
	if !ok {
		return Result{Error: fmt.Sprintf("unknown language: %s", lang)}
	}

	tmpDir, err := mkWorkDir(ld, code)
	if err != nil {
		return Result{Error: err.Error()}
	}
	defer os.RemoveAll(tmpDir)

	return RunWithDependencies(ld, tmpDir)
}

func RunWithDependencies(ld LangDefinition, tempDir string) Result {
	var buildDuration float64

	if ld.buildCommand != nil {
		buildCommand := mkBuildCommand(ld, tempDir)

		start := time.Now()
		stdout, stderr, err := runCommand(buildCommand, config.PROCESS_TIMEOUT)
		buildDuration = time.Since(start).Seconds()

		if err != nil {
			log.Errorf("Failed to build: %s", err)
			return Result{
				Stdout:        stdout,
				Stderr:        stderr,
				BuildDuration: buildDuration,
				Error:         err.Error(),
			}
		}
	}

	execCommand := mkExecCommand(ld, tempDir)

	start := time.Now()
	stdout, stderr, err := runCommand(execCommand, config.PROCESS_TIMEOUT)
	execDuration := time.Since(start).Seconds()

	result := Result{
		Stdout:        stdout,
		Stderr:        stderr,
		BuildDuration: buildDuration,
		ExecDuration:  execDuration,
	}

	if err != nil {
		log.Errorf("Execution error: %s", err)
		result.Error = err.Error()
	}
	return result
}
