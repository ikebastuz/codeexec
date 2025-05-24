package runner

import (
	"codeexec/internal/config"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Runner struct {
	lang Lang
	code string
}

func NewRunner(lang Lang, code string) *Runner {
	return &Runner{lang: lang, code: code}
}

func (r *Runner) Run() Result {
	fs := &TempDirFS{}
	return r.RunWithDeps(fs)
}

func (r *Runner) RunWithDeps(fs FS) Result {
	ld, ok := LangDefinitions[r.lang]
	if !ok {
		return Result{Error: fmt.Sprintf("unknown language: %s", r.lang)}
	}

	tempDir, err := fs.Create(ld.sourceFileName, r.code)
	if err != nil {
		return Result{Error: err.Error()}
	}
	defer fs.Cleanup()

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
